import { GetLinkResponse } from "~/models/link";

let apiHost =
  process.env.NODE_ENV === "production"
    ? "https://jaxsax.co"
    : "http://localhost:8081";

export async function ListLinks() {
  let resp = await fetch(apiURL("/api/links"), {
    method: "GET",
  });
  let body = await resp.json();

  return body;
}

export async function GetLink(url: string): Promise<GetLinkResponse> {
  let resp = await fetch(
    withQuery(apiURL("/api/links/get"), {
      url,
    }),
    {
      method: "GET",
    }
  );

  return resp.json();
}

function apiURL(url: string) {
  return new URL(url, apiHost).href;
}

function withQuery(url: string, parameters: Record<string, string>) {
  let u = new URL(url);
  u.search = new URLSearchParams(parameters).toString();
  return u.href;
}