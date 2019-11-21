package scheduler

import (
	rpc "github.com/Mr-Linus/Pump2/rpc"
	"strconv"
	"sync"
)

type Node struct {
	IP       string
	Active   bool
	NodeStat rpc.NodeStat
}

type Task struct {
	workNode string
	task rpc.BuildInfo
	isBuild bool
	state bool
}

var (
	Nodes   []Node
	Tasks   []Task
	IPs     []string
	workers = 10
)


func CalculateHighestScore(activeNodes []Node, workers int) (string, int, error){
	var (
		nScore int = 0
		maxScore = 0
		IP = ""
	)
	cpuMax, cpuFreqMax, freeMemMax,buildNumMax, err := CollectNodeInfo(activeNodes)
	var stop <-chan struct{}
	pieces := len(activeNodes)
	toProcess := make(chan Node, pieces)
	for _, node := range activeNodes {
		toProcess <- node
	}
	close(toProcess)
	if pieces < workers {
		workers = pieces
	}
	wg := sync.WaitGroup{}
	wg.Add(workers)
	for i := 0; i < workers; i++ {
		go func() {
			defer wg.Done()
			for n := range toProcess {
				select {
				case <-stop:
					return
				default:
					nScore,err = CalculateNodePerform(n,cpuMax,cpuFreqMax,freeMemMax,buildNumMax)
					if err != nil {
						LogErrPrint(err)
						return
					}
					if maxScore < nScore {
						maxScore = nScore
						IP = n.IP
					}
				}
			}
		}()
	}
	wg.Wait()
	return IP, maxScore, err
}

func DefaultSchedule() (IP string, err error) {
	if UpdateCache(workers) != nil{
		LogErrPrint(err)
		return "", err
	}
	activeNodes, err := FilterNodes(Nodes)
	if err != nil {
		LogErrPrint(err)
		return "", err
	}
	IP, maxScore, err := CalculateHighestScore(activeNodes,workers)
	LogPrint("The BestNode is "+IP+" ,Score: "+strconv.Itoa(maxScore))
	return IP,nil
}

