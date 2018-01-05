package main

import (
	"context"
	"log"

	pb "github.com/colek42/grpclangserv/api"
	"github.com/davecgh/go-spew/spew"
	"google.golang.org/grpc"
)

var serverAddr = "localhost:4534"

func main() {
	conn, err := grpc.Dial("localhost:4534", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Unable to connect to server %s", serverAddr)
	}
	defer conn.Close()

	client := pb.NewLanguageServerClient(conn)
	ctx := context.Background()
	q := &pb.Query{
		LineNumber: 54,
		CharNumber: 3,
		Pkg:        "github.com/colek42/streamingDemo/packetsender",
		FileName:   "server.go",
	}
	res, err := client.GetDefinition(ctx, q)
	if err != nil {
		log.Println(err)
	}

	spew.Dump(res)
}
