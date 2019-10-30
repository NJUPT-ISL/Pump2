package server

import (
	"bytes"
	"context"
	e "errors"
	"fmt"
	"github.com/Mr-Linus/Pump2/pkg/operations"
	"github.com/Mr-Linus/Pump2/pkg/pump2"
	pu "github.com/Mr-Linus/Pump2/pkg/pump2"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"log"
	"net"
	"os"
	"strings"
)
var (
	nodeStats bool
	nodeHealth string
)
type P2Server struct {
	pump2.UnimplementedPump2Server
}

func nodeStatsToString(stats bool) string {
	if stats{
		return "true"
	} else {
		return "false"
	}
}

func (s *P2Server) BuildImages(ctx context.Context, in *pump2.BuildInfo) (*pump2.BuildResult, error)  {
	if nodeStats && nodeHealth == "ready"{
		_, err := operations.ConfigDockerfile(in)
		if err != nil {
			return &pump2.BuildResult{BuildStats:false,ImageName:""}, err
		}
		args := operations.ConfigBuildArgs(in)
		nodeStats = false
		log.Println("Start Build Image: "+in.GetName())
		res, err := operations.ImageBuild(in.GetName(), args)
		if err != nil {
			nodeStats = true
			return &pump2.BuildResult{BuildStats:false,ImageName:""}, err
		}
		var buf = new(bytes.Buffer)
		_, err = buf.ReadFrom(res.Body)
		if err != nil {
			fmt.Println(err)
		}
		log.Println(buf.String())
		nodeStats = true
		err = os.Remove(os.Getenv("HOME")+"/Archive.tar")
		if err != nil {
			fmt.Println(err)
		}
		if strings.Contains(buf.String(),"error"){
			return &pump2.BuildResult{BuildStats:false, ImageName:""},nil
		}
		return &pump2.BuildResult{BuildStats:true, ImageName:in.GetName()},nil
	} else {
		return &pump2.BuildResult{BuildStats:false,ImageName:""}, e.New("Build image failed. Node Status is"+
			nodeStatsToString(nodeStats)+" and node is "+nodeHealth)
	}
}

func (s *P2Server) NodeStatus(ctx context.Context, in *pump2.NodeInfo) (*pump2.NodeStat, error)  {
	return &pump2.NodeStat{NodeHealth:nodeHealth,NodeStats:nodeStats},nil
}

func StartWithoutTLS (IP string, Port string) {
	nodeHealth = "ready"
	nodeStats = true
	listen, err := net.Listen("tcp",IP+":"+Port)
	log.Println("Pump2 Server is running at: "+IP+":"+string(Port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	pu.RegisterPump2Server(s, &P2Server{})
	if err := s.Serve(listen); err != nil{
		log.Fatalf("failed to serve: %v", err)
	}
}

func StartWithTLS (IP string, Port string, tlsCrtFile string, tlsKeyfile string) {
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
