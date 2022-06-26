import { json, LoaderFunction } from "@remix-run/node";
import {
  Form,
  Link,
  useLoaderData,
  useSearchParams,
  useSubmit,
  useTransition,
} from "@remix-run/react";
import formatDistance from "date-fns/formatDistance";
import formatISO from "date-fns/formatISO";
import { useMemo } from "react";
import { ClientOnly } from "remix-utils";
import { classNames } from "~/lib/classnames";

import { LinkItem } from "~/models/link";
import { ListLinks } from "~/utils/api.server";

type LoaderData = {
  items: Array<LinkItem>;
  totalCount: number;
  itemsPerPage: number;
  q: string | undefined;
};

export const loader: LoaderFunction = async ({ request }) => {
  const url = new URL(request.url);
  let term = url.searchParams.get("q") ?? "";
  term = term.trim();

  let limit = url.searchParams.get("limit") ?? "";
  let page = url.searchParams.get("page") ?? "";

  let body = await ListLinks(term, page, limit);

  return json<LoaderData>({
    items: body.links,
    q: term,
    totalCount: body.total,
    itemsPerPage: body.items_per_page,
  });
};

const LinkTimingInfo: React.FC<{ created_ts: number }> = ({ created_ts }) => {
  const date = new Date(created_ts * 1000);
  return (
    <span title={formatISO(date)}>
      {formatDistance(date, new Date(), { addSuffix: true })}
    </span>
  );
};

const LinkItem: React.FC<LinkItem> = ({
  id,
  link,
  domain,
  created_ts,
  title,
}) => {
  return (
    <div className="truncate">
      <a href={link} className="text-lg">
        {title}
      </a>
      {domain ? (
        <div className="block text-gray-400 md:inline md:ml-2">({domain})</div>
      ) : null}

      <div className="text-gray-400">
        <ClientOnly fallback={<span>...</span>}>
          {() => <LinkTimingInfo created_ts={created_ts} />}
        </ClientOnly>
        {" | "}
        <span className="hover:underline">
          <Link to={`/detail?url=${encodeURIComponent(link)}&from=${id}`}>
            details
          </Link>
        </span>
      </div>
    </div>
  );
};

export default function Index() {
  const {
    items,
    totalCount,
    itemsPerPage,
    q = "",
  } = useLoaderData<LoaderData>();
  const submit = useSubmit();
  const transition = useTransition();
  const [searchParams, setSearchParams] = useSearchParams();

  const pagination = useMemo(() => {
    const currentPage = parseInt(searchParams.get("page") ?? "1", 10);

    const prevPageNumber = currentPage - 1 || 1;
    const prevPageParam = new URLSearchParams(searchParams);
    prevPageParam.set("page", prevPageNumber.toString());

    const nextPageNumber = currentPage + 1;
    const nextPageParam = new URLSearchParams(searchParams);
    nextPageParam.set("page", nextPageNumber.toString());

    return {
      prev: {
        disabled: currentPage - 1 <= 0,
        link: `?${prevPageParam.toString()}`,
      },
      next: {
        disabled: currentPage * itemsPerPage >= totalCount,
        link: `?${nextPageParam.toString()}`,
      },
    };
  }, [searchParams]);

  function onChangeLimit(event: React.ChangeEvent<HTMLSelectElement>) {
    const newParam = new URLSearchParams(searchParams);
    newParam.set("limit", event.target.value);
    setSearchParams(newParam);
  }

  function handleChange(event: React.ChangeEvent<HTMLFormElement>) {
    submit(event.currentTarget, { replace: true });
  }

  const limitOptions = [1, 15, 30, 50, 100];
  return (
    <div className="mx-2 xl:container xl:mx-auto mt-24">
      <h1 className="text-center text-7xl text-bold">link search</h1>
      <div className="mt-4">
        <Form onChange={handleChange}>
          <input
            type="text"
            name="q"
            placeholder="Enter search terms"
            className="w-full px-4 py-2 border-2 border-gray-400 outline-none focus:border-blue-400"
            defaultValue={q}
          />
        </Form>
        <div className="flex justify-between mt-1">
          <div className="text-gray-400 flex items-center">
            {totalCount} results
          </div>
          <div>
            <label htmlFor="limit">Items per page</label>
            <select
              name="limit"
              defaultValue={parseInt(
                searchParams.get("limit") || itemsPerPage.toString(),
                10
              )}
              onChange={onChangeLimit}
              className="ml-2 inline-block p-2 text-sm text-gray-900 bg-gray-50 rounded-lg border border-gray-300 focus:ring-blue-500 focus:border-blue-500 dark:bg-gray-700 dark:border-gray-600 dark:placeholder-gray-400 dark:text-white dark:focus:ring-blue-500 dark:focus:border-blue-500"
            >
              {limitOptions.map((x) => (
                <option key={x} value={x}>
                  {x}
                </option>
              ))}
            </select>
          </div>
        </div>
      </div>
      <div className="mt-2 mb-4">
        {transition.state === "submitting" ? (
          <p>Loading...</p>
        ) : (
          <>
            {items && items.map((i) => <LinkItem key={i.id} {...i} />)}
            {items && items.length > 0 && (
              <div className="mt-4">
                <ul className="flex justify-center space-x-2">
                  <li>
                    <Link
                      to={pagination.prev.link}
                      className={classNames(
                        "py-2 px-3 ml-0 leading-tight rounded-l-lg border ",
                        pagination.prev.disabled
                          ? "pointer-events-none cursor-not-allowed bg-slate-300 text-white opacity-75"
                          : null,
                        !pagination.prev.disabled
                          ? "text-gray-500 bg-white"
                          : null
                      )}
                    >
                      Previous
                    </Link>
                  </li>
                  <li>
                    <Link
                      to={pagination.next.link}
                      className={classNames(
                        "py-2 px-3 ml-0 leading-tight rounded-r-lg border",
                        pagination.next.disabled
                          ? "pointer-events-none cursor-not-allowed bg-slate-300 text-white"
                          : null,
                        !pagination.next.disabled
                          ? "text-gray-500 bg-white"
                          : null
                      )}
                    >
                      Next
                    </Link>
                  </li>
                </ul>
              </div>
            )}
            {items?.length == 0 && <div>No results</div>}
          </>
        )}
      </div>
    </div>
  );
}
