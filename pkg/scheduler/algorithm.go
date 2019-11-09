package scheduler

import (
	e "errors"
	"github.com/pkg/errors"
	"log"
)

func FilterNodes(n []Node) (nodes []Node, err error) {
	nodes = []Node{}
	for _, no := range Nodes {
		if no.Active == true {
			nodes = append(nodes, no)
		}
	}
	if len(nodes) == 0 {
		return nodes, e.New("The number of Active Node is " + string(len(nodes)))
	}
	return nodes, nil
}

func CalculateNodePerform(n Node, cpuMax int32, freqMax float32, freeMemMax int32) (score int, err error) {
	sCPU, err := CalculateCPUScore(n, cpuMax)
	if err != nil {
		log.Println(err)
		return 0, err
	}
	sCPUFreq, err := CalculateCPUFreqScore(n, freqMax)
	if err != nil {
		log.Println(err)
		return 0, err
	}
	sLoad, err := CalculateLoadAvgScore(n)
	if err != nil {
		log.Println(err)
		return 0, err
	}
	sMem, err := CalculateMemoryScore(n)
	sMemFree, err := CalculateMemoryFreeScore(n,freeMemMax)
	score = ((sCPU * 5 + sCPUFreq * 5) / 10 - sLoad) * 3 + (sMem * 3 + sMemFree * 7) * 7
	return score, nil
}

func CalculateCPUScore(n Node, cpuMax int32) (score int, err error) {
	if n.NodeStat.Cpu == 0 {
		return 0, errors.New(n.IP + ":The Number of the CPU is 0")
	}
	return int(float32(n.NodeStat.Cpu) / float32(cpuMax) * 100), nil
}

func CalculateCPUFreqScore(n Node, freqMax float32) (score int, err error) {
	if n.NodeStat.CpuFreq == 0 {
		return 0, errors.New("The CPU Frequency of the node" + n.IP + "is 0")
	}
	return int(n.NodeStat.CpuFreq / freqMax * 100), nil
}

func CalculateLoadAvgScore(n Node) (score int, err error) {
	if n.NodeStat.LoadAvg == 0 {
		return 0, errors.New("The Load of the Node is 0")
	}
	return int(n.NodeStat.LoadAvg / float32(n.NodeStat.Cpu) * 100), nil
}

func CalculateMemoryScore(n Node) (score int, err error) {
	if n.NodeStat.Memory == 0 {
		return 0, errors.New("The Memory of the Node is 0")
	}
	return int(n.NodeStat.MemoryFree / n.NodeStat.Memory), nil
}

func CalculateMemoryFreeScore(n Node, freeMemMax int32) (score int, err error) {
	if n.NodeStat.MemoryFree == 0 {
		return 0, errors.New("The Free Memory of the Node is 0")
	}
	return int(n.NodeStat.MemoryFree / freeMemMax), nil
}
