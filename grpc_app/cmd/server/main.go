package main

import (
	"google.golang.org/grpc"
	"grpc_app/pkg/api"
	nums_streamer "grpc_app/pkg/nums-streamer"
	"log"
	"net"
)

func main() {
	s := grpc.NewServer()
	srv := nums_streamer.NewGRPCServer()
	api.RegisterNumsStreamerServer(s, srv)

	l, err := net.Listen("tcp", ":8080")
	if err != nil {
		log.Fatal(err)
	}

	if err := s.Serve(l); err != nil {
		log.Fatal(err)
	}
}
