package myContext

import (
	"TradeMatching/common/glog/log"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"path"
	"runtime"
)

var (
	ClientType = "client-type"
)

type Handler func(ctx *MyContext)

type MyContext struct {
	*gin.Context
	Trace string // 追蹤用的Request-Id
}

func Background(c *gin.Context) *MyContext {
	return &MyContext{c, c.GetHeader(viper.GetString("http.headerTrace"))}
}

func (c *MyContext) Info(args ...interface{}) {
	log.ZapLog.Named(funcName(c.Trace)).Info(args...)
}

func (c *MyContext) Infof(template string, args ...interface{}) {
	log.ZapLog.Named(funcName(c.Trace)).Infof(template, args...)
}

func (c *MyContext) Warn(args ...interface{}) {
	log.ZapLog.Named(funcName(c.Trace)).Warn(args...)
}

func (c *MyContext) Warnf(template string, args ...interface{}) {
	log.ZapLog.Named(funcName(c.Trace)).Warnf(template, args...)
}

func (c *MyContext) Error(args ...interface{}) {
	log.ZapLog.Named(funcName(c.Trace)).Error(args...)
}

func (c *MyContext) Errorf(template string, args ...interface{}) {
	log.ZapLog.Named(funcName(c.Trace)).Errorf(template, args...)
}

func funcName(trace string) string {
	pc, _, _, _ := runtime.Caller(2)
	funcName := runtime.FuncForPC(pc).Name()
	return path.Base(funcName) + " " + trace
}

func (c *MyContext) ClientType() string {
	return c.Request.Header.Get(ClientType)
}

func Wrap(h Handler) gin.HandlerFunc {
	return func(c *gin.Context) {
		ctxAny, exists := c.Get("myContext")
		if !exists {
			c.JSON(500, gin.H{"error": "MyContext missing"})
			return
		}
		ctx := ctxAny.(*MyContext)
		h(ctx)
	}
}
