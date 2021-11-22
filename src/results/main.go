package results

import (
	"fmt"
	gen "github.com/sonraisecurity/sonrai-asearch/src/proto"
	"google.golang.org/grpc"
	"log"
	"net"
)

func Main() {
	// TODO
	// - not hardcode port
	// - make this a util or something reusable (maybe)

	port := 9882
	lis, err  := net.Listen("tcp", fmt.Sprintf("localhost:%d", port))
	if err != nil {
		log.Fatalf("failed to listen %v", err)
	}

	svc, err := NewResultsService()
	if err != nil {
		log.Fatalf("failed to create service: %v", err)
	}

	s := grpc.NewServer()
	gen.RegisterResultsServer(s, svc)

	log.Printf("server listening at %v", lis.Addr())
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}