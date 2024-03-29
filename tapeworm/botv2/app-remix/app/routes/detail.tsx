import { json, LoaderFunction } from "@remix-run/node";
import { useLoaderData } from "@remix-run/react";
import formatDistance from "date-fns/formatDistance";
import formatISO from "date-fns/formatISO";
import { ClientOnly } from "remix-utils";
import { LinkDimension, LinkItem } from "~/models/link";
import { GetLink, ListLinksByDomain } from "~/utils/api.server";
import { Hn } from "~/icons"

type LoaderData = {
    links: Array<LinkItem>;
    fromLink: LinkItem;
    domain: string;
    fromSameDomain: Array<LinkItem>;
};

export const loader: LoaderFunction = async ({ request }) => {
    let searchParams = new URL(request.url).searchParams;
    let url = searchParams.get("url") ?? "";

    let links = await GetLink(url);

    let originLinkID = parseInt(searchParams.get("from") ?? "", 10);
    let originLinks = links.links.filter((x) => x.id == originLinkID);
    let otherLinks = links.links.filter((x) => x.id !== originLinkID);
    let originLink = originLinks[0];

    let linksByDomain = await ListLinksByDomain(originLink.domain);
    let sameDomainLinks = linksByDomain.links.filter(
        (x) => x.id !== originLink.id
    );

    return json<LoaderData>({
        links: otherLinks,
        fromLink: originLink,
        domain: originLink.domain,
        fromSameDomain: sameDomainLinks,
    });
};

const OtherLinkItems: React.FC<{ items: Array<LinkItem> }> = ({ items }) => {
    return (
        <div className="mt-8">
            <h1 className="text-4xl text-bold">Also added previously</h1>
            {items.map((x) => (
                <div key={x.id} className="mt-2">
                    <SingleLinkV2 {...x} />
                </div>
            ))}
        </div>
    );
};

const SameDomainLinkItems: React.FC<{
    domain: string;
    items: Array<LinkItem>;
}> = ({ domain, items }) => {
    return (
        <div className="mt-8">
            <h1 className="text-4xl text-bold">Also from {domain}</h1>
            {items.map((x) => (
                <div key={x.id} className="mt-2">
                    <SingleLinkV2 {...x} />
                </div>
            ))}
        </div>
    );
};

const SingleLinkV2: React.FC<LinkItem> = ({
    title,
    link,
    created_ts,
    created_by,
}) => {
    return (
        <div>
            <div>
                <a className="text-blue-500" href={link}>
                    {title}
                </a>
                <div className="flex md:inline-flex md:ml-2 flex-row mt-2 justify-start gap-x-2">
                    <ClientOnly fallback={"..."}>
                        {() => {
                            const date = new Date(created_ts * 1000);
                            const distanceDate = formatDistance(date, new Date(), {
                                addSuffix: true,
                            });
                            return (
                                <span
                                    title={formatISO(date)}
                                    className="border-2 border-lime-600 p-2"
                                >
                                    {distanceDate}
                                </span>
                            );
                        }}
                    </ClientOnly>
                    <span
                        className="inline-block border-2 border-red-400 p-2"
                        title="User ID"
                    >
                        {created_by}
                    </span>
                </div>
            </div>
        </div>
    );
};

let colors = [
    "bg-red-500",
    "bg-orange-400",
    "bg-yellow-500",
    "bg-green-400",
    "bg-blue-400",
    "bg-slate-400",
];

function gencolor(x) {
    return colors[Math.floor(Math.random() * colors.length)];
}

const SingleLink: React.FC<LinkItem> = ({
    title,
    link,
    created_ts,
    created_by,
    labels,
    dimensions,
}) => {
    return (
        <div>
            <div>
                <h1 className="text-2xl text-bold">title</h1>
                <p>{title}</p>
                {labels ? (
                    <div className="flex flex-wrap gap-2">
                        {labels.map((x) => (
                            <span
                                key={x.name}
                                className={` ${gencolor(x.name)} shadow border-1 p-2 rounded`}
                            >
                                {x.name}
                            </span>
                        ))}
                    </div>
                ) : null}
            </div>
            <div>
                <h1 className="text-2xl text-bold">link</h1>
                <p>{link}</p>
            </div>
            <div>
                <h1 className="text-2xl text-bold">created at</h1>
                <ClientOnly fallback={"..."}>
                    {() => {
                        const date = new Date(created_ts * 1000);
                        return (
                            <p title={formatDistance(date, new Date(), { addSuffix: true })}>
                                {formatISO(date)}
                            </p>
                        );
                    }}
                </ClientOnly>
            </div>
            <div>
                <h1 className="text-2xl text-bold">created by</h1>
                <p>{created_by}</p>
            </div>
            {dimensions && <div>
                <h1 className="text-2xl text-bold">aggregators</h1>
                {dimensions.map((x) => <Dimension {...x} />)}
            </div>}
        </div>
    );
};

function Dimension({ kind, data }: LinkDimension) {
    return <div className="mt-2 flex">
        {kind === "hn_thread" && (
            <a href={data.url}>
                <div style={{ background: "#ff6600", border: "solid", borderColor: "#ff6600" }}>
                    <Hn width={24} height={24} fill="#fff" />
                </div>
            </a>
        )}
    </div>
}

export default function LinkDetailPage() {
    let { links, fromLink, domain, fromSameDomain } = useLoaderData<LoaderData>();
    return (
        <div className="mx-2 xl:container xl:mx-auto mt-12 min-h-screen">
            <h1 className="text-center text-7xl text-bold">link detail page</h1>
            <div className="mt-4">
                <SingleLink {...fromLink} />
            </div>
            {links.length > 0 ? <OtherLinkItems items={links} /> : null}
            {fromSameDomain.length > 0 ? (
                <SameDomainLinkItems domain={domain} items={fromSameDomain} />
            ) : null}
        </div>
    );
}
