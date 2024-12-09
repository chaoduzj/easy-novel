package tools

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/767829413/easy-novel/internal/config"
	"github.com/767829413/easy-novel/internal/definition"
	"github.com/767829413/easy-novel/internal/model"
	"github.com/767829413/easy-novel/pkg/utils"
	"github.com/fatih/color"
	"github.com/go-resty/resty/v2"
	"github.com/go-shiori/go-epub"
)

func ProcessSaveHandler(book *model.Book, dirPath string) error {
	conf := config.GetConf()
	s := fmt.Sprintf("\n<== 《%s》（%s）下载完毕，", book.BookName, book.Author)

	switch conf.Base.Extname {
	case definition.NovelExtname_TXT, definition.NovelExtname_EPUB:
		s += fmt.Sprintf("正在合并为 %s", strings.ToUpper(conf.Base.Extname))
	case definition.NovelExtname_HTML:
		s += "正在生成 HTML 目录文件"
	}
	fmt.Println(s + " ...")

	switch conf.Base.Extname {
	case definition.NovelExtname_TXT:
		return txtMergeHandler(book, dirPath)
	case definition.NovelExtname_EPUB:
		return epubMergeHandler(book, dirPath)
	case definition.NovelExtname_HTML:
		return htmlMergeHandler(book, dirPath)
	default:
		return fmt.Errorf("unsupported extension: %s", conf.Base.Extname)
	}
}

func txtMergeHandler(book *model.Book, dirPath string) error {
	outputPath := filepath.Join(dirPath, fmt.Sprintf("%s（%s）.txt", book.BookName, book.Author))
	homePageFile, err := os.Create(outputPath)
	if err != nil {
		return err
	}
	defer homePageFile.Close()
	// 首页添加书籍信息
	bookInfo := []string{
		fmt.Sprintf("书名：%s", book.BookName),
		fmt.Sprintf("作者：%s", book.Author),
		fmt.Sprintf("简介：%s", book.Intro),
		strings.Repeat("\u3000", 2),
	}

	for _, line := range bookInfo {
		if _, err := homePageFile.WriteString(line + "\n"); err != nil {
			return err
		}
	}

	// 合并章节文件
	filePaths, err := utils.GetSortedFilePaths(dirPath)
	if err != nil {
		return err
	}
	for _, f := range filePaths {
		content, err := os.ReadFile(f)
		if err != nil {
			return err
		}
		if _, err := homePageFile.Write(content); err != nil {
			return err
		}
	}
	return nil
}
func epubMergeHandler(book *model.Book, dirPath string) error {
	var (
		err       error
		filePaths []string
		attempts  = 7
	)
	utils.SpinWaitMaxRetryAttempts(func() bool {
		filePaths, err = utils.GetSortedFilePaths(dirPath)
		if len(filePaths) == 0 || err != nil {
			return false
		}
		return true
	}, attempts)
	if len(filePaths) == 0 {
		utils.GetColorIns(color.FgRed).
			Printf("《%s》（%s）下载章节数为 0，取消生成 EPUB", book.BookName, book.Author)
		return nil
	}

	// 等待文件系统更新索引
	epubIns, err := epub.NewEpub(book.BookName)
	if err != nil {
		return err
	}
	epubIns.SetAuthor(book.Author)
	epubIns.SetDescription(book.Intro)
	epubIns.SetLang("zh")

	for _, filePath := range filePaths {
		content, err := os.ReadFile(filePath)
		if err != nil {
			return err
		}
		// 获取文件名
		fileName := filepath.Base(filePath)

		// 从文件名中提取标题
		title := strings.SplitN(fileName, "_", 2)[1]
		title = strings.TrimSuffix(title, filepath.Ext(title))
		_, err = epubIns.AddSection(string(content), title, filePath, "")
		if err != nil {
			return err
		}
	}

	// 下载封面
	if len(book.CoverURL) != 0 {
		fmt.Printf("<== 正在下载封面：%s", book.CoverURL)
		client := resty.New()
		resp, err := client.R().Get(book.CoverURL)
		if err != nil {
			utils.GetColorIns(color.FgRed).
				Printf("封面下载失败：%s\n", err.Error())
		} else {
			coverPath := filepath.Join(dirPath, "cover.jpg")
			if err := os.WriteFile(coverPath, resp.Body(), 0644); err != nil {
				utils.GetColorIns(color.FgRed).
					Printf("保存封面失败：%s", err)
			} else {
				epubIns.AddImage(coverPath, "cover.jpg")
				epubIns.SetCover("cover.jpg", "")
			}
		}
	}
	// 保存 EPUB 文件
	savePath := filepath.Join(filepath.Dir(dirPath), book.BookName+".epub")
	err = epubIns.Write(savePath)
	if err != nil {
		return err
	}
	// TODO: 处理 EPUB 格式的临时 HTML ，删除文件
	return nil
}
func htmlMergeHandler(book *model.Book, dirPath string) error {
	return nil
}
