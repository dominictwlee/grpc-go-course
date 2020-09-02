package main

import (
	"context"
	"fmt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/status"
	"grpc-go-course/greet/greetpb"
	"io"
	"log"
	"net"
	"strconv"
	"time"
)

type server struct{}

func (s *server) GreetWithDeadline(ctx context.Context, req *greetpb.GreetWithDeadlineRequest) (*greetpb.GreetWithDeadlineResponse, error) {
	fmt.Println("Greet function was invoked")
	for i := 0; i < 3; i++ {
		if ctx.Err() == context.Canceled {
			fmt.Println("Client canceled the request!")
			return nil,status.Error(codes.DeadlineExceeded, "The client canceled the request")
		}
		time.Sleep(1 * time.Second)
	}
	firstName := req.GetGreeting().GetFirstName()
	lastName := req.GetGreeting().GetLastName()
	greeting := fmt.Sprintf("Hello %s %s", firstName, lastName)
	res := &greetpb.GreetWithDeadlineResponse{
		Result: greeting,
	}
	return res, nil
}

func (s *server) GreetEveryone(stream greetpb.GreetService_GreetEveryoneServer) error {
	for {
		req, err := stream.Recv()
		if err == io.EOF {
			return nil
		}
		if err != nil {
			log.Fatalf("Error while reading client stream : %v", err)
			return err
		}
		greeting := req.GetGreeting()
		firstName := greeting.GetFirstName()
		lastName := greeting.GetLastName()

		res := &greetpb.GreetEveryoneResponse{Result: fmt.Sprintf("Hello %v, %v!", firstName, lastName)}
		if err := stream.Send(res); err != nil {
			log.Fatalf("Error while sending data to client: %v", err)
			return err
		}
	}
}

func (s *server) LongGreet(stream greetpb.GreetService_LongGreetServer) error {
	fmt.Println("LongGreet Func was invoked with a streaming request")
	result := ""
	for {
		req, err := stream.Recv()
		if err == io.EOF {
			return stream.SendAndClose(&greetpb.LongGreetResponse{Result: result})
		}
		if err != nil {
			log.Fatalf("Error while reading client stream %v", err)
		}

		firstName := req.GetGreeting().GetFirstName()
		result += fmt.Sprintf("Hello %s ! ", firstName)
	}
}

func (s *server) GreetManyTimes(req *greetpb.GreetManyTimesRequest, stream greetpb.GreetService_GreetManyTimesServer) error {
	greeting := fmt.Sprintf("Hello %s %s", req.GetGreeting().GetFirstName(), req.GetGreeting().GetLastName())
	for i := 0; i < 10; i++ {
		res := &greetpb.GreetManyTimesResponse{
			Result: fmt.Sprintf("%s number %s", greeting, strconv.Itoa(i)),
		}
		if err := stream.Send(res); err != nil {
			return err
		}
		time.Sleep(1000 * time.Millisecond)
	}
	return nil
}

func (*server) Greet(ctx context.Context, req *greetpb.GreetRequest) (*greetpb.GreetResponse, error) {
	fmt.Println("Greet function was invoked")
	firstName := req.GetGreeting().GetFirstName()
	lastName := req.GetGreeting().GetLastName()
	greeting := fmt.Sprintf("Hello %s %s", firstName, lastName)
	res := &greetpb.GreetResponse{
		Result: greeting,
	}
	return res, nil
}

//func (*server) GreetManyTimes(ctx context.Context, in *greetpb.GreetManyTimesRequest, opts ...grpc.CallOption) (GreetService_GreetManyTimesClient, error)

func main() {
	fmt.Println("Hello Server")

	listen, err := net.Listen("tcp", "0.0.0.0:50051")
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

	greetpb.RegisterGreetServiceServer(s, &server{})

	if err := s.Serve(listen); err != nil {
		log.Fatalf("Failed to serve %v", err)
	}
}
