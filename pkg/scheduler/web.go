package scheduler

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/satori/go.uuid"
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
	t := Task{
		ID:       uuid.NewV4().String(),
		WorkerNode: address,
		BuildInfo:     task,
		IsBuild:  true,
	}
	Tasks = append(Tasks,t)
	go func() {
		state,err := DoTask(address,task)
		for i,task := range Tasks{
			if task.BuildInfo.Name == t.BuildInfo.Name{
				Tasks = append(Tasks[:i],Tasks[i+1:]...)
				t.IsBuild = false
				t.State = state
				Tasks = append(Tasks,t)
				break
			}
		}
		if err != nil{
			LogErrPrint(err)
			return
		}
	}()
	c.JSON(200, gin.H{
		"state": "ok",
		"TaskID":t.ID,
	})
}

func GetListTask(c *gin.Context) {
	c.JSON(200,Tasks)
}

func GetListNode(c *gin.Context) {
	c.JSON(200,Nodes)
}

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
		buildGroup.GET("list",GetListTask)
	}
	nodeGroup := router.Group("/node")
	{
		nodeGroup.GET("list",GetListNode)
	}
	return router
}

func RunScheduler(File string){
	LogPrint("Start the Pump2 Scheduler.")
	LogPrint("Init Scheduler Cache.")
	if err := InitCache(File,workers);err != nil{
		LogErrPrint(err)
	}
	gin.DisableConsoleColor()
	Addr := ":5021"
	LogPrint("Pump2 Scheduler is running at" + Addr)
	r := InitRouter()
	_ = r.Run(Addr) // listen and serve on 0.0.0.0:8081
}