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
	s := fmt.Sprintf(
		"\n<== 《%s》（%s）下载完毕，正在合并为 %s ...\n",
		book.BookName,
		book.Author,
		strings.ToUpper(conf.Base.Extname),
	)
	fmt.Println(s)
	switch conf.Base.Extname {
	case definition.NovelExtname_TXT:
		return txtMergeHandler(book, dirPath)
	case definition.NovelExtname_EPUB:
		return epubMergeHandler(book, dirPath)
	default:
		return fmt.Errorf("unsupported extension: %s", conf.Base.Extname)
	}
}
