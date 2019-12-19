package server

import (
	"github.com/Mr-Linus/Pump2/pkg/scheduler"
	"github.com/gin-gonic/gin"
	uuid "github.com/satori/go.uuid"
)

func PostAddTask(c *gin.Context) {
	if err := scheduler.UpdateCache(scheduler.Workers);err != nil {
		scheduler.LogErrPrint(err)
		c.JSON(200, gin.H{
			"state": "failed",
			"error": err.Error(),
		})
		return
	}
	address,err := scheduler.DefaultSchedule()
	if err != nil {
		scheduler.LogErrPrint(err)
		c.JSON(200, gin.H{
			"state": "failed",
			"error": err.Error(),
		})
		return
	}
	task,err := scheduler.CheckTask(c)
	if err != nil {
		scheduler.LogErrPrint(err)
		c.JSON(200, gin.H{
			"state": "failed",
			"error": err.Error(),
		})
		return
	}
	t := scheduler.Task{
		ID:       uuid.NewV4().String(),
		WorkerNode: address,
		BuildInfo:     task,
		IsBuild:  true,
	}
	scheduler.Tasks = append(scheduler.Tasks,t)
	go func() {
		state,err := scheduler.DoTask(address,task)
		for i,task := range scheduler.Tasks {
			if task.BuildInfo.Name == t.BuildInfo.Name{
				scheduler.Tasks = append(scheduler.Tasks[:i], scheduler.Tasks[i+1:]...)
				t.IsBuild = false
				t.State = state
				scheduler.Tasks = append(scheduler.Tasks,t)
				break
			}
		}
		if err != nil{
			scheduler.LogErrPrint(err)
			return
		}
	}()
	c.JSON(200, gin.H{
		"state": "ok",
		"TaskID":t.ID,
	})
}

func GetListTask(c *gin.Context) {
	c.JSON(200, scheduler.Tasks)
}

func PostTaskInfo(c *gin.Context) {
	for _,t := range scheduler.Tasks{
		if t.ID == c.PostForm("uuid"){
			c.JSON(200,t)
			return
		}
	}
	c.JSON(403,gin.H{
		"state":"failed",
		"error":"Can't find the task ID.",
	})
}

func GetListNode(c *gin.Context) {
	c.JSON(200, scheduler.Nodes)
}