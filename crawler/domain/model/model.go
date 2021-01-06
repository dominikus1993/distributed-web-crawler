package model

type CrawlWebsite struct {
	Url string `json:"url,omitempty"`
}

type CrawledWebsite struct {
	Url string `json:"url,omitempty"`
}

func NewCrawlWebsite(url string) *CrawlWebsite {
	return &CrawlWebsite{Url: url}
}

func NewCrawledWebsite(url string) *CrawledWebsite {
	return &CrawledWebsite{Url: url}
}
