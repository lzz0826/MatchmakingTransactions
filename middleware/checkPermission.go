package middleware

import (
	"TradeMatching/common/myContext"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/spf13/viper"
)

func TraceMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		traceHeader := viper.GetString("http.headerTrace") // 例如 "X-Request-Id"
		traceId := c.GetHeader(traceHeader)
		if traceId == "" {
			traceId = uuid.NewString() // 沒有就產生一個 UUID
		}

		// Response Header 中也上這個 trace id
		c.Header(traceHeader, traceId)

		// 包裝成 MyContext，設置 traceId
		ctx := &myContext.MyContext{
			Context: c,
			Trace:   traceId,
		}

		// 設定到 gin.Context 的鍵中，方便 controller 拿出來用
		c.Set("myContext", ctx)

		// 繼續後面的 middleware/handler
		c.Next()
	}
}
