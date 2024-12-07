package model

type SearchResult struct {
	Url           string
	BookName      string
	Author        string
	Intro         string
	LatestChapter string
	LatestUpdate  string
}

type CrawlResult struct {
	TakeTime int64
}
