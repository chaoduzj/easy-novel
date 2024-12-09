package parse

import (
	"bytes"
	"fmt"
	"strings"

	"github.com/767829413/easy-novel/internal/config"
	"github.com/767829413/easy-novel/internal/model"
	"github.com/767829413/easy-novel/internal/source"
	"github.com/767829413/easy-novel/internal/tools"
	"github.com/767829413/easy-novel/pkg/utils"
	"github.com/gocolly/colly/v2"
)

type ChapterParser struct {
	rule *model.Rule
}

func NewChapterParser(sourceID int) *ChapterParser {
	return &ChapterParser{
		rule: source.GetRuleBySourceID(sourceID),
	}
}

func (b *ChapterParser) Parse(
	chapter *model.Chapter,
	res *model.SearchResult,
	book *model.Book,
	bookDir string,
) (err error) {
	conf := config.GetConf()
	downOk := false
	// 抓取内容
	utils.SpinWaitMaxRetryAttempts(
		func() bool {
			var err error
			// 防止重复获取
			if !downOk {
				chapter.Content, err = b.crawl(chapter.URL)
				if err != nil {
					// 尝试重试
					fmt.Printf(
						"==> 正在重试下载失败章节内容: 【%s】，最大尝试次数: %d，失败原因：%s\n",
						chapter.Title,
						conf.Crawl.Max,
						err.Error(),
					)
					return false
				} else {
					downOk = true
				}
			}
			err = tools.ConvertChapter(chapter, conf.Base.Extname, b.rule)
			if err != nil {
				// 尝试重试
				fmt.Printf(
					"==> 正在重试转换失败章节: 【%s】，最大尝试次数: %d，失败原因：%s\n",
					chapter.Title,
					conf.Crawl.Max,
					err.Error(),
				)
				return false
			}
			return true
		},
		conf.Retry.MaxAttempts,
	)
	return nil
}

func (b *ChapterParser) crawl(url string) (string, error) {
	nextUrl := url
	sb := bytes.NewBufferString("")

	for {
		collector := getCollector(nil)
		collector.OnHTML(b.rule.Chapter.Content, func(e *colly.HTMLElement) {
			html, err := e.DOM.Html()
			if err == nil {
				sb.WriteString(html)
			} else {
				// 打印错误
				fmt.Printf("ChapterParser crawl Error parsing HTML: %v\n", err)
			}
		})
		if !b.rule.Chapter.Pagination {
			err := collector.Visit(nextUrl)
			if err != nil {
				return "", err
			}
			collector.Wait()
			return sb.String(), nil
		} else {
			collector.OnHTML(b.rule.Chapter.NextPage, func(e *colly.HTMLElement) {
				if strings.Contains(e.Text, "下一章") {
					nextUrl = ""
					return
				}
				href := e.Attr("href")
				nextUrl = utils.NormalizeURL(href, b.rule.URL)
			})
			err := collector.Visit(nextUrl)
			if err != nil {
				return "", err
			}
			collector.Wait()
		}
		if nextUrl == "" {
			break
		}
	}
	return sb.String(), nil
}
