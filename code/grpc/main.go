package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net"
	"time"

	pb "code/grpc/api"
	"code/grpc/middleware"

	"google.golang.org/grpc"
)

var (
	port = flag.Int("port", 8080, "The server port")
)

type server struct {
	pb.UnimplementedHelloServer
}

func (s *server) Say(ctx context.Context, in *pb.SayRequest) (*pb.SayResponse, error) {
	log.Printf("Received: %v", in.GetName())
	return &pb.SayResponse{Message: "Hello " + in.GetName()}, nil
}

func main() {
	flag.Parse()
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", *port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	interceptor := grpc.ChainUnaryInterceptor(
		middleware.Recovery(),
		middleware.Ratelimit(&middleware.AlwaysPassLimiter{}),
		middleware.Timeout(middleware.WithTimeout(100*time.Millisecond)),
	)

	s := grpc.NewServer(interceptor)
	pb.RegisterHelloServer(s, &server{})
	defer func() {
		s.Stop()
		lis.Close()
	}()
	log.Printf("server listening at %v", lis.Addr())
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
