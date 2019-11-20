package scheduler

import (
	"fmt"
	"time"
)

func LogErrPrint(err error) {
	fmt.Printf("[Pump2 Scheduler] - [%s] \"Error: %s\"\n", time.Now(), err)
}

func LogPrint(log string) {
	fmt.Printf("[Pump2 Scheduler] - [%s] \" %s\"\n", time.Now(), log)
}