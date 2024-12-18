package chapter

import (
	"fmt"
	"regexp"
	"strings"
)

func txtConvert(title, content string) string {
	// 全角空格，用于首行缩进
	ident := strings.Repeat("\u3000", 2)
	re := regexp.MustCompile(`<p>(.*?)</p>`)
	matches := re.FindAllStringSubmatch(content, -1)

	var result strings.Builder
	for _, match := range matches {
		if len(match) > 1 {
			result.WriteString(ident)
			result.WriteString(match[1])
			result.WriteString("\n")
		}
	}

	return fmt.Sprintf("%s\n\n%s", title, result.String())
}
