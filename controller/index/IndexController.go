package index

import (
	"TradeMatching/common/myContext"
	"github.com/gin-gonic/gin"
	"net/http"
)

func GetById(ctx *myContext.MyContext) {
	ctx.JSON(http.StatusOK, gin.H{
		"data": "test",
	})
}
