import { useState } from 'react'
import { useQuery } from 'react-query';
import lunr from 'lunr';
import { useAsync } from 'react-async-hook'
import useConstant from 'use-constant'
import AwesomeDebouncePromise from 'awesome-debounce-promise'
import formatDistance from 'date-fns/formatDistance';

async function getLinks() {
    return fetch('https://jaxsax.co/api/links').then((r) => r.json())
}

let index = lunr(() => { })

function useLinks() {
    return useQuery(['links'], getLinks, {
        onSuccess: (data) => {
            let t0 = performance.now()
            const tmpIndex = lunr((inst) => {
                inst.ref('id')
                inst.field('title')

                data.links.forEach((e: any) => {
                    inst.add(e)
                })
            })
            console.log(`index rebuilt in ${performance.now() - t0} ms`)

            index = tmpIndex
        },
    })
}


type Link = {
    id: string
    created_ts: number
    created_by: number
    link: string
    title: string
}

type Props = {
    links: Link[]
    searchDurationMs: number | null
}

function LinkItem(l: Link) {
    let title = l.title.replace(/(\r\n|\n|\r)/gm, " ");
    title = l.title.replace(/(\s+)/g, " ");

    let linkHostname: string | null = null;
    try {
        const { hostname } = new URL(l.link)
        linkHostname = hostname
    } catch (error) {
    }

    return (
        <div className="truncate" style={{ flexBasis: '100%' }}>
            <a href={l.link} className="text-lg">
                {title}
            </a>
            {linkHostname !== null ?
                (
                    <div className="text-gray-400 md:inline md:ml-2">
                        ({linkHostname})
                    </div>
                ) : null}
            <p className="text-gray-400">
                {formatDistance(new Date(l.created_ts * 1000), new Date(), { addSuffix: true })}
            </p>
        </div>
    )
}

function LinksContainer({ links, searchDurationMs }: Props) {
    const Wrapper = ({ children }) => {
        return <div>
            {children}
        </div>
    }

    if (links.length === 0) {
        return <Wrapper>
            <span className="text-xl">No results found</span>
        </Wrapper>
    }

    return <Wrapper>
        <p className="text-gray-500">{links.length} records found {searchDurationMs ? `in ${searchDurationMs} ms` : ''} </p>
        <div className="mt-4 flex flex-row flex-wrap gap-4">
            {links.map((l) => <LinkItem key={l.id} {...l} />)}
        </div>
    </Wrapper>
}

// Generic reusable hook
// https://stackoverflow.com/questions/23123138/perform-debounce-in-react-js
const useDebouncedSearch = (searchFunction: (term: string) => any) => {

    // Handle the input text state
    const [inputText, setInputText] = useState('');

    // Debounce the original search async function
    const debouncedSearchFunction = useConstant(() =>
        AwesomeDebouncePromise(searchFunction, 100)
    );

    // The async callback is run each time the text changes,
    // but as the search function is debounced, it does not
    // fire a new request on each keystroke
    const searchResults = useAsync(
        async () => {
            if (inputText.length === 0) {
                return {};
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

        return {
            data: r,
            elapsed: t1 - t0,
        }
    })

    let items = data?.links
    let searchDuration = null
    if (!searchResults.loading && data) {
        const filteredItemIds = new Set(searchResults.result?.data?.map(x => x.ref))
        items = items.filter(x => filteredItemIds.has(String(x.id)))
        searchDuration = searchResults.result?.elapsed
    }

    if (!items || items.length === 0) {
        if (inputText === '') {
            items = data?.links
        }
    }

    const loading = isLoading || searchResults.loading
    const done = isSuccess && searchResults.result
    return (
        <>
            <div className="mx-2 xl:container xl:mx-auto mt-24">
                <h1 className="text-center text-7xl text-bold">link search</h1>
                <div className="mt-4">
                    <input
                        type="text" name="query" placeholder="Enter search terms"
                        value={inputText}
                        onChange={(e) => setInputText(e.target.value)}
                        className="w-full px-4 py-2 border-2 border-gray-400 outline-none focus:border-gray-400 focus:border-blue-400" />
                </div>
                <div className="mt-2 mb-4">
                    {loading ? <h1 className="text-center text-2xl">Loading...</h1> : null}
                    {done ? <LinksContainer links={items} searchDurationMs={searchDuration} /> : null}
                </div>
            </div>
        </>
    )
}

export default IndexPage;
