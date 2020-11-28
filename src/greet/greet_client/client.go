package main

import (
	"context"
	"fmt"
	"io"
	"log"

	"github.com/royzh/grpc/src/greet/greetpb"

	"google.golang.org/grpc"
)

func doUnary(c greetpb.GreetServiceClient) {
	log.Println("Doing Unary RPC")
	// Define our request
	req := &greetpb.GreetRequest{
		Greeting: &greetpb.Greeting{
			FirstName: "Roy",
			LastName:  "Zhang",
		},
	}

	// Call greet. Look at the GreetServiceClient to see how its defined
	res, err := c.Greet(context.Background(), req)
	if err != nil {
		log.Fatalln("Error while calling greet", err)
	}
	log.Println("Response from server", res.Result)
}

func doServerStreaming(c greetpb.GreetServiceClient) {
	log.Println("Doing server streaming rpc")

	req := &greetpb.GreetManyTimesRequest{
		Greeting: &greetpb.Greeting{
			FirstName: "Roy",
			LastName:  "Zhang",
		},
	}
	// Call GreetManyTimes. Look at the GreetServiceClient to see how its defined
	resStream, err := c.GreetManyTimes(context.Background(), req)
	if err != nil {
		log.Fatalln("Error while calling GreetManyTimes", err)
	}
	// Setting up a stream gets me GreetService_GreetManyTimesServer and GreetService_GreetManyTimesClient structs
	// This method is in GreetService_GreetManyTimesClient
	for {
		msg, err := resStream.Recv()
		if err == io.EOF {
			// When the stream reaches the end, it throws an EOF error
			// Break out of loop when its done
			break
		}
		if err != nil {
			log.Fatalln("Error while reading stream", err)
		}
		log.Println("Response from GreetManyTimes", msg.GetResult())
	}
}

func doClientStreaming(c greetpb.GreetServiceClient) {
	log.Println("Doing client streaming rpc")

	// Slice of data we're going to be adding
	names := []string{"List", "Of", "Strings", "Can't", "Relate"}
	// Create a list of requests
	var requests []*greetpb.LongGreetRequest

	// Start adding some data into our slice
	for _, name := range names {
		req := &greetpb.LongGreetRequest{
			Greeting: &greetpb.Greeting{
				FirstName: name,
			},
		}
		requests = append(requests, req)
	}

	stream, err := c.LongGreet(context.Background())
	if err != nil {
		log.Fatalln("Error while calling LongGreet", err)
	}

	// Range over our slice of requests to send data over a stream
	for _, req := range requests {
		stream.Send(req)
	}

	// Close and receive stream
	res, err := stream.CloseAndRecv()
	if err != nil {
		log.Fatalln("Can't close and receive stream for some reason", err)
	}
	log.Println("Response from server", res.GetResult())
}

func doBidirectionalStreaming(c greetpb.GreetServiceClient) {
	log.Println("Doing bidrectional streaming")

	// Slice of strings we're going to be adding to greetings
	names := []string{"Hello", "These", "Are", "Some", "Weird", "Names"}
	// List of requests
	var requests []*greetpb.GreetEveryoneRequest

	for _, name := range names {
		req := &greetpb.GreetEveryoneRequest{
			Greeting: &greetpb.Greeting{
				FirstName: name,
			},
		}
		requests = append(requests, req)
	}

	// initialize wait channel
	waitc := make(chan struct{})

	stream, err := c.GreetEveryone(context.Background())
	if err != nil {
		log.Fatalln("Error in greeting everyone", err)
	}

	// Send data with one goroutine
	go func() {
		for _, req := range requests {
			log.Println("Sending message", req.GetGreeting().GetFirstName())
			stream.Send(req)
		}
		stream.CloseSend()
	}()

	// Receive data along another goroutine
	go func() {
		for {
			msg, err := stream.Recv()
			if err == io.EOF {
				// When the stream reaches the end, it throws an EOF error
				// Close the wait channel when we're done
				// That will unblock everything
				break
			}
			if err != nil {
				log.Fatalln("Error while reading stream", err)
				break
			}
			log.Println("Response from GreetEveryone", msg.GetResult())
		}
		close(waitc)
	}()

	// block until everything is done
	<-waitc
}

func main() {
	fmt.Println("I am client")
	// Connect to a running server.
	// Use insecure connection because we don't have a certificate
	conn, err := grpc.Dial("localhost:50051", grpc.WithInsecure())
	if err != nil {
		log.Fatalln("Cannot connect to server")
	}

	// recall that defer means that this executes at the very end
	defer conn.Close()

	// Create the client
	c := greetpb.NewGreetServiceClient(conn)
	//log.Println("Created client", c)

	//doUnary(c)
	//doServerStreaming(c)
	//doClientStreaming(c)
	doBidirectionalStreaming(c)
}
