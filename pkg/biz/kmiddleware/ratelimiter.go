// @Description 限流中间件

package kmiddleware

import (
	"github.com/LuoHongLiang0921/kuaigo/pkg/biz/tcontext"
	"github.com/LuoHongLiang0921/kuaigo/pkg/core/server/ginserver"
	"github.com/LuoHongLiang0921/kuaigo/pkg/ecode"
	"github.com/LuoHongLiang0921/kuaigo/pkg/errs"
)

// RateLimiter
//  @Description: RateLimiter 接口限流
//  @Receiver m
//  @Param rater
//  @Param udid 请求方唯一值 如udid等
//  @Return xgin.HandlerFunc
func (m *Middleware) RateLimiter(keys ...string) ginserver.HandlerFunc {
	return func(c *ginserver.TContext) {
		path := c.Request.URL.Path
		var udid string
		// todo 待优化
		if len(keys) == 0 && tcontext.GetPParam(c) != nil {
			udid = tcontext.GetPParam(c).DID
		} else if len(keys) > 0 {
			udid = keys[0]
		} else {
			udid = c.ClientIP()
		}

		ans := m.rater.Find(path)
		for _, an := range ans {
			if m.rater != nil {
				if m.rater.Execute(c, an, udid, path, c.ClientIP()).RateLimiter() {
					m.abortWithErrorJSON(c, errs.NewCustomError(ecode.CodeRateLimitError, "RateLimiter limit"))
				}
			}
		}
		c.Next()
	}
}
