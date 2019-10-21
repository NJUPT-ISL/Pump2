package client

import (
	"context"
	pb "github.com/Mr-Linus/Pump2/pkg/pump2"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"log"
	"time"
)

const (
	address     = "localhost:10088"
)

func RunTestWithoutTLS() {
	// Set up a connection to the server.
	conn, err := grpc.Dial(address, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	c := pb.NewPump2Client(conn)
	// Contact the server and print out its response.
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	r, err := c.BuildImages(ctx,
		&pb.BuildInfo{
			Name:"test",
			Gpu:true,
			Tf:true,
			Torch:false,
			TfVersion:"123",
			TorchVersion:"23",
			Dependence:"DEP",
			UseToTest:false,
		})
	if err != nil {
		log.Fatalf("could not greet: %v", err)
	}
	log.Printf("Recv: %s", r.GetImageName())
}

func RunTestWithTLS(tlsKeyfile string, domainName string) {
	creds, err := credentials.NewClientTLSFromFile(tlsKeyfile, domainName)
	if err != nil {
		log.Fatalf("failed to load credentials: %v", err)
	}
	conn, err := grpc.Dial(address, grpc.WithTransportCredentials(creds))
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	c := pb.NewPump2Client(conn)
	// Contact the server and print out its response.
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	r, err := c.BuildImages(ctx,
		&pb.BuildInfo{
			Name:"test",
			Gpu:true,
			Tf:true,
			Torch:false,
			TfVersion:"123",
			TorchVersion:"23",
			Dependence:"DEP",
			UseToTest:false,
		})
	if err != nil {
		log.Fatalf("could not greet: %v", err)
	}
	log.Printf("Recv: %s", r.GetImageName())
}