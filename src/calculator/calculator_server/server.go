package main

import (
	"context"
	"io"
	"log"
	"math"
	"net"

	"github.com/royzh/grpc/src/calculator/calculatorpb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type server struct{}

// Sum sums two numbers together
func (*server) Sum(ctx context.Context, req *calculatorpb.SumRequest) (*calculatorpb.SumResponse, error) {
	log.Println("Sum function invoked")
	result := req.GetNums().GetX() + req.GetNums().GetY()
	response := &calculatorpb.SumResponse{
		Result: result,
	}
	return response, nil
}

func (*server) SquareRoot(ctx context.Context, req *calculatorpb.SquareRootRequest) (*calculatorpb.SquareRootResponse, error) {
	log.Println("SquareRoot function invoked")
	number := req.GetNumber()
	if number < 0 {
		// We can't sqrt a number < 0 so if we get an argument like that
		// Return an error
		// The error comes from the grpc/status package
		return nil, status.Errorf(codes.InvalidArgument, "Received negative number")
	}
	res := &calculatorpb.SquareRootResponse{
		NumberRoot: math.Sqrt(float64(number)),
	}
	return res, nil
}

// PrimeDecomp gives us the prime decomp of a number
func (*server) PrimeDecomp(req *calculatorpb.PrimeDecompRequest,
	stream calculatorpb.CalculatorService_PrimeDecompServer) error {

	log.Println("Prime Decomp function invoked")
	n := req.GetNum()
	var k int32 = 2
	for n > 1 {
		if n%k == 0 {
			res := &calculatorpb.PrimeDecompResponse{
				Result: k,
			}
			stream.Send(res)
			n = n / k
		} else {
			k = k + 1
		}
	}
	return nil
}

func (*server) Average(stream calculatorpb.CalculatorService_AverageServer) error {
	log.Println("Average called")
	var sum int64
	var count int64
	for {
		req, err := stream.Recv()
		if err == io.EOF {
			average := sum / count
			res := &calculatorpb.AverageResponse{
				Result: average,
			}
			return stream.SendAndClose(res)
		}
		if err != nil {
			log.Println("Stream error", err)
		}
		sum += req.GetNum()
		count++
	}
}

func main() {
	log.Println("Starting server")

	lis, err := net.Listen("tcp", "0.0.0.0:50051")
	if err != nil {
		log.Fatalln("Can't listen to this port", err)
	}

	// Start the server
	s := grpc.NewServer()
	calculatorpb.RegisterCalculatorServiceServer(s, &server{})

	if err := s.Serve(lis); err != nil {
		log.Fatalln("Failed to serve", err)
	}
}
