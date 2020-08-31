package main

import (
	"google.golang.org/grpc"
	"grpc-go-course/calculator/calculatorpb"
	"log"
	"context"
)

func main() {
	cc, err := grpc.Dial("localhost:50051", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Could not connect %v", err)
	}

	defer cc.Close()

	c := calculatorpb.NewSumServiceClient(cc)
	doUnary(c)
}
func doUnary(c calculatorpb.SumServiceClient) {
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


