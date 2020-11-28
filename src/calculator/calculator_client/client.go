package main

import (
	"context"
	"io"
	"log"

	"github.com/royzh/grpc/src/calculator/calculatorpb"
	"google.golang.org/grpc"
)

func doSum(c calculatorpb.CalculatorServiceClient, x int32, y int32) {
	log.Println("Doing unary rpc")
	req := &calculatorpb.SumRequest{
		Nums: &calculatorpb.TwoNumbers{
			X: x,
			Y: y,
		},
	}
	res, err := c.Sum(context.Background(), req)
	if err != nil {
		log.Fatalln("Can't do shit with this server")
	}
	log.Println("Summed two numbers:", res.Result)
}

func doPrimeDecomp(c calculatorpb.CalculatorServiceClient, n int32) {
	log.Println("Doing prime decomp rpc")
	req := &calculatorpb.PrimeDecompRequest{
		Num: n,
	}

	resStream, err := c.PrimeDecomp(context.Background(), req)
	if err != nil {
		log.Fatalln("Error while calling PrimeDecomp", err)
	}

	// Accept stream of data from server
	for {
		msg, err := resStream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatalln("Error while reading stream", err)
		}
		log.Println(msg.GetResult())
	}
}

func doAverage(c calculatorpb.CalculatorServiceClient, nums []int64) {
	log.Println("Doing average")

	// Create a list of requests
	var requests []*calculatorpb.AverageRequest

	// Add the numbers from nums into requests
	for _, num := range nums {
		req := &calculatorpb.AverageRequest{
			Num: num,
		}
		requests = append(requests, req)
	}

	stream, err := c.Average(context.Background())
	if err != nil {
		log.Fatalln("Error while calling average request", err)
	}

	// Send data over the stream
	for _, req := range requests {
		stream.Send(req)
	}

	// Close and receive stream
	res, err := stream.CloseAndRecv()
	if err != nil {
		log.Fatalln("Can't close and receive", err)
	}
	log.Println("Response from server", res.GetResult())
}

func main() {
	log.Println("Starting client")

	conn, err := grpc.Dial("localhost:50051", grpc.WithInsecure())
	if err != nil {
		log.Fatalln("Could not connect to server", err)
	}

	defer conn.Close()

	c := calculatorpb.NewCalculatorServiceClient(conn)

	//doSum(c, 5, 3)
	//doPrimeDecomp(c, 120)
	doAverage(c, []int64{1, 2, 3, 4, 5})
}
