import { useQuery } from 'react-query';

function sleep(ms: number) {
    return new Promise((resolve) => setTimeout(resolve, ms))
}

async function getLinks() {
    await sleep(1000)

    return fetch('https://jaxsax.co/api/links').then((r) => r.json())
}

function useLinks() {
    return useQuery(['links'], getLinks)
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

function IndexPage() {
    const { isSuccess, isLoading, data } = useLinks()
    return (
        <>
            <div className="mx-2 xl:container xl:mx-auto mt-24">
                <h1 className="text-center text-7xl text-bold">link search</h1>
                <form className="mt-4">
                    <input
                        type="text" name="query" placeholder="Enter search terms"
                        className="w-full px-4 py-2 border-2 border-gray-400 outline-none focus:border-gray-400" />
                </form>
                <div className="mt-4">
                    {isLoading ? <h1 className="text-center text-2xl">Loading...</h1> : null}
                    {isSuccess ? <LinksContainer links={data.links.slice(0, 5)} /> : null}
                </div>
            </div>
        </>
    )
}

export default IndexPage;
