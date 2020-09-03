package main

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/status"
	"grpc-go-course/blog/blogpb"
	"grpc-go-course/blog/models"
	"grpc-go-course/db"
	"log"
	"net"
	"os"
	"os/signal"
)

var collection *mongo.Collection

type server struct{}

func (s *server) ListBlogs(req *blogpb.ListBlogsRequest, stream blogpb.BlogService_ListBlogsServer) error {
	blogs, err := models.ListAll(collection)
	if err != nil {
		return status.Errorf(codes.Internal, "Failed to list blogs %v", err)
	}

	for _, b := range blogs {
		if err := stream.Send(&blogpb.ListBlogsResponse{Blog: mapDataToBlogpb(b)}); err != nil {
			return err
		}
	}

	return nil
}

func (s *server) DeleteBlog(ctx context.Context, req *blogpb.DeleteBlogRequest) (*blogpb.DeleteBlogResponse, error) {
	id := req.GetBlogId()
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "Cannot parse ID")
	}

	item := models.BlogItem{
		ID: oid,
	}

	res, err := item.Delete(collection)
	if err != nil {
		return nil, status.Error(codes.Internal, fmt.Sprintf("Failed to delete blog %v", err))
	}
	fmt.Printf("Deleted %v blog(s)", res.DeletedCount)
	return &blogpb.DeleteBlogResponse{Blog: mapDataToBlogpb(item)}, nil
}

func (s *server) UpdateBlog(ctx context.Context, req *blogpb.UpdateBlogRequest) (*blogpb.UpdateBlogResponse, error) {
	blog := req.GetBlog()
	id := blog.GetId()
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "Cannot parse ID")
	}

	item := models.BlogItem{
		ID:       oid,
		AuthorID: blog.GetAuthorId(),
		Content:  blog.GetContent(),
		Title:    blog.GetTitle(),
	}

	data, err := item.Update(collection)

	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, status.Error(codes.NotFound, fmt.Sprintf("Cannot find blog with specified ID %v", err))
		}
		return nil, status.Error(codes.Internal, fmt.Sprintf("Unexpected Error %v", err))
	}
	return &blogpb.UpdateBlogResponse{Blog: mapDataToBlogpb(*data)}, nil
}

func (s *server) ReadBlog(ctx context.Context, req *blogpb.ReadBlogRequest) (*blogpb.ReadBlogResponse, error) {
	blogId := req.GetBlogId()
	oid, err := primitive.ObjectIDFromHex(blogId)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "Cannot parse ID")
	}
	data := models.BlogItem{
		ID: oid,
	}
	_, err = data.ById(collection)
	if err != nil {
		return nil, status.Error(codes.NotFound, fmt.Sprintf("Cannot find blog with specified ID %v", err))
	}
	return &blogpb.ReadBlogResponse{
		Blog: mapDataToBlogpb(data),
	}, nil
}

func (s *server) CreateBlog(ctx context.Context, req *blogpb.CreateBlogRequest) (*blogpb.CreateBlogResponse, error) {
	blog := req.GetBlog()
	item := models.BlogItem{
		AuthorID: blog.GetAuthorId(),
		Content:  blog.GetContent(),
		Title:    blog.GetTitle(),
	}

	res, err := item.Create(collection)
	if err != nil {
		return nil, status.Errorf(codes.Internal, fmt.Sprintf("Internal error: %v", err))
	}
	oid, ok := res.InsertedID.(primitive.ObjectID)
	if !ok {
		return nil, status.Errorf(codes.Internal, fmt.Sprintf("Cannot convert to OID", err))
	}

	return &blogpb.CreateBlogResponse{Blog: &blogpb.Blog{
		Id:       oid.Hex(),
		AuthorId: blog.GetAuthorId(),
		Title:    blog.GetTitle(),
		Content:  blog.GetContent(),
	}}, nil
}

func mapDataToBlogpb(data models.BlogItem) *blogpb.Blog {
	return &blogpb.Blog{
		Id:       data.ID.Hex(),
		AuthorId: data.AuthorID,
		Title:    data.Title,
		Content:  data.Content,
	}
}

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	address := "0.0.0.0:50051"

	client, err := db.InitClient()
	if err != nil {
		log.Fatalf("Failed to connect to mongoDB %v", err)
	}

	collection = client.Database("mydb").Collection("blog")

	listen, err := net.Listen("tcp", address)
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	opts := []grpc.ServerOption{}
	tls := true
	if tls {
		certFile := "ssl/server.crt"
		keyFile := "ssl/server.pem"
		creds, err := credentials.NewServerTLSFromFile(certFile, keyFile)
		if err != nil {
			log.Fatalf("Failed loading certs %v", err)
		}

		opts = append(opts, grpc.Creds(creds))
	}

	s := grpc.NewServer(opts...)
	blogpb.RegisterBlogServiceServer(s, &server{})

	go func() {
		fmt.Println("Starting Server at " + address)
		if err := s.Serve(listen); err != nil {
			log.Fatalf("Failed to serve %v", err)
		}
	}()

	//Wait for ctrl+c to exit
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, os.Interrupt)

	//Block until a signal is received
	<-ch
	fmt.Println("Stopping the server")
	s.Stop()
	fmt.Println("Closing the listener")
	listen.Close()
	fmt.Println("Closing MongoDB connection")
	client.Disconnect(context.TODO())
	fmt.Println("End of Program")
}
