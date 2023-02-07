export type LinkItem = {
    id: number;
    link: string;
    title: string;
    domain: string;
    path: string;
    created_ts: number;
    created_by: number;
    labels: Array<LinkLabel>;
    dimensions: Array<LinkDimension>;
};

export type LinkLabel = {
    name: string;
};

export type LinkDimension = {
    kind: string;
    data: any
}

export type GetLinkResponse = {
    links: Array<LinkItem>;
};
