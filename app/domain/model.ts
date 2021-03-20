
export interface CrawlWebsite {
    readonly url: string;
}

export interface CrawledMedia {
    readonly url: string;
    readonly contents: CrawlWebsite[];
}