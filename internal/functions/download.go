package functions

import (
	"fmt"
	"math"
	"os"
	"strconv"
	"strings"

	"github.com/767829413/easy-novel/internal/crawler"
	"github.com/767829413/easy-novel/internal/definition"
	"github.com/767829413/easy-novel/internal/model"
	"github.com/767829413/easy-novel/pkg/utils"
	"github.com/chzyer/readline"
	"github.com/fatih/color"
	"github.com/olekukonko/tablewriter"
	"github.com/sirupsen/logrus"
)

type download struct {
	log *logrus.Logger
}

func NewDownload(l *logrus.Logger) App {
	return &download{log: l}
}

func (d *download) Execute() error {
	d.log.Info("Starting novel download")
	utils.GetColorIns(color.BgGreen).Println("开始下载小说...")
	// 实现下载小说的逻辑
	rl, err := readline.New("")
	if err != nil {
		return err
	}
	defer rl.Close()

	// 1. Query
	prompt := utils.GetColorIns(color.BgBlue).Sprint("==> 请输入书名或作者（宁少字别错字）: ")
	rl.SetPrompt(prompt)

	keyword, err := rl.Readline()
	if err != nil {
		return err
	}
	keyword = strings.TrimSpace(keyword)
	if keyword == "" {
		return nil
	}

	crawler := crawler.NewNovelCrawler()
	results := crawler.Search(keyword)

	if len(results) == 0 {
		return nil
	}

	// 2. Print search results
	d.printSearchResult(results)

	// 3. Select download chapter
	for {
		rl.SetPrompt("==> 请输入下载序号（首列的数字，或输入 0 返回）：")
		input, err := rl.Readline()
		if err != nil {
			return err
		}
		input = strings.TrimSpace(input)

		num, err := strconv.Atoi(input)
		if err != nil {
			continue
		}

		if num == 0 {
			return nil
		}
		if num < 0 || num > len(results) {
			continue
		}

		sr := results[num-1]
		fmt.Printf("<== 你选择了《%s》(%s)\n", sr.BookName, sr.Author)
		fmt.Printf("==> %d: 重新选择功能", definition.ActionDownload_REOPEN)
		fmt.Printf("==> %d: 开始下载当前全本", definition.ActionDownload_START)
		fmt.Printf("==> %d: 重新输入序号", definition.ActionDownload_RESELECT)

		rl.SetPrompt("==> 请输入数字：")
		actionInput, err := rl.Readline()
		if err != nil {
			return err
		}

		action, err := strconv.Atoi(strings.TrimSpace(actionInput))
		if err != nil {
			continue
		}

		switch action {
		case definition.ActionDownload_REOPEN:
			return nil
		case definition.ActionDownload_START:
			start, end := 1, math.MaxInt // Max int
			res := crawler.Crawl(sr, start, end)
			fmt.Printf("<== 完成！总耗时 %d s\n", res.TakeTime)
			return nil
		case definition.ActionDownload_RESELECT:
			continue
		default:
			utils.GetColorIns(color.FgHiRed).Println("无效的选项，请重新输入")
		}
	}
}

func (d *download) printSearchResult(results []*model.SearchResult) {

	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"序号", "书名", "作者", "最新章节", "最后更新时间"})

	for i, r := range results {
		table.Append([]string{
			strconv.Itoa(i + 1),
			r.BookName,
			r.Author,
			r.LatestChapter,
			r.LatestUpdate,
		})
	}

	table.SetBorder(true) // 启用边框
	table.SetCenterSeparator("|")
	table.SetColumnSeparator("|")
	table.SetRowSeparator("-")
	table.SetHeaderLine(true)
	table.SetTablePadding("\t")                  // tab-separated columns
	table.SetAlignment(tablewriter.ALIGN_CENTER) // 设置内容居中

	// 启用每行内容上下的边框
	table.SetRowLine(true)

	table.Render()
}
