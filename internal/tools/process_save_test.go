package tools

import (
	"encoding/json"
	"testing"

	"github.com/767829413/easy-novel/internal/model"
)

var (
	bookJsonStr = ``
	dirPath     = ``
)

func TestEpubMergeHandler(t *testing.T) {
	book := model.Book{}
	json.Unmarshal([]byte(bookJsonStr), &book)
	err := epubMergeHandler(&book, dirPath)
	if err == nil {
		t.Fatal("expected an error, but got nil")
	}
}
