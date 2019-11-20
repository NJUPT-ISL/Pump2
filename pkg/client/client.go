package client

import (
	"context"
	pb "github.com/Mr-Linus/Pump2/rpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"log"
	"time"
)

const (
	address = "0.0.0.0:10088"
)

func RunTestWithoutTLS() {
	// Set up a connection to the builder.
	conn, err := grpc.Dial(address, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	c := pb.NewPump2Client(conn)
	// Contact the builder and print out its response.
	ctx := context.Background()
	r, err := c.BuildImages(ctx,
		&pb.BuildInfo{
			Name:         "test:latest",
			Gpu:          false,
			Tf:           true,
			Torch:        false,
			TfVersion:    "1.14.0",
			TorchVersion: "",
			Dependence:   "",
			UseToTest:    true,
		})
	if err != nil {
		log.Fatalf("could not greet: %v", err)
	}
	log.Printf("Recv: %s", r.GetImageName())
}

func RunTestWithTLS(tlsKeyfile string, domainName string) {
	cred, err := credentials.NewClientTLSFromFile(tlsKeyfile, domainName)
	if err != nil {
		log.Fatalf("failed to load credentials: %v", err)
	}
	conn, err := grpc.Dial(address, grpc.WithTransportCredentials(cred))
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	c := pb.NewPump2Client(conn)
	// Contact the builder and print out its response.
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	r, err := c.BuildImages(ctx,
		&pb.BuildInfo{
			Name:         "test",
			Gpu:          true,
			Tf:           true,
			Torch:        false,
			TfVersion:    "123",
			TorchVersion: "23",
			Dependence:   "DEP",
			UseToTest:    false,
		})
	if err != nil {
		log.Fatalf("could not greet: %v", err)
	}
	log.Printf("Recv: %s", r.GetImageName())
}
