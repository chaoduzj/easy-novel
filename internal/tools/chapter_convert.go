package tools

import (
	"bytes"
	"fmt"
	"regexp"
	"strings"
	"sync"
	"text/template"

	"github.com/767829413/easy-novel/internal/definition"
	"github.com/767829413/easy-novel/internal/model"
)

var templates sync.Map

func ConvertChapter(
	chapter *model.Chapter,
	extName string,
	rule *model.Rule,
) error {
	var content string
	var err error
	content = FormatForChapter(FilterForChapter(chapter, rule), rule)

	switch extName {
	case definition.NovelExtname_TXT:
		content = txtConvert(chapter.Title, content)
	case definition.NovelExtname_EPUB, definition.NovelExtname_HTML:
		content, err = templateConvert(chapter.Title, content, extName)
		if err != nil {
			return err
		}
	}
	chapter.Content = content
	return nil
}

func txtConvert(title, content string) string {
	// 全角空格，用于首行缩进
	indent := strings.Repeat("\u3000", 2)
	// 创建正则表达式
	re := regexp.MustCompile(`<p>(.*?)</p>`)
	// 使用strings.Builder来高效构建字符串
	var result strings.Builder
	// 查找所有匹配项
	matches := re.FindAllStringSubmatch(content, -1)
	for _, match := range matches {
		result.WriteString(indent)
		result.WriteString(match[1])
		result.WriteString("\n")
	}
	// 构建最终内容
	finalContent := fmt.Sprintf("%s\n\n%s", title, result.String())

	return finalContent
}

func templateConvert(title, content, extName string) (string, error) {
	tmpl, err := getTemplate(extName)
	if err != nil {
		return "", err
	}

	data := struct {
		Title   string
		Content string
	}{
		Title:   title,
		Content: content,
	}

	var renderedContent bytes.Buffer
	err = tmpl.Execute(&renderedContent, data)
	if err != nil {
		return "", err
	}

	return renderedContent.String(), nil
}

func getTemplate(extName string) (*template.Template, error) {
	if tmpl, ok := templates.Load(extName); ok {
		return tmpl.(*template.Template), nil
	}

	var templateContent string
	switch extName {
	case definition.NovelExtname_HTML: // 通用HTML模板
		templateContent = definition.NovelTemp_HTML
	case definition.NovelExtname_EPUB:
		templateContent = definition.NovelTemp_EPUB
	default:
		return nil, fmt.Errorf("unsupported template type: %s", extName)
	}

	tmpl, err := template.New(extName).Parse(templateContent)
	if err != nil {
		return nil, err
	}

	actual, loaded := templates.LoadOrStore(extName, tmpl)
	if loaded {
		// Another goroutine has already stored a value, use that
		return actual.(*template.Template), nil
	}

	return tmpl, nil
}
