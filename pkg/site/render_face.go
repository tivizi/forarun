package site

import (
	"encoding/hex"
	"fmt"
	"regexp"
	"strings"
)

var facereg *regexp.Regexp = regexp.MustCompile("\\{([\u4e00-\u9fa5]+)\\}")

// RenderFace 渲染表情
func RenderFace(text string) string {
	matchs := facereg.FindAllStringSubmatch(text, -1)
	for _, m := range matchs {
		text = strings.ReplaceAll(text, m[0], fmt.Sprintf("![face-%s](https://hu60.cn/img/face/%s)", m[1], hex.EncodeToString([]byte(m[1]))+".gif"))
	}
	return text
}
