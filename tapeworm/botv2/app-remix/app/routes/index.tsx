import { json, LoaderFunction } from "@remix-run/node";
import {
  Form,
  useLoaderData,
  useSubmit,
  useTransition,
} from "@remix-run/react";
import formatDistance from "date-fns/formatDistance";
import formatISO from "date-fns/formatISO";
import React from "react";
import { ClientOnly } from "remix-utils";

let apiHost =
  process.env.NODE_ENV === "production"
    ? "https://jaxsax.co"
    : "http://localhost:8081";

type Link = {
  id: number;
  link: string;
  title: string;
  domain: string;
  created_ts: number;
  created_by: number;
};

type LoaderData = {
  items: Array<Link>;
  q: string | undefined;
};

export const loader: LoaderFunction = async ({ request }) => {
  const url = new URL(request.url);
  let term = url.searchParams.get("q") ?? "";
  term = term.trim();

  const apiRequest = new Request(`${apiHost}/api/links`, {
    method: "GET",
  });

  const resp = await fetch(apiRequest);
  const body = await resp.json();

  let filteredLinks = [];
  if (term) {
    filteredLinks = body["links"].filter((x: Link) => {
      return x.title.toLowerCase().includes(term.toLowerCase());
    });
  } else {
    filteredLinks = body["links"];
  }

  return json<LoaderData>({
    items: filteredLinks,
    q: term,
  });
};

const LinkTimingInfo: React.FC<{ created_ts: number }> = ({ created_ts }) => {
  const date = new Date(created_ts * 1000);
  return (
    <p className="text-gray-400" title={formatISO(date)}>
      {formatDistance(date, new Date(), { addSuffix: true })}
    </p>
  );
};

const LinkItem: React.FC<Link> = ({ link, domain, created_ts, title }) => {
  return (
    <div className="truncate">
      <a href={link} className="text-lg">
        {title}
      </a>
      {domain ? (
        <div className="block text-gray-400 md:inline md:ml-2">({domain})</div>
      ) : null}

      <ClientOnly fallback={<div>...</div>}>
        {() => <LinkTimingInfo created_ts={created_ts} />}
      </ClientOnly>
    </div>
  );
};

export default function Index() {
  const { items, q = "" } = useLoaderData<LoaderData>();
  const submit = useSubmit();
  const transition = useTransition();

  function handleChange(event: React.ChangeEvent<HTMLFormElement>) {
    submit(event.currentTarget, { replace: true });
  }

  return (
    <div className="mx-2 xl:container xl:mx-auto mt-24 min-h-screen">
      <h1 className="text-center text-7xl text-bold">link search</h1>
      <div className="mt-4">
        <Form onChange={handleChange}>
          <input
            autoFocus
            type="text"
            name="q"
            placeholder="Enter search terms"
            className="w-full px-4 py-2 border-2 border-gray-400 outline-none focus:border-gray-400 focus:border-blue-400"
            defaultValue={q}
          />
        </Form>
      </div>
      <div className="mt-2 mb-4">
        {transition.state === "submitting" ? (
          <p>Loading...</p>
        ) : (
          <>
            {items && items.map((i) => <LinkItem key={i.id} {...i} />)}
            {items?.length == 0 && <div>No results</div>}
          </>
        )}
      </div>
    </div>
  );
}
