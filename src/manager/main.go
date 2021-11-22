package manager

import (
	"fmt"
	gen "github.com/sonraisecurity/sonrai-asearch/src/proto"
	"google.golang.org/grpc"
	"log"
	"net"
)

func Main() {
	port := 9881 // TODO not have this hard-coded
	lis, err := net.Listen("tcp", fmt.Sprintf("localhost:%d", port))
	if err != nil {
		log.Fatalf("failed to listen %v", err)
	}

	svc, err := NewManager()
	if err != nil {
		log.Fatalf("failed to create service: %v", err)
	}

	s := grpc.NewServer()
	gen.RegisterManagerServer(s, svc)

	log.Printf("server listening at %v", lis.Addr())
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}