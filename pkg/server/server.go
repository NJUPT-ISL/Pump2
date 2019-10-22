package server

import (
	"context"
	"github.com/Mr-Linus/Pump2/pkg/pump2"
	pu "github.com/Mr-Linus/Pump2/pkg/pump2"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
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

func StartWithoutTLS (IP string, Port int) {
	listen, err := net.Listen("tcp",IP+":"+string(Port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	pu.RegisterPump2Server(s, &P2Server{})
	if err := s.Serve(listen); err != nil{
		log.Fatalf("failed to serve: %v", err)
	}
}

func StartWithTLS (IP string, Port int, tlsCrtFile string, tlsKeyfile string) {
	listen, err := net.Listen("tcp", IP+":"+string(Port))
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
