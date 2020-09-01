package main

import (
	"context"
	"fmt"
	"google.golang.org/grpc"
	"grpc-go-course/calculator/calculatorpb"
	"log"
	"net"
)

type server struct{}

func (s *server) DecomposePrimeNumber(req *calculatorpb.PrimeNumberDecompositionRequest, stream calculatorpb.CalculatorService_DecomposePrimeNumberServer) error {
	k := int64(2)
	n := req.GetNumber()
	for n > 1 {
		if n%k == 0 {
			res := &calculatorpb.PrimeNumberDecompositionResponse{Result: k}
			if err := stream.Send(res); err != nil {
				return err
			}
			n = n / k
		} else {
			k++
		}

	}
	return nil
}

func (s *server) Sum(ctx context.Context, req *calculatorpb.SumRequest) (*calculatorpb.SumResponse, error) {
	fmt.Println("Sum function was invoked")
	sum := req.GetNum_1() + req.GetNum_2()
	res := &calculatorpb.SumResponse{Result: sum}
	return res, nil
}

func main() {
	address := "0.0.0.0:50051"
	listener, err := net.Listen("tcp", address)
	if err != nil {
		log.Fatalf("Failed to listen %v", err)
	}

	s := grpc.NewServer()

	calculatorpb.RegisterCalculatorServiceServer(s, &server{})

	fmt.Printf("Server started at %s\n", address)

	if err := s.Serve(listener); err != nil {
		log.Fatalf("Failed to serve %v", err)
	}
}
