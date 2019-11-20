package scheduler

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"time"
)

func PostAddTask(c *gin.Context) {
	if err := UpdateCache(workers);err != nil {
		LogErrPrint(err)
		c.JSON(200, gin.H{
			"state": "failed",
			"error": err.Error(),
		})
		return
	}
	address,err := DefaultSchedule()
	if err != nil {
		LogErrPrint(err)
		c.JSON(200, gin.H{
			"state": "failed",
			"error": err.Error(),
		})
		return
	}
	task,err := CheckTask(c)
	if err != nil {
		LogErrPrint(err)
		c.JSON(200, gin.H{
			"state": "failed",
			"error": err.Error(),
		})
		return
	}
	state,err := DoTask(address,task)
	if err != nil{
		LogErrPrint(err)
		c.JSON(200, gin.H{
			"state": "failed",
			"error": err.Error(),
		})
		return
	}
	Tasks = append(Tasks,Task{WorkNode:address,task:task,state:state})
	c.JSON(200, gin.H{
		"state": "ok",
	})
}

func GetListTask(c *gin.Context) {
	c.JSON(200, Tasks)
}

func InitRouter() *gin.Engine {
	router := gin.New()
	//router.Use(gin.Logger())
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
		buildGroup.GET("list",GetListTask)
	}
	return router
}

func RunScheduler(File string,){
	LogPrint("Start the Pump2 Scheduler.")
	LogPrint("Init Scheduler Cache.")
	if err := InitCache(File,workers);err != nil{
		LogErrPrint(err)
	}
	gin.DisableConsoleColor()
	Addr := ":60886"
	LogPrint("Pump2 Scheduler is running at" + Addr)
	r := InitRouter()
	_ = r.RunTLS(Addr, "pem/server.crt", "pem/server.key") // listen and serve on 0.0.0.0:8081
}