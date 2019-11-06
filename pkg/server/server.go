package server

import (
	"bytes"
	"context"
	e "errors"
	"fmt"
	"github.com/Mr-Linus/Pump2/pkg/operations"
	"github.com/Mr-Linus/Pump2/pkg/rpc"
	pu "github.com/Mr-Linus/Pump2/pkg/rpc"
	"github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/mem"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"log"
	"net"
	"os"
	"runtime"
	"strings"
)

var (
	nodeStats  bool
	nodeHealth string
	buildNum   int32 = 1
)

type P2Server struct {
	rpc.UnimplementedPump2Server
}

func nodeStatsToString(stats bool) string {
	if stats {
		return "true"
	} else {
		return "false"
	}
}

func (s *P2Server) BuildImages(ctx context.Context, in *rpc.BuildInfo) (*rpc.BuildResult, error) {
	if nodeStats && nodeHealth == "ready" {
		_, err := operations.ConfigDockerfile(in)
		if err != nil {
			return &rpc.BuildResult{BuildStats: false, ImageName: ""}, err
		}
		args := operations.ConfigBuildArgs(in)
		log.Println("Start Build Image: " + in.GetName())
		buildNum++
		res, err := operations.ImageBuild(in.GetName(), args)
		if err != nil {
			buildNum--
			nodeStats = false
			return &rpc.BuildResult{BuildStats: false, ImageName: ""}, err
		}
		var buf = new(bytes.Buffer)
		_, err = buf.ReadFrom(res.Body)
		if err != nil {
			fmt.Println(err)
		}
		log.Println(buf.String())
		nodeStats = true
		err = os.Remove(os.Getenv("HOME") + "/Archive.tar")
		if err != nil {
			fmt.Println(err)
		}
		if strings.Contains(buf.String(), "error") {
			buildNum--
			return &rpc.BuildResult{BuildStats: false, ImageName: ""}, nil
		}
		buildNum--
		return &rpc.BuildResult{BuildStats: true, ImageName: in.GetName()}, nil
	} else {
		return &rpc.BuildResult{BuildStats: false, ImageName: ""}, e.New("Build image failed. Node Status is" +
			nodeStatsToString(nodeStats) + " and node is " + nodeHealth)
	}
}

func (s *P2Server) NodeStatus(ctx context.Context, in *rpc.NodeInfo) (*rpc.NodeStat, error) {
	var (
		cpuUsage float64 = 0
		mhz      float64 = 0
	)
	cpusPer, err := cpu.Percent(1, true)
	if err != nil {
		fmt.Println(err)
	}
	for _, i := range cpusPer {
		cpuUsage += i
	}
	cpuUsage /= float64(len(cpusPer))
	cpuInfo, err := cpu.Info()
	if err != nil {
		fmt.Println(err)
	}
	for _, i := range cpuInfo {
		mhz += i.Mhz
	}
	mhz /= float64(len(cpuInfo))
	memInfo, err := mem.VirtualMemory()
	if err != nil {
		fmt.Println(err)
	}
	return &rpc.NodeStat{
		NodeStats:  nodeStats,
		NodeHealth: nodeHealth,
		BuildNum:   buildNum,
		Cpu:        int32(runtime.NumCPU()),
		CpuUsage:   float32(cpuUsage),
		CpuFreq:    float32(mhz / 1024),
		Memory:     int32(memInfo.Total / 1073741824),
		MemoryFree: int32(memInfo.Available / 1073741824),
	}, nil
}

func StartWithoutTLS(IP string, Port string) {
	nodeHealth = "ready"
	nodeStats = true
	listen, err := net.Listen("tcp", IP+":"+Port)
	log.Println("Pump2 Server is running at: " + IP + ":" + string(Port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	pu.RegisterPump2Server(s, &P2Server{})
	if err := s.Serve(listen); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}

func StartWithTLS(IP string, Port string, tlsCrtFile string, tlsKeyfile string) {
	nodeHealth = "ready"
	nodeStats = true
	listen, err := net.Listen("tcp", IP+":"+Port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	cred, err := credentials.NewClientTLSFromFile(tlsCrtFile, tlsKeyfile)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer(grpc.Creds(cred))
	pu.RegisterPump2Server(s, &P2Server{})
	if err := s.Serve(listen); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
