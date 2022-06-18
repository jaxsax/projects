export type LinkItem = {
  id: number;
  link: string;
  title: string;
  domain: string;
  created_ts: number;
  created_by: number;
};

export type GetLinkResponse = {
  links: Array<LinkItem>;
};
