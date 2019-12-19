package server

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"time"
)

func InitRouter() *gin.Engine {
	router := gin.New()
	router.Use(gin.LoggerWithFormatter(func(param gin.LogFormatterParams) string {
		return fmt.Sprintf("[Pump2 Scheduler]%s - [%s] \"%s %s %s %d %s \"%s\" %s\"\n",
			param.ClientIP,
			time.Now(),
			param.Method,
			param.Path,
			param.Request.Proto,
			param.StatusCode,
			param.Latency,
			param.Request.UserAgent(),
			param.ErrorMessage,
		)
	}))
	router.Use(gin.Recovery())
	buildGroup := router.Group("/build")
	{
		buildGroup.POST("add", PostAddTask)
		buildGroup.GET("list", GetListTask)
		buildGroup.POST("info", PostTaskInfo)
	}
	nodeGroup := router.Group("/node")
	{
		nodeGroup.GET("list", GetListNode)
	}
	return router
}