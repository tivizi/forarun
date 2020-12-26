package daemon

import (
	"github.com/tivizi/forarun/pkg/domain"
	"github.com/tivizi/forarun/pkg/site"
)

func init() {
	go func() {
		for {
			ctx := <-*(site.OnlineStateChan())
			if ctx.Site != nil {
				domain.NewOnlineActive(ctx.Site.ID, ctx.UserAgent, ctx.ClientIP, ctx.Session)
			}
		}
	}()
}
