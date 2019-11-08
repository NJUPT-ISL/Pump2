package scheduler

import (
	"github.com/Mr-Linus/Pump2/rpc"
	"log"
)

type Node struct {
	IP       string
	Active   bool
	NodeStat rpc.NodeStat
	Score    int
}

var (
	Nodes   []Node
	IPs     []string
	workers = 10
)

func Schedule() (IP string, err error) {
	ActiveNode, err := FilterNodes(Nodes)
	if err != nil {
		log.Println(err)
		return "", err
	}
	//TODO
	return ActiveNode[0].IP, nil
}
