package main

import (
	"context"
	"log"
	"time"

	pb "github.com/liuchang8877/go-web-word/document"
	"google.golang.org/grpc"
)

func main() {
	conn, err := grpc.Dial("localhost:50051", grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	c := pb.NewDocumentServiceClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	r, err := c.GenerateDocument(ctx, &pb.GenerateRequest{Title: "Test Title", Content: "Test Content"})
	if err != nil {
		log.Fatalf("could not generate document: %v", err)
	}
	log.Printf("Download URL: %s", r.GetDownloadUrl())
}
