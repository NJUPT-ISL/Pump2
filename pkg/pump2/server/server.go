package server

import (
	"context"
	"github.com/Mr-Linus/Pump2/pkg/pump2"
	"log"
)

type P2Server struct {
	pump2.UnimplementedPump2Server
}

func (s *P2Server) BuildImages(ctx context.Context, in *pump2.BuildInfo) (*pump2.BuildResult, error)  {
	log.Printf("Received: %v", in.GetName())
	return &pump2.BuildResult{ Stats: false, ImageName:"Test"}, nil
}
