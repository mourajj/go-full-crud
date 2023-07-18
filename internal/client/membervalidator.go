package grpcclient

import (
	"log"
	"os"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func StartGRPC() *grpc.ClientConn {
	hostname := os.Getenv("DOCKER_INTERNAL")
	conn, err := grpc.Dial(hostname+":9000", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatal(err)
	}

	return conn
}
