package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"
	"path/filepath"
	"time"

	pb "github.com/liuchang8877/go-web-word/document"
	"github.com/nguyenthenguyen/docx"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

type server struct {
	pb.UnimplementedDocumentServiceServer
}

func (s *server) GenerateDocument(ctx context.Context, req *pb.GenerateRequest) (*pb.GenerateResponse, error) {
	// 检查模板文件是否存在
	templatePath := "./template.docx"
	if _, err := os.Stat(templatePath); os.IsNotExist(err) {
		log.Printf("Error: Template file %s does not exist", templatePath)
		return nil, fmt.Errorf("template file does not exist")
	}

	// 读取模板文件
	r, err := docx.ReadDocxFile(templatePath)
	if err != nil {
		log.Printf("Error reading template file %s: %v", templatePath, err)
		return nil, fmt.Errorf("failed to read template file")
	}
	defer r.Close()

	doc := r.Editable()
	// 替换占位符
	doc.Replace("{{title}}", req.Title, -1)
	doc.Replace("{{content}}", req.Content, -1)

	// 确保documents目录存在
	outputDir := "./documents"
	if _, err := os.Stat(outputDir); os.IsNotExist(err) {
		os.Mkdir(outputDir, 0755)
	}

	// 保存到文件
	filename := fmt.Sprintf("document-%d.docx", time.Now().Unix())
	filepath := filepath.Join(outputDir, filename)
	err = doc.WriteToFile(filepath)
	if err != nil {
		log.Printf("Error writing to file %s: %v", filepath, err)
		return nil, fmt.Errorf("failed to write to file")
	}

	// 返回下载链接
	downloadURL := fmt.Sprintf("%s/%s", "/documents", filename)

	return &pb.GenerateResponse{DownloadUrl: downloadURL}, nil
}

func main() {
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	s := grpc.NewServer()
	pb.RegisterDocumentServiceServer(s, &server{})
	// Register reflection service on gRPC server.
	reflection.Register(s)

	log.Println("gRPC server is running on port 50051...")
	if err := s.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}
