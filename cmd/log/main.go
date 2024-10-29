package main

import (
	"fmt"
	"net"
	api "server/api/v1"
	tlsconfig "server/config"
	log "server/log"
	servergrpc "server/logService"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

func main() {
	lis, err := net.Listen("tcp", ":8080")
	if err != nil {
		fmt.Println(err)
	}

	severTLSConfig, err := tlsconfig.SetupTLSConfig(tlsconfig.TLSConfig{
		CertFile:      tlsconfig.ServerCertFile,
		KeyFile:       tlsconfig.ServerKeyFile,
		CAFile:        tlsconfig.CAFile,
		ServerAddress: lis.Addr().String(),
		Server:        true,
	})

	if err != nil {
		fmt.Println(err)
	}

	serverCreds := credentials.NewTLS(severTLSConfig)

	s := grpc.NewServer(grpc.Creds(serverCreds))
	config := log.Config{}

	commitLog, err := log.NewLog("./data/logs", config)
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println("Server is running on port :8080")

	api.RegisterLogServer(s, &servergrpc.GrpcServer{CommitLog: commitLog})
	s.Serve(lis)
}
