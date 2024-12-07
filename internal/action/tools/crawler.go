package tools

import (
	"fmt"
	"time"

	"github.com/767829413/easy-novel/internal/action/model"
)

type Crawler interface {
	Search(key string) []*model.SearchResult
	Crawl(res *model.SearchResult, start, end int) *model.CrawlResult
}

type novelCrawler struct{}

func NewNovelCrawler() Crawler {
	return &novelCrawler{}
}

func (nc *novelCrawler) Search(key string) []*model.SearchResult {
	// TODO: Implement search logic using a search engine
	fmt.Println("<== 正在搜索...")
	start := time.Now()
	// TODO: 解析

	duration := time.Since(start)
	return nil
}

func (nc *novelCrawler) Crawl(res *model.SearchResult, start, end int) *model.CrawlResult {
	// TODO: Implement chapter crawling logic
	return &model.CrawlResult{TakeTime: 0}
}
