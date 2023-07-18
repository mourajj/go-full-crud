package grpcclient

import (
	"log"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func StartGRPC() *grpc.ClientConn {
	conn, err := grpc.Dial("host.docker.internal:9000", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatal(err)
	}

	return conn
}
