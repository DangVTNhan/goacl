package handler

import (
	"context"
	"github.com/DangVTNhan/goacl/api"
)

type PingServer struct {
	api.UnimplementedPingServiceServer
}

func NewPingServer() *PingServer {
	return &PingServer{}
}

func (s *PingServer) Ping(_ context.Context, _ *api.PingRequest) (*api.PingResponse, error) {
	return &api.PingResponse{Message: "Pong"}, nil
}
