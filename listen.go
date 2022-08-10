package grpcutil

import (
	"net"
	"log"
  "google.golang.org/grpc"
)

func Serve(s *grpc.Server, port string) {

  lis, err := net.Listen("tcp", port)
  if err != nil {
    log.Fatalf("failed to listen: %v", err)
  }

  log.Print("Starting gRPC..." + port)
  if errsrv := s.Serve(lis); errsrv != nil {
    log.Fatalf("Could not start gRPC server: %v", errsrv)
  }
}
