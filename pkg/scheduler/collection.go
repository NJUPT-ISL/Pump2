package scheduler

import (
	"context"
	e "errors"
	"github.com/Mr-Linus/Pump2/pkg/yaml"
	rpc "github.com/Mr-Linus/Pump2/rpc"
	"github.com/pkg/errors"
	"google.golang.org/grpc"
	"sync"
)

func AddNodeInfo(IP string) error {
	conn, err := grpc.Dial(IP, grpc.WithInsecure())
	if err != nil {
		LogErrPrint(err)
		return err
	}
	defer conn.Close()
	c := rpc.NewPump2Client(conn)
	// Contact the server and print out its response.
	nodeStat, err := c.NodeStatus(context.Background(), &rpc.NodeInfo{})
	if err != nil {
		LogErrPrint(err)
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
		LogErrPrint(err)
		return err
	}
	defer conn.Close()
	c := rpc.NewPump2Client(conn)
	nodeStat, err := c.NodeStatus(context.Background(), &rpc.NodeInfo{})
	if err != nil {
		LogErrPrint(err)
		return err
	}
	Nodes[index].NodeStat = *nodeStat
	return nil
}

func GetNodesIP(File string) error {
	conf, err := yaml.ReadNodeYaml(File)
	if err != nil {
		LogErrPrint(err)
		return err
	}
	IPs = conf.Nodes.IP
	return nil
}

func InitCache(File string,workers int) error {
	if err := GetNodesIP(File); err != nil {
		LogErrPrint(err)
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
						LogErrPrint(err)
					}
				}
			}
		}()
	}
	wg.Wait()
	return nil
}

func UpdateCache(workers int) error {
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
						LogErrPrint(err)
					}
				}
			}
		}()
	}
	wg.Wait()
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

func GetMaxFreeMemory(ns []Node) (Mem int32, err error) {
	Mem = 0
	for _, n := range ns {
		if n.NodeStat.MemoryFree > Mem {
			Mem = n.NodeStat.MemoryFree
		}
	}
	if Mem == 0 || len(ns) == 0 {
		return 0, errors.New(
			"The number of node is " + string(len(ns)) + " and the Max Free Memory is " + string(int(Mem)))
	}
	return Mem, nil
}

func GetMaxBuildNum(ns []Node) (BuildNum int32, err error) {
	BuildNum = 0
	for _, n := range ns {
		if n.NodeStat.BuildNum > BuildNum {
			BuildNum = n.NodeStat.BuildNum
		}
	}
	if BuildNum == 0 || len(ns) == 0 {
		return 0, errors.New(
			"The number of node is " + string(len(ns)) + " and the Max Build Number is " + string(int(BuildNum)))
	}
	return BuildNum, nil
}

func CollectNodeInfo (activeNodes []Node) (cpuMax int32,cpuFreqMax float32,freeMemMax int32,buildNumMax int32, err error){
	cpuMax, err = GetMaxCpu(activeNodes)
	if err != nil {
		LogErrPrint(err)

	}
	cpuFreqMax, err = GetMaxCpuFreq(activeNodes)
	if err != nil {
		LogErrPrint(err)
	}
	freeMemMax, err = GetMaxFreeMemory(activeNodes)
	if err != nil {
		LogErrPrint(err)
	}
	buildNumMax,err = GetMaxBuildNum(activeNodes)
	return cpuMax, cpuFreqMax, freeMemMax, buildNumMax, err
}
