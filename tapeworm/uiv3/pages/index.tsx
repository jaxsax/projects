import { useState, Fragment, useEffect, useRef } from 'react'
import { useQuery } from 'react-query';
import lunr from 'lunr';
import { useAsync } from 'react-async-hook'
import useConstant from 'use-constant'
import AwesomeDebouncePromise from 'awesome-debounce-promise'
import formatDistance from 'date-fns/formatDistance';
import formatISO from 'date-fns/formatISO';
import { Menu, Transition } from '@headlessui/react'
import { ChevronDownIcon } from '@heroicons/react/solid'

async function getLinks() {
    return fetch('/api/links').then((r) => r.json())
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
    searchDurationSeconds: number | null
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

    const date = new Date(l.created_ts * 1000)
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
            <p className="text-gray-400" title={formatISO(date)}>
                {formatDistance(date, new Date(), { addSuffix: true })}
            </p>
        </div>
    )
}

function LinksContainer({ links, searchDurationSeconds }: Props) {
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
        <div className="flex justify-between">
            <p className="text-gray-500">{links.length} records found {searchDurationSeconds ? `in ${searchDurationSeconds} ms` : ''} </p>
            <div className="text-right">
                <Menu as="div" className="relative inline-block text-left border-2 border-blue-600 px-4 py-2 rounded-md">
                    <Menu.Button>Sort by</Menu.Button>
                    <Menu.Items className="origin-top-right absolute right-0 mt-2 w-56 rounded-md shadow-lg bg-white ring-1 ring-black ring-opacity-5 focus:outline-none">
                        <div className="py-1 px-1">
                            <Menu.Item>
                                <div className="px-4 py-2 text-sm rounded-md hover:bg-purple-600 hover:text-white">
                                    Newest
                                </div>
                            </Menu.Item>
                            <Menu.Item>
                                <div className="px-4 py-2 text-sm rounded-md hover:bg-purple-600 hover:text-white">
                                    Oldest
                                </div>
                            </Menu.Item>
                        </div>
                    </Menu.Items>
                </Menu>
            </div>
        </div>
        <div className="mt-4 flex flex-row flex-wrap gap-4">
            {links.map((l) => <LinkItem key={l.id} {...l} />)}
        </div>
    </Wrapper >
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
    const { isSuccess, isLoading, isError, data, dataUpdatedAt } = useLinks()
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

    const loading = (isLoading || searchResults.loading) && !isError
    const done = isSuccess && searchResults.result
    return (
        <>
            <div className="mx-2 xl:container xl:mx-auto mt-24">
                <h1 className="text-center text-7xl text-bold">link search</h1>
                {dataUpdatedAt
                    ? (
                        <p className="text-center text-gray-400 text-md">
                            index last rebuilt {' '}
                            {
                                formatDistance(new Date(dataUpdatedAt), new Date(),
                                    { addSuffix: true, includeSeconds: true })
                            }
                        </p>
                    )
                    : null
                }
                <div className="mt-4">
                    <input
                        type="text" name="query" placeholder="Enter search terms"
                        value={inputText}
                        onChange={(e) => setInputText(e.target.value)}
                        className="w-full px-4 py-2 border-2 border-gray-400 outline-none focus:border-gray-400 focus:border-blue-400" />
                </div>
                <div className="mt-2 mb-4">
                    {isError ? <h1 className="text-center text-2xl">Failed to retrieve links</h1> : null}
                    {loading ? <h1 className="text-center text-2xl">Loading...</h1> : null}
                    {done
                        ? <LinksContainer
                            links={items}
                            searchDurationSeconds={searchDuration} />
                        : null}
                </div>
            </div>
        </>
    )
}

export default IndexPage;
