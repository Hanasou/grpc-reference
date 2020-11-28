package main

import (
	"context"
	"fmt"
	"log"

	"github.com/royzh/grpc/src/blog/blogpb"
	"google.golang.org/grpc"
)

func createBlog(c blogpb.BlogServiceClient, b blogpb.Blog) {
	log.Println("Client making create blog request")
	req := &blogpb.CreateBlogRequest{
		Blog: &b,
	}
	res, err := c.CreateBlog(context.Background(), req)
	if err != nil {
		log.Fatalln("Error in creating blog post", err)
	}
	log.Println("Response from server:", res)
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
	c := blogpb.NewBlogServiceClient(conn)
	//log.Println("Created client", c)
}
