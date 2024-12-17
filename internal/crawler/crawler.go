package crawler

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"sync"
	"sync/atomic"
	"time"

	"github.com/767829413/easy-novel/internal/config"
	"github.com/767829413/easy-novel/internal/model"
	"github.com/767829413/easy-novel/internal/parse"
	"github.com/767829413/easy-novel/internal/tools"
	"github.com/767829413/easy-novel/pkg/utils"
	"github.com/fatih/color"
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
	conf := config.GetConf()
	// Implement search logic using a search engine
	fmt.Println("<== 正在搜索...")
	startTime := time.Now()
	// 解析
	searchResults, err := parse.NewSearchResultParser(conf.Base.SourceID).Parse(key)
	if err != nil {
		fmt.Printf("<== 执行搜索索失败：%s\n", err.Error())
		return nil
	}
	duration := time.Since(startTime)
	fmt.Printf("<== 搜索到 %d 条记录，耗时 %f s\n", len(searchResults), duration.Seconds())
	return searchResults
}

func (nc *novelCrawler) Crawl(res *model.SearchResult, start, end int) *model.CrawlResult {
	conf := config.GetConf()
	// 小说详情页抓取解析
	book, err := parse.NewBookParser(conf.Base.SourceID).Parse(res.Url)
	if err != nil {
		utils.GetColorIns(color.BgRed).Printf("<== 执行详情页获取失败：%s\n", err.Error())
		return nil
	}

	// Format the directory name as "BookName (Author)"
	bookDir := fmt.Sprintf("%s (%s)", book.BookName, book.Author)
	dirPath := filepath.Join(conf.Base.DownloadPath, bookDir)
	// Create the directory
	err = os.MkdirAll(dirPath, os.ModePerm)
	if err != nil {
		utils.GetColorIns(color.BgRed).
			Printf("创建下载目录失败\n1. 检查下载路径是否合法\n2. 尝试以管理员身份运行（某些目录需要管理员权限）\n")
		return nil
	}
	fmt.Printf("<== 开始获取章节目录：%s\n", book.BookName)

	// 获取小说目录
	catalogsParser := parse.NewCatalogsParser(conf.Base.SourceID)
	catalogs, err := catalogsParser.Parse(res.Url, start, end)
	if err != nil {
		utils.GetColorIns(color.BgRed).Printf("<== 执行获取章节目录失败：%s\n", err.Error())
		return nil
	}
	if len(catalogs) == 0 {
		utils.GetColorIns(color.FgHiYellow).Println("获取章节内容为空")
		return nil
	}

	startTime := time.Now()
	// 解析下载内容
	var wg sync.WaitGroup
	cpuNum := runtime.NumCPU()
	threads := conf.Crawl.Threads
	if threads == -1 {
		threads = cpuNum * 2
	}
	semaphore := make(chan struct{}, threads)
	var nowCatalogsCount int32 = int32(len(catalogs))

	fmt.Printf(
		"<== 开始下载《%s》（%s） 共计 %d 章 | 协程数：%d",
		book.BookName,
		book.Author,
		len(catalogs),
		threads,
	)

	for _, chapter := range catalogs {
		wg.Add(1)
		go func(chapter *model.Chapter, bookDir string) {
			defer wg.Done()
			semaphore <- struct{}{}
			defer func() { <-semaphore }()
			// 下载逻辑
			fmt.Printf("<== 正在下载: 【%s】\n", chapter.Title)
			parse.NewChapterParser(conf.Base.SourceID).Parse(chapter, res, book, bookDir)
			tools.CreateFileForChapter(chapter, bookDir)
			atomic.AddInt32(&nowCatalogsCount, -1)
			fmt.Printf("<== 下载结束,待下载章节数：%d\n", atomic.LoadInt32(&nowCatalogsCount))
		}(
			chapter,
			bookDir,
		)
	}
	wg.Wait()

	// 合并生成小说文件格式
	err = tools.ProcessSaveHandler(book, dirPath)
	if err != nil {
		utils.GetColorIns(color.BgRed).Printf("<== 执行合并生成小说文件格式失败：%s\n", err.Error())
		return nil
	}
	return &model.CrawlResult{TakeTime: int64(time.Since(startTime).Seconds())}
}
