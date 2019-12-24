package server

import (
	"github.com/Mr-Linus/Pump2/pkg/scheduler"
	"github.com/gin-gonic/gin"
)

func RunScheduler(File string){
	scheduler.LogPrint("Start the Pump2 Scheduler.")
	scheduler.LogPrint("Init Scheduler Cache.")
	if err := scheduler.InitCache(File, scheduler.Workers);err != nil{
		scheduler.LogErrPrint(err)
	}
	gin.DisableConsoleColor()
	Addr := ":5021"
	scheduler.LogPrint("Pump2 Scheduler is running at" + Addr)
	r := InitRouter()
	_ = r.Run(Addr) // listen and serve on 0.0.0.0:5021
}