package handler

import (
	"context"
	"github.com/DangVTNhan/goacl/api"
)

type PingServer struct {
	api.UnimplementedPingServer
}

func NewPingServer() *PingServer {
	return &PingServer{}
}

func (s *PingServer) Ping(_ context.Context, _ *api.PingRequest) (*api.PingReply, error) {
	return &api.PingReply{Message: "Pong"}, nil
}
