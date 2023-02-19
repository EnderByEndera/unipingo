package scheduler

import (
	"codefridge/rpc/protocols"
	"context"
	"log"
	"net"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

type MySchedulerServer struct {
	protocols.UnimplementedSchedulerServer
}

func (s *MySchedulerServer) GetTask(ctx context.Context, in *protocols.Service) (*protocols.Task, error) {
	return &protocols.Task{TaskType: protocols.TaskType_BUILD_DOCKER, Message: "should!"}, nil
}

func Main() {
	listener, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("监听失败: %v", err)
	}

	// 创建gRPC服务
	grpcServer := grpc.NewServer()

	// Tester 注册服务实现者
	// 此函数在 test.pb.go 中，自动生成
	protocols.RegisterSchedulerServer(grpcServer, &MySchedulerServer{})

	// 在 gRPC 服务上注册反射服务
	// func Register(s *grpc.Server)
	reflection.Register(grpcServer)

	err = grpcServer.Serve(listener)
	if err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
