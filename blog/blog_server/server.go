package main

import (
	"context"
	"fmt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"grpc-go-course/blog/blogpb"
	"grpc-go-course/db"
	"log"
	"net"
	"os"
	"os/signal"
)

type server struct{}

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	address := "0.0.0.0:50051"

	client, err := db.InitClient()
	if err != nil {
		log.Fatalf("Failed to connect to mongoDB %v", err)
	}

	client.Database("mydb").Collection("blog")

	listen, err := net.Listen("tcp", address)
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	opts := []grpc.ServerOption{}
	tls := false
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
