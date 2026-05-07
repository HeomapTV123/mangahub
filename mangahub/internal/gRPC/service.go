package grpc

import (
	"context"
	"fmt"

	pb "mangahub/proto/mangahub/proto"
)

type MangaServer struct {
	pb.UnimplementedMangaServiceServer
}

func (s *MangaServer) GetManga(
	ctx context.Context,
	req *pb.GetMangaRequest,
) (*pb.MangaResponse, error) {

	if req.Id == "" {
		return nil, fmt.Errorf("manga id required")
	}

	return &pb.MangaResponse{
		Id:     req.Id,
		Title:  "One Piece",
		Author: "Oda Eiichiro",
	}, nil
}

func (s *MangaServer) SearchManga(
	ctx context.Context,
	req *pb.SearchRequest,
) (*pb.SearchResponse, error) {

	results := []*pb.MangaResponse{
		{
			Id:     "naruto",
			Title:  "Naruto",
			Author: "Kishimoto",
		},
		{
			Id:     "bleach",
			Title:  "Bleach",
			Author: "Kubo",
		},
	}

	return &pb.SearchResponse{
		Results: results,
	}, nil
}

func (s *MangaServer) UpdateProgress(
	ctx context.Context,
	req *pb.ProgressRequest,
) (*pb.ProgressResponse, error) {

	if req.UserId == "" || req.MangaId == "" {
		return nil, fmt.Errorf("missing fields")
	}

	return &pb.ProgressResponse{
		Status: "progress updated",
	}, nil
}
