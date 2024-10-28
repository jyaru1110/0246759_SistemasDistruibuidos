package main

import (
	"fmt"
	"net"
	api "server/api/v1"
	"server/todo"
	servergrpc "server/todoService"

	"google.golang.org/grpc"
)

func main() {
	lis, err := net.Listen("tcp", ":8081")
	if err != nil {
		fmt.Println(err)
	}

	database := todo.NewDatabase()

	s := grpc.NewServer()

	todoController := todo.NewTodoController(database)

	if err != nil {
		fmt.Println(err)
	}

	fmt.Println("Server is running on port :8081")

	api.RegisterTodoServiceServer(s, &servergrpc.GrpcServer{Todo: todoController})
	s.Serve(lis)
}
