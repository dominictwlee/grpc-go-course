package main

import (
	"context"
	"fmt"
	"google.golang.org/grpc"
	"grpc-go-course/calculator/calculatorpb"
	"io"
	"log"
	"net"
)

type server struct{}

func (s *server) FindMaximum(stream calculatorpb.CalculatorService_FindMaximumServer) error {
	var curMax int64
	for {
		req, err := stream.Recv()
		if err == io.EOF {
			return nil
		}
		if err != nil {
			return err
		}

		curNum := req.GetNumber()
		if curNum > curMax {
			curMax = curNum
			if err := stream.Send(&calculatorpb.FindMaximumResponse{MaxNumber: curMax}); err != nil {
				return err
			}
		}
	}
}

func (s *server) ComputeAverage(stream calculatorpb.CalculatorService_ComputeAverageServer) error {
	var nums []int64
	for {
		req, err := stream.Recv()
		if err == io.EOF {
			var sum int64
			for _, n := range nums {
				sum += n
			}
			mean := float64(sum) / float64(len(nums))
			return stream.SendAndClose(&calculatorpb.ComputeAverageResponse{Mean: mean})
		}
		if err != nil {
			log.Fatalf("Failed to read client stream %v", err)
		}
		nums = append(nums, req.GetNumber())
	}
}

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
