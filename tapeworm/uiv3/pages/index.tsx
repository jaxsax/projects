import { useQuery  } from 'react-query';

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
    title: Map<string, any>
}

type Props = {
    links: Link[]
}

function LinksContainer(data: Props) {
    console.debug(data)
    return (
        <div className="flex">
        </div>
    )
}

function IndexPage() {
    const { isSuccess, isLoading, data } = useLinks()
    return (
        <>
            <div className="container mx-auto mt-12">
                <h1 className="text-3xl text-bold">link search</h1>
                <form className="mt-4">
                    <input
                        type="text" name="query" placeholder="Enter search terms"
                        className="w-1/3 px-4 py-2 border-2 border-gray-400 outline-none  focus:border-gray-400" />
                </form>
                <div className="mt-4">
                    {isLoading ? <h1 className="text-2xl font-bold">Loading...</h1> : null}
                    {isSuccess ? <LinksContainer links={data} /> : null}
                </div>
            </div>
        </>
    )
}

export default IndexPage;
