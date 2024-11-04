package main

import (
	"fmt"
	"net"
	"server/auth"
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
		CertFile: tlsconfig.ServerCertFile,
		KeyFile:  tlsconfig.ServerKeyFile,
		CAFile:   tlsconfig.CAFile,
		Server:   true,
	})

	if err != nil {
		fmt.Println(err)
	}

	serverCreds := credentials.NewTLS(severTLSConfig)

	clog, err := log.NewLog("./data/logs", log.Config{})

	authorizer := auth.New(tlsconfig.ACLModelFile, tlsconfig.ACLPolicyFile)

	if err != nil {
		fmt.Println(err)
	}

	config := &servergrpc.Config{
		CommitLog:  clog,
		Authorizer: authorizer,
	}

	s, err := servergrpc.NewGRPCServer(config, grpc.Creds(serverCreds))

	if err != nil {
		fmt.Println(err)
	}

	fmt.Println("Server is running on port:8080")
	s.Serve(lis)
}
