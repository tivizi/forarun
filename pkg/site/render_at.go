package site

import (
	"fmt"
	"regexp"
	"strings"
)

var atreg *regexp.Regexp = regexp.MustCompile("@(.+?)([\\s,])")

// RenderAt 渲染表情
func RenderAt(text string) string {
	matchs := atreg.FindAllStringSubmatch(text, -1)
	for _, m := range matchs {
		text = strings.ReplaceAll(text, m[0], fmt.Sprintf("@[%s](/profile/%s)%s", m[1], m[1], m[2]))
	}
	return text
}
