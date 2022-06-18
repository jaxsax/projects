import { json, LoaderFunction } from "@remix-run/node";
import { useLoaderData } from "@remix-run/react";
import formatDistance from "date-fns/formatDistance";
import formatISO from "date-fns/formatISO";
import { ClientOnly } from "remix-utils";
import { LinkItem } from "~/models/link";
import { GetLink } from "~/utils/api.server";

type LoaderData = {
  links: Array<LinkItem>;
  fromLink: LinkItem;
};

export const loader: LoaderFunction = async ({ request }) => {
  let searchParams = new URL(request.url).searchParams;
  let url = searchParams.get("url");

  let links = await GetLink(url);

  let originLinkID = parseInt(searchParams.get("from"), 10);
  let originLink = links.links.filter((x) => x.id == originLinkID);
  let otherLinks = links.links.filter((x) => x.id !== originLinkID);
  return json<LoaderData>({
    links: otherLinks,
    fromLink: originLink[0],
  });
};

const OtherLinkItems: React.FC<{ items: Array<LinkItem> }> = ({ items }) => {
  return (
    <div className="mt-4">
      <h1 className="text-4xl text-bold">Related links</h1>
      {items.map((x) => (
        <SingleLink key={x.id} {...x} />
      ))}
    </div>
  );
};

const SingleLink: React.FC<LinkItem> = ({
  title,
  link,
  created_ts,
  created_by,
}) => {
  return (
    <div>
      <div>
        <h1 className="text-2xl text-bold">title</h1>
        <p>{title}</p>
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
    </div>
  );
};

export default function LinkDetailPage() {
  let { links, fromLink } = useLoaderData<LoaderData>();
  return (
    <div className="mx-2 xl:container xl:mx-auto mt-12 min-h-screen">
      <h1 className="text-center text-7xl text-bold">link detail page</h1>
      <div className="mt-4">
        <SingleLink {...fromLink} />
      </div>
      {links.length > 0 ? <OtherLinkItems items={links} /> : null}
    </div>
  );
}
