package main

import (
	"context"
	"fmt"
	"google.golang.org/grpc"
	"grpc-go-course/calculator/calculatorpb"
	"io"
	"log"
)

func main() {
	cc, err := grpc.Dial("localhost:50051", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Could not connect %v", err)
	}

	defer cc.Close()

	c := calculatorpb.NewCalculatorServiceClient(cc)
	//doUnary(c)
	//doServerStream(c)
	doClientStream(c)
}
func doUnary(c calculatorpb.CalculatorServiceClient) {
	req := &calculatorpb.SumRequest{
		Num_1: 3,
		Num_2: 7,
	}

	res, err := c.Sum(context.Background(), req)
	if err != nil {
		log.Fatalf("error while calling Sum RPC: %v", err)
	}
	log.Printf("Response from Sum: %v", res.Result)
}

func doServerStream(c calculatorpb.CalculatorServiceClient) {
	req := &calculatorpb.PrimeNumberDecompositionRequest{Number: 120000300000}
	stream, err := c.DecomposePrimeNumber(context.Background(), req)
	if err != nil {
		log.Fatalf("Failed to stream %v", err)
	}

	for {
		msg, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatalf("Failed to read stream %v", err)
		}
		fmt.Println("Prime Num Decomposition: ", msg.GetResult())
	}
}

func doClientStream(c calculatorpb.CalculatorServiceClient) {
	nums := []int64{3, 5, 9, 54, 23}

	stream, err := c.ComputeAverage(context.Background())
	if err != nil {
		log.Fatalf("Error while calling ComputeAverage %v", err)
	}

	for _, number := range nums {
		if err := stream.Send(&calculatorpb.ComputeAverageRequest{
			Number: number,
		}); err != nil {
			log.Printf("Failed to send request %v", err)
		}
	}

	res, err := stream.CloseAndRecv()
	if err != nil {
		log.Fatalf("Failed to receive response from ComputeAverage %v", err)
	}
	fmt.Printf("Mean is: %v\n", res.GetMean())
}
