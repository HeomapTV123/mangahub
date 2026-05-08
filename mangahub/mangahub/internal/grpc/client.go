package grpc

import (
	"context"
	"fmt"
	"log"
	"time"

	"google.golang.org/grpc"

	pb "mangahub/proto/mangahub/proto"
)

func TestClient() {
	conn, err := grpc.Dial(
		"localhost:50051",
		grpc.WithInsecure(),
	)

	if err != nil {
		log.Fatal(err)
	}

	defer conn.Close()

	client := pb.NewMangaServiceClient(conn)

	ctx, cancel := context.WithTimeout(
		context.Background(),
		3*time.Second,
	)
	defer cancel()

	res, err := client.GetManga(
		ctx,
		&pb.GetMangaRequest{
			Id: "one-piece",
		},
	)

	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Manga:", res.Title)
}
