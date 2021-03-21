package model

type CrawlWebsite struct {
	Url string `json:"url,omitempty" bson:"url,omitempty`
}

type Content struct {
	Url string `json:"url,omitempty"`
}
type CrawledWebsite struct {
	Url      string     `json:"url,omitempty"`
	Contents *[]Content `json:"contents,omitempty"`
}

func NewCrawlWebsite(url string) CrawlWebsite {
	return CrawlWebsite{Url: url}
}

func NewCrawledWebsite(url string, contents *[]Content) *CrawledWebsite {
	return &CrawledWebsite{Url: url, Contents: contents}
}

func NewContent(url string) Content {
	return Content{Url: url}
}
