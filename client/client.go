package main

import (
	"context"
	"flag"
	"fmt"
	api "server/api/v1"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var addr = flag.String("addr", "localhost:8080", "the address to connect to")

func main() {
	conn, err := grpc.NewClient(*addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		fmt.Printf("could not connect: %v", err)
	}
	defer conn.Close()
	c := api.NewLogClient(conn)

	// Contact the server and print out its response.
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	r, err := c.Produce(ctx, &api.ProduceRequest{Record: &api.Record{Value: []byte("hello world")}})
	if err != nil {
		fmt.Printf("could not produce: %v", err)
	}

	fmt.Println(r.Offset)

	consumeClient, err := c.Consume(ctx, &api.ConsumeRequest{Offset: r.Offset})

	if err != nil {
		fmt.Printf("could not consume: %v", err)
	}

	fmt.Println(string(consumeClient.Record.Value))
}
