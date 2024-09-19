package server

import (
	"fmt"
	"net"
	api "server/api/v1"
	log "server/log"

	"google.golang.org/grpc"
)

func main() {
	lis, err := net.Listen("tcp", ":8080")
	if err != nil {
		fmt.Println(err)
	}

	s := grpc.NewServer()
	config := log.Config{}

	commitLog, err := log.NewLog("/data", config)
	if err != nil {
		fmt.Println(err)
	}

	api.RegisterLogServer(s, &grpcServer{CommitLog: commitLog})
	s.Serve(lis)
}