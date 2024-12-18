package merge

import (
	"fmt"
	"strings"

	"github.com/767829413/easy-novel/internal/config"
	"github.com/767829413/easy-novel/internal/definition"
	"github.com/767829413/easy-novel/internal/model"
)

func MergeSaveHandler(book *model.Book, dirPath string) error {
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
