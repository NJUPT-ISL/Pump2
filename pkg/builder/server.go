package builder

import (
	"bytes"
	"context"
	e "errors"
	"github.com/Mr-Linus/Pump2/pkg/operations"
	pu "github.com/Mr-Linus/Pump2/rpc"
	"github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/load"
	"github.com/shirou/gopsutil/mem"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"net"
	"os"
	"runtime"
	"strings"
)

var (
	nodeStats  bool
	nodeHealth string
	buildNum   int32 = 0
)

type P2Server struct {
	pu.UnimplementedPump2Server
}

func nodeStatsToString(stats bool) string {
	if stats {
		return "true"
	} else {
		return "false"
	}
}

func (s *P2Server) BuildImages(ctx context.Context, in *pu.BuildInfo) (*pu.BuildResult, error) {
	if nodeStats && nodeHealth == "ready" {
		_, err := operations.ConfigDockerfile(in)
		if err != nil {
			return &pu.BuildResult{BuildStats: false, ImageName: ""}, err
		}
		args := operations.ConfigBuildArgs(in)
		LogPrint("Start Build Image: " + in.GetName())
		buildNum++
		res, err := operations.ImageBuild(in.GetName(), &args)
		if err != nil {
			LogErrPrint(err)
			buildNum--
			nodeStats = false
			return &pu.BuildResult{BuildStats: false, ImageName: ""}, err
		}
		var buf = new(bytes.Buffer)
		_, err = buf.ReadFrom(res.Body)
		if err != nil {
			LogErrPrint(err)
		}
		LogPrint(buf.String())
		nodeStats = true
		err = os.Remove(os.Getenv("HOME") + "/Archive.tar")
		if err != nil {
			LogErrPrint(err)
		}
		if strings.Contains(buf.String(), "error") {
			buildNum--
			return &pu.BuildResult{BuildStats: false, ImageName: ""}, nil
		}
		buildNum--
		return &pu.BuildResult{BuildStats: true, ImageName: in.GetName()}, nil
	} else {
		return &pu.BuildResult{BuildStats: false, ImageName: ""}, e.New("Build image failed. Node Status is" +
			nodeStatsToString(nodeStats) + " and node is " + nodeHealth)
	}
}

func (s *P2Server) NodeStatus(ctx context.Context, in *pu.NodeInfo) (*pu.NodeStat, error) {
	var mhz float64 = 0
	lo, err := load.Avg()
	if err != nil {
		LogErrPrint(err)
	}
	cpuInfo, err := cpu.Info()
	if err != nil {
		LogErrPrint(err)
	}
	for _, i := range cpuInfo {
		mhz += i.Mhz
	}
	mhz /= float64(len(cpuInfo))
	memInfo, err := mem.VirtualMemory()
	if err != nil {
		LogErrPrint(err)
	}
	return &pu.NodeStat{
		NodeStats:  nodeStats,
		NodeHealth: nodeHealth,
		BuildNum:   buildNum,
		Cpu:        int32(runtime.NumCPU()),
		LoadAvg:    float32(lo.Load5),
		CpuFreq:    float32(mhz / 1024),
		Memory:     int32(memInfo.Total / 1073741824),
		MemoryFree: int32(memInfo.Available / 1073741824),
	}, nil
}

func StartWithoutTLS(IP string, Port string) {
	nodeHealth = "ready"
	nodeStats = true
	listen, err := net.Listen("tcp", IP+":"+Port)
	LogPrint("Pump2 Server is running at: " + IP + ":" + string(Port))
	if err != nil {
		LogErrPrint(err)
	}
	s := grpc.NewServer()
	pu.RegisterPump2Server(s, &P2Server{})
	if err := s.Serve(listen); err != nil {
		LogErrPrint(err)
	}
}

func StartWithTLS(IP string, Port string, tlsCrtFile string, tlsKeyfile string) {
	nodeHealth = "ready"
	nodeStats = true
	listen, err := net.Listen("tcp", IP+":"+Port)
	if err != nil {
		LogErrPrint(err)
	}
	cred, err := credentials.NewClientTLSFromFile(tlsCrtFile, tlsKeyfile)
	if err != nil {
		LogErrPrint(err)
	}
	s := grpc.NewServer(grpc.Creds(cred))
	pu.RegisterPump2Server(s, &P2Server{})
	if err := s.Serve(listen); err != nil {
		LogErrPrint(err)
	}
}
