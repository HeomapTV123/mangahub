package grpc

import (
	"log"
	"net"

	"google.golang.org/grpc"

	pb "mangahub/proto/mangahub/proto"
)

func StartGRPCServer() {
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatal("[gRPC] Listen error:", err)
	}

	server := grpc.NewServer()

	pb.RegisterMangaServiceServer(
		server,
		&MangaServer{},
	)

	log.Println("[gRPC] Server running on :50051")

	if err := server.Serve(lis); err != nil {
		log.Fatal("[gRPC] Serve error:", err)
	}
}
