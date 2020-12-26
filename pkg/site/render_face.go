package site

import (
	"fmt"
	"regexp"
	"strings"
)

var faces = map[string]string{
	"不开心": "bkx.gif",
	"冷":   "leng.gif",
	"乖":   "guai.gif",
	"勉强":  "mq.gif",
	"吐舌":  "ts.gif",
	"吐":   "tu.gif",
	"呵呵":  "hh.gif",
	"呼":   "hu.gif",
	"咦":   "yi.gif",
	"哈哈":  "haha.gif",
	"啊":   "a.gif",
	"喷":   "pen.gif",
	"太开心": "tkx.gif",
	"委屈":  "wq.gif",
	"怒":   "nu.gif",
	"惊哭":  "jk.gif",
	"惊讶":  "jy.gif",
	"汗":   "han.gif",
	"泪":   "lei.gif",
	"滑稽":  "hj.gif",
	"狂汗":  "kh.gif",
	"疑问":  "yw.gif",
	"真棒":  "zb.gif",
	"睡觉":  "sj.gif",
	"笑眼":  "xy.gif",
	"花心":  "hx.gif",
	"鄙视":  "bs.gif",
	"酷":   "ku.gif",
	"钱":   "qian.gif",
	"阴险":  "yx.gif",
	"黑线":  "heix.gif",
}
var facereg *regexp.Regexp = regexp.MustCompile("\\{([\u4e00-\u9fa5]+)\\}")

// RenderFace 渲染表情
func RenderFace(text string) string {
	matchs := facereg.FindAllStringSubmatch(text, -1)
	for _, m := range matchs {
		text = strings.ReplaceAll(text, m[0], fmt.Sprintf("![face-%s](/common/static/face/%s)", m[1], faces[m[1]]))
	}
	return text
}

// GetFaces 获取所有表情
func GetFaces() *map[string]string {
	return &faces
}
