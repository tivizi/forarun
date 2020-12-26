package site

import (
	"fmt"
	"strings"
	"time"

	"github.com/tivizi/forarun/pkg/domain"
)

// IntelliTime 智能显示时间
func IntelliTime(t time.Time) string {
	return fmtDuration(time.Now().Sub(t))
}

func fmtDuration(d time.Duration) string {
	h := int64(d / time.Hour)
	m := int64(d / time.Minute)
	day := int64(h / 24)
	if day != 0 {
		if day == 1 {
			return "昨天"
		}
		if day == 2 {
			return "前天"
		}
		return fmt.Sprintf("%d天前", day)
	}
	if h != 0 {
		return fmt.Sprintf("%d小时前", h)
	}
	if m != 0 {
		return fmt.Sprintf("%d分钟前", m)
	}
	return "刚刚"
}

// ThreadContent 帖子内容
func ThreadContent(text string) string {
	return strings.ReplaceAll(text, "\n", "<br />")
}

// FromID SiteFromID
func FromID(siteid string) *domain.Site {
	return siteUniqueCache[siteid]
}
