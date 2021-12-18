import { useState } from 'react'
import { useQuery } from 'react-query';
import lunr from 'lunr';
import { useAsync } from 'react-async-hook'
import useConstant from 'use-constant'
import AwesomeDebouncePromise from 'awesome-debounce-promise'

function sleep(ms: number) {
    return new Promise((resolve) => setTimeout(resolve, ms))
}

async function getLinks() {
    await sleep(1000)

    return fetch('https://jaxsax.co/api/links').then((r) => r.json())
}

let index = lunr(() => { })

function useLinks() {
    return useQuery(['links'], getLinks, {
        onSuccess: (data) => {
            const tmpIndex = lunr((inst) => {
                inst.ref('id')
                inst.field('title')

                data.links.forEach((e) => {
                    inst.add(e)
                })
            })

            index = tmpIndex
        },
    })
}


type Link = {
    id: string
    created_ts: string
    created_by: number
    link: string
    title: string
}

type Props = {
    links: Link[]
}

function LinkItem(l: Link) {
    let title = l.title.replace(/(\r\n|\n|\r)/gm, " ");
    title = l.title.replace(/(\s+)/g, " ");

    let linkHostname = null;
    try {
        const { hostname } = new URL(l.link)
        linkHostname = hostname
    } catch (error) {
    }

    return (
        <div className="w-full">
            <div>
                <a href={l.link}>
                    {title}
                    {linkHostname !== null ?
                        (
                            <div className="text-gray-400">
                                ({linkHostname})
                            </div>
                        ) : null}
                </a>
            </div>
        </div>
    )
}

function LinksContainer({ links }: Props) {
    return (
        <div className="flex flex-row flex-wrap gap-4">
            {links.map((l) => <LinkItem key={l.id} {...l} />)}
        </div>
    )
}

// Generic reusable hook
const useDebouncedSearch = (searchFunction: (term: string) => any) => {

    // Handle the input text state
    const [inputText, setInputText] = useState('');

    // Debounce the original search async function
    const debouncedSearchFunction = useConstant(() =>
        AwesomeDebouncePromise(searchFunction, 1000)
    );

    // The async callback is run each time the text changes,
    // but as the search function is debounced, it does not
    // fire a new request on each keystroke
    const searchResults = useAsync(
        async () => {
            if (inputText.length === 0) {
                return [];
            } else {
                return debouncedSearchFunction(inputText);
            }
        },
        [debouncedSearchFunction, inputText]
    );

    // Return everything needed for the hook consumer
    return {
        inputText,
        setInputText,
        searchResults,
    };
};

function IndexPage() {
    const { isSuccess, isLoading, data } = useLinks()
    const { inputText, setInputText, searchResults } = useDebouncedSearch((term: string) => {
        const t0 = performance.now()
        let r = index.search(term)
        const t1 = performance.now()

        console.log(`finished in ${t1 - t0} ms`)

        return r
    })

    let items = data?.links
    if (!searchResults.loading && data) {
        const filteredItemIds = new Set(searchResults.result.map(x => x.ref))
        items = items.filter(x => filteredItemIds.has(String(x.id)))
    }

    return (
        <>
            <div className="mx-2 xl:container xl:mx-auto mt-24">
                <h1 className="text-center text-7xl text-bold">link search</h1>
                <form className="mt-4">
                    <input
                        type="text" name="query" placeholder="Enter search terms"
                        value={inputText}
                        onChange={(e) => setInputText(e.target.value)}
                        className="w-full px-4 py-2 border-2 border-gray-400 outline-none focus:border-gray-400" />
                </form>
                <div className="mt-4">
                    {isLoading || searchResults.loading ? <h1 className="text-center text-2xl">Loading...</h1> : null}
                    {isSuccess || !searchResults.loading ? <LinksContainer links={items} /> : null}
                </div>
            </div>
        </>
    )
}

export default IndexPage;
