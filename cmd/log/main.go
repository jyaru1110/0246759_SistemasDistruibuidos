package main

import (
	"fmt"
	"net"
	api "server/api/v1"
	log "server/log"
	servergrpc "server/logService"

	"google.golang.org/grpc"
)

func main() {
	lis, err := net.Listen("tcp", ":8080")
	if err != nil {
		fmt.Println(err)
	}

	s := grpc.NewServer()
	config := log.Config{}

	commitLog, err := log.NewLog("./data/logs", config)
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println("Server is running on port :8080")

	api.RegisterLogServer(s, &servergrpc.GrpcServer{CommitLog: commitLog})
	s.Serve(lis)
}
