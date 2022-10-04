package middleware

import (
	"context"
	"net/http"
	"server/pkg/code"
	"server/pkg/response"

	"time"

	"github.com/gin-gonic/gin"
	"github.com/gogf/gf/frame/g"
	"github.com/juju/ratelimit"
)

// 每次向桶中放令牌时，是放 quantum 个令牌
// fillInterval放令牌的间隔时间 capacity桶容量 quantum每次放令牌的数量
func RateLimitMiddleware(fillInterval time.Duration, capacity int64, quantum int64) gin.HandlerFunc {
	bucket := ratelimit.NewBucketWithQuantum(fillInterval, capacity, quantum)
	return func(c *gin.Context) {
		if bucket.TakeAvailable(1) < 1 {
			response.Error(c, http.StatusBadRequest, code.RateLimit, code.GetErrMsg(code.RateLimit), g.I18n().Translate(context.Background(), `{#ratelimit}`))
			c.Abort()
			return
		}
		c.Next()
	}
}

/**
限流：漏桶，令牌桶

漏桶：每隔一段时间释放一个令牌，拿到令牌的才可以继续请求，没拿到的等待先拿到请求令牌

令牌桶：匀速向桶中添加令牌，请求服务时先从桶中获取令牌，有令牌再去请求服务

区别：漏桶释放令牌是匀速的，用户等待令牌的下发；令牌桶中的令牌数量固定，用户主动去令牌桶中取令牌；比如10个用户来请求，漏桶会挨个给10个用户按照间隔时间匀速发放，但令牌桶的话，10个用户可以同时去桶里获取令牌，只有桶中数量有十个，就可以一次性全部通过，即并发；

令牌桶算法：系统以恒定的速率产生令牌，然后把令牌放到令牌桶中，令牌桶有一个容量，当令牌桶满了的时候，再向其中放令牌，那么多余的令牌会被丢弃；当想要处理一个请求的时候，需要从令牌桶中取出一个令牌，如果此时令牌桶中没有令牌，那么则拒绝该请求。
**/
