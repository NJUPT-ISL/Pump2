package scheduler

import (
	"context"
	e "errors"
	"github.com/Mr-Linus/Pump2/pkg/yaml"
	"github.com/Mr-Linus/Pump2/rpc"
	"github.com/pkg/errors"
	"google.golang.org/grpc"
	"log"
	"sync"
)

func AddNodeInfo(IP string) error {
	conn, err := grpc.Dial(IP, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
		return err
	}
	defer conn.Close()
	c := rpc.NewPump2Client(conn)
	// Contact the server and print out its response.
	nodeStat, err := c.NodeStatus(context.Background(), &rpc.NodeInfo{})
	if err != nil {
		log.Println(err)
		return err
	}
	n := Node{IP: IP, Active: true, NodeStat: *nodeStat}
	Nodes = append(Nodes, n)
	return nil
}

func UpdateNodeInfo(IP string) error {
	var index int
	for i, n := range Nodes {
		if n.IP == IP {
			index = i
			break
		}
	}
	conn, err := grpc.Dial(IP, grpc.WithInsecure())
	if err != nil {
		Nodes[index].Active = false
		log.Fatalf("did not connect: %v", err)
		return err
	}
	defer conn.Close()
	c := rpc.NewPump2Client(conn)
	nodeStat, err := c.NodeStatus(context.Background(), &rpc.NodeInfo{})
	if err != nil {
		log.Println(err)
		return err
	}
	Nodes[index].NodeStat = *nodeStat
	return nil
}

func GetNodesIP(File string) error {
	conf, err := yaml.ReadNodeYaml(File)
	if err != nil {
		log.Println(err)
		return err
	}
	IPs = conf.Nodes.IP
	return nil
}

func InitCache(File string) error {
	if err := GetNodesIP(File); err != nil {
		log.Println(err)
		return err
	}
	var stop <-chan struct{}
	pieces := len(IPs)
	toProcess := make(chan string, pieces)
	for _, IP := range IPs {
		toProcess <- IP
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
			for ip := range toProcess {
				select {
				case <-stop:
					return
				default:
					if err := AddNodeInfo(ip); err != nil {
						log.Println(err)
					}
				}
			}
		}()
	}
	return nil
}

func UpdateCache() error {
	if len(IPs) == 0 {
		return e.New("Error: the Node List is:" + string(len(IPs)))
	}
	var stop <-chan struct{}
	pieces := len(IPs)
	toProcess := make(chan string, pieces)
	for _, IP := range IPs {
		toProcess <- IP
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
			for ip := range toProcess {
				select {
				case <-stop:
					return
				default:
					if err := UpdateNodeInfo(ip); err != nil {
						log.Println(err)
					}
				}
			}
		}()
	}
	return nil
}

func GetMaxCpu(ns []Node) (cpu int32, err error) {
	cpu = 0
	for _, n := range ns {
		if n.NodeStat.Cpu > cpu {
			cpu = n.NodeStat.Cpu
		}
	}
	if cpu == 0 || len(ns) == 0 {
		return 0, errors.New(
			"The number of node is " + string(len(ns)) + " and the Max CPU is " + string(cpu))
	}
	return cpu, nil
}

func GetMaxCpuFreq(ns []Node) (cpuFreq float32, err error) {
	cpuFreq = 0
	for _, n := range ns {
		if n.NodeStat.CpuFreq > cpuFreq {
			cpuFreq = n.NodeStat.CpuFreq
		}
	}
	if cpuFreq == 0 || len(ns) == 0 {
		return 0, errors.New(
			"The number of node is " + string(len(ns)) + " and the Max CPUFreq is " + string(int(cpuFreq)))
	}
	return cpuFreq, nil
}
