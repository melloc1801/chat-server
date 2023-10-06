package main

import (
	desc "chat_server/pkg/chat_v1"
	"context"
	"fmt"
	"github.com/fatih/color"
	"github.com/golang/protobuf/ptypes/empty"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"log"
	"net"
)

const grpcPort = 50052

type server struct {
	desc.UnimplementedChatV1Server
}

func (s *server) Create(_ context.Context, req *desc.CreateChatRequest) (*desc.CreateChatResponse, error) {
	fmt.Println("Create request")
	fmt.Println(color.GreenString("", req.Usernames))
	fmt.Println("========================================")

	return &desc.CreateChatResponse{Id: 2}, nil
}
func (s *server) Delete(_ context.Context, req *desc.DeleteRequest) (*empty.Empty, error) {
	fmt.Println("Delete request")
	fmt.Println(color.GreenString("Id", req.Id))
	fmt.Println("========================================")

	return &empty.Empty{}, nil
}
func (s *server) SendMessage(_ context.Context, req *desc.SendMessageRequest) (*empty.Empty, error) {
	fmt.Println("SendMessage request")
	fmt.Println(color.GreenString("From", req.From))
	fmt.Println(color.GreenString("Text", req.Text))
	fmt.Println(color.GreenString("Timestamp", req.Timestamp))
	fmt.Println("========================================")

	return &empty.Empty{}, nil
}

func main() {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", grpcPort))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	s := grpc.NewServer()
	reflection.Register(s)
	desc.RegisterChatV1Server(s, &server{})

	if err = s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
