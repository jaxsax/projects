export type LinkItem = {
  id: number;
  link: string;
  title: string;
  domain: string;
  path: string;
  created_ts: number;
  created_by: number;
  labels: Array<LinkLabel>;
};

export type LinkLabel = {
  name: string;
};

export type GetLinkResponse = {
  links: Array<LinkItem>;
};
