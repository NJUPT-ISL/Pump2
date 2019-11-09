package scheduler

import (
	rpc "github.com/Mr-Linus/Pump2/rpc"
	"log"
	"sync"
)

type Node struct {
	IP       string
	Active   bool
	NodeStat rpc.NodeStat
}

var (
	Nodes   []Node
	IPs     []string
	workers = 10
)

func Schedule() (IP string, err error) {
	var (
		maxScore int = 0
		nScore int = 0
	)
	if UpdateCache() != nil{
		log.Println(err)
		return "", err
	}
	activeNodes, err := FilterNodes(Nodes)
	if err != nil {
		log.Println(err)
		return "", err
	}
	cpuMax, err :=  GetMaxCpu(activeNodes)
	if err != nil {
		log.Println(err)
		return "", err
	}
	cpuFreqMax, err := GetMaxCpuFreq(activeNodes)
	if err != nil {
		log.Println(err)
		return "", err
	}
	freeMemMax, err := GetMaxFreeMemory(activeNodes)
	if err != nil {
		log.Println(err)
		return "", err
	}
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
					nScore,err = CalculateNodePerform(n,cpuMax,cpuFreqMax,freeMemMax)
					if err != nil {
						log.Println(err)
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
	return IP,nil
}
