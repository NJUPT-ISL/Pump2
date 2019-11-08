package scheduler

import e "errors"

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

//func CalculateNodePerform(n Node) (score int, err error) {
//
//}
