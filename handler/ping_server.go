package handler

import (
	"context"
	"github.com/DangVTNhan/goacl/pb"
)

type PingServer struct {
	pb.UnimplementedPingServer
}

func NewPingServer() *PingServer {
	return &PingServer{}
}

func (s *PingServer) Ping(ctx context.Context, req *pb.PingRequest) (*pb.PingReply, error) {
	return &pb.PingReply{Message: "Pong"}, nil
}
