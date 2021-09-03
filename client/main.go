package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"google.golang.org/grpc"

	pb "github.com/rcrick/grpc-gateway/proto"
)

const PORT = "8000"

func main() {
	conn, err := grpc.Dial(":"+PORT, grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	client := pb.NewHelloServiceClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	res, err := client.SayHello(ctx, &pb.HelloRequest{Name: "I'm client"})
	if err != nil {
		log.Fatal(err)
	}
	fmt.Print(res.GetMessage())
}
