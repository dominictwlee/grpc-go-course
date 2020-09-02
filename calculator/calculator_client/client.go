package main

import (
	"context"
	"fmt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"grpc-go-course/calculator/calculatorpb"
	"io"
	"log"
	"time"
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
	//doClientStream(c)
	//doBiDiStream(c)
	doErrorUnary(c)
}

func doErrorUnary(c calculatorpb.CalculatorServiceClient) {
	doErrorCall(c, 10)
	doErrorCall(c, -20)
}

func doErrorCall(c calculatorpb.CalculatorServiceClient, n int32) {
	req := &calculatorpb.SquareRootRequest{Number: n}
	res, err := c.SquareRoot(context.Background(), req)
	if err != nil {
		resErr, ok := status.FromError(err)
		if ok {
			// error from gRPC
			if resErr.Code() == codes.InvalidArgument {
				fmt.Println(resErr.Message())
			}
		} else {
			log.Fatalf("Error calling SquareRoot: %v", resErr)
		}
		return
	}
	fmt.Println("SquareRoot: ", res.GetResult())
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

func doBiDiStream(c calculatorpb.CalculatorServiceClient) {
	stream, err := c.FindMaximum(context.Background())
	if err != nil {
		log.Fatalf("Failed to call FindMaximum %v", err)
	}

	waitChan := make(chan struct{})

	go func() {
		nums := []int64{1, 5, 3, 6, 2, 20}
		for _, n := range nums {
			if err := stream.Send(&calculatorpb.FindMaximumRequest{Number: n}); err != nil {
				log.Fatalf("Failed to send message %v", err)
			}
			time.Sleep(1000 * time.Millisecond)
		}
		if err := stream.CloseSend(); err != nil {
			log.Fatalf("Failed to close send stream %v", err)
		}
	}()

	go func() {
		for {
			res, err := stream.Recv()
			if err == io.EOF {
				break
			}
			if err != nil {
				log.Printf("Failed to receive message %v", err)
				break
			}
			fmt.Printf("Current Max is: %v\n", res.GetMaxNumber())
		}
		close(waitChan)
	}()
	<-waitChan
}
