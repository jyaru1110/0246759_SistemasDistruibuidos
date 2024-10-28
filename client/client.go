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

var addr = flag.String("addr", "localhost:8081", "the address to connect to")

func main() {
	conn, err := grpc.NewClient(*addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		fmt.Printf("could not connect: %v", err)
	}
	defer conn.Close()
	c := api.NewTodoServiceClient(conn)

	// Contact the server and print out its response.
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	r, err := c.ProduceTodo(ctx, &api.ProduceTodoRequest{Todo: &api.Todo{Id: "todo_1_gg", Value: "esto es un todooo"}})
	if err != nil {
		fmt.Printf("could not produce: %v", err)
	}

	fmt.Printf("inserted id: %s", r.Id)

	consumeClient, err := c.Get(ctx, &api.GetRequest{Id: "todo_1_gg"})

	if err != nil {
		fmt.Printf("could not consume: %v", err)
	}

	fmt.Printf("id: %s, value: %s", consumeClient.Todo.Id, consumeClient.Todo.Value)

	fmt.Println()
}
