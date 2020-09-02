package main

import (
	"context"
	"fmt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"grpc-go-course/blog/blogpb"
	"log"
)

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	certFile := "ssl/ca.crt"
	tls := true
	opts := grpc.WithInsecure()
	if tls {
		creds, err := credentials.NewClientTLSFromFile(certFile, "")
		if err != nil {
			log.Fatalf("Could not construct credentials %v", err)
		}
		opts = grpc.WithTransportCredentials(creds)
	}

	cc, err := grpc.Dial("localhost:50051", opts)
	if err != nil {
		log.Fatalf("Could not connect: %v", err)
	}
	defer cc.Close()

	c := blogpb.NewBlogServiceClient(cc)

	fmt.Println("Creating blog")
	blog := &blogpb.Blog{
		AuthorId: "Dom",
		Title:    "My First Blog",
		Content:  "Content of first blog",
	}
	createRes, err := c.CreateBlog(context.Background(), &blogpb.CreateBlogRequest{Blog: blog})
	if err != nil {
		log.Fatalf("Unexpected error: %v", err)
	}
	fmt.Printf("Blog has been created %v", createRes)

}
