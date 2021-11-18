package main

import (
	"fmt"
	gen "github.com/sonraisecurity/sonrai-asearch/src/proto"
	"google.golang.org/grpc"
	"log"
	"net"
)

func main() {
	port := 9881 // TODO not have this hard-coded
	lis, err := net.Listen("tcp", fmt.Sprintf("localhost:%d", port))
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	mgr, err := NewManager()
	if err != nil {
		log.Fatalf("Failed to create Manager: %v", err)
	}

	s := grpc.NewServer()
	gen.RegisterManagerServer(s, mgr)

	log.Printf("server listening at %v", lis.Addr())
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}