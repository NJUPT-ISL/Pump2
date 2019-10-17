package server

import (
	"context"
	"github.com/Mr-Linus/Pump2/pkg/pump2"
	pu "github.com/Mr-Linus/Pump2/pkg/pump2"
	"google.golang.org/grpc"
	"log"
	"net"
)

type P2Server struct {
	pump2.UnimplementedPump2Server
}

func (s *P2Server) BuildImages(ctx context.Context, in *pump2.BuildInfo) (*pump2.BuildResult, error)  {
	log.Printf("Received: %v", in.GetName())
	return &pump2.BuildResult{ BuildStats: false, ImageName:"Test"}, nil
}

func Start (IP string, Port string) {
	lis, err := net.Listen("tcp",IP+":"+Port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	pu.RegisterPump2Server(s, &P2Server{})
	if err := s.Serve(lis); err != nil{
		log.Fatalf("failed to serve: %v", err)
	}
}