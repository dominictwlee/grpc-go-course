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
	readBlog(c)
}

func readBlog(c blogpb.BlogServiceClient) {
	fmt.Println("Getting blog")
	readRes, err := c.ReadBlog(context.Background(), &blogpb.ReadBlogRequest{BlogId: "5f50030183dcd90c273af889"})
	if err != nil {
		log.Fatalf("Unexpected error: %v", err)
	}
	fmt.Printf("Blog received %v\n", readRes)
}

func createBlog(c blogpb.BlogServiceClient) {
	fmt.Println("Creating blog")
	blog := &blogpb.Blog{
		AuthorId: "Jane",
		Title:    "My 2nd Blog",
		Content:  "Content of 2nd blog",
	}
	createRes, err := c.CreateBlog(context.Background(), &blogpb.CreateBlogRequest{Blog: blog})
	if err != nil {
		log.Fatalf("Unexpected error: %v", err)
	}
	fmt.Printf("Blog has been created %v\n", createRes)
}
