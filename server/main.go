package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"net/http"

	"google.golang.org/grpc"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	pb "github.com/rcrick/grpc-gateway/proto"
)

type server struct {
	pb.UnimplementedHelloServiceServer
}

func (s *server) SayHello(ctx context.Context, in *pb.HelloRequest) (*pb.HelloReply, error) {
	log.Printf("Receive msg: %s", in.GetName())
	if err := in.Validate(); err != nil {
		log.Print(err)
		return nil, err
	}
	return &pb.HelloReply{Message: fmt.Sprintf("You send server: %s", in.GetName())}, nil
}

const PORT = "8000"

func main() {
	s := grpc.NewServer()
	pb.RegisterHelloServiceServer(s, &server{})

	lis, err := net.Listen("tcp", ":"+PORT)
	if err != nil {
		log.Fatal(err)
	}

	go func() {
		log.Fatal(s.Serve(lis))
	}()

	// Create a client connection to the gRPC server we just started
	// This is where the gRPC-Gateway proxies the requests
	conn, err := grpc.DialContext(
		context.Background(),
		":"+PORT,
		grpc.WithInsecure(),
		grpc.WithBlock())
	if err != nil {
		log.Fatal(err)
	}

	gwmux := runtime.NewServeMux()
	err = pb.RegisterHelloServiceHandler(context.Background(), gwmux, conn)
	if err != nil {
		log.Fatal(err)
	}

	gwServer := &http.Server{
		Addr:    ":8001",
		Handler: gwmux,
	}
	log.Println("Serving gRPC-Gateway on http://localhost:8001/")
	log.Fatalln(gwServer.ListenAndServe())
}
