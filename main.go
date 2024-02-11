package main

import (
	"database/sql"
	"log"
	"net"

	mydb "github.com/aryyawijaya/simple-bank/db/sqlc"
	"github.com/aryyawijaya/simple-bank/pb"
	servergrpc "github.com/aryyawijaya/simple-bank/server-grpc"
	"github.com/aryyawijaya/simple-bank/server-http"
	"github.com/aryyawijaya/simple-bank/util"
	_ "github.com/lib/pq"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

// const (
// 	driverName     = "postgres"
// 	dataSourceName = "postgresql://root:secretpassword@localhost:5433/simple_bank?sslmode=disable"
// 	serverAddress  = "localhost:8080"
// )

func main() {
	config, err := util.LoadConfig(".")
	if err != nil {
		log.Fatalf("Cannot load config: %v\n", err)
	}
	// conn, err := sql.Open(driverName, dataSourceName)
	conn, err := sql.Open(config.DBDriver, config.DBSource)
	if err != nil {
		log.Fatalf("Cannot connect to database: %v\n", err)
	}

	store := mydb.NewStore(conn)
	runGRPCServer(config, store)
	// runGinServer(config, store)
}

func runGRPCServer(config util.Config, store mydb.Store) {
	// create Implementation RPC server
	server, err := servergrpc.NewServer(store, &config)
	if err != nil {
		log.Fatalf("Cannot create the server: %v\n", err)
	}

	// create new gRPC server object
	grpcServer := grpc.NewServer()

	// register Implementation RPC server to gRPC server
	pb.RegisterSimpleBankServer(grpcServer, server)

	// self documentation gRPC server
	reflection.Register(grpcServer)

	// start server to listen gRPC request
	// create listener
	listener, err := net.Listen("tcp", config.GRPCServerAddress)
	if err != nil {
		log.Fatalf("Cannot create listener: %v\n", err)
	}
	// server gRPC request
	log.Printf("Start gRPC server at %s\n", listener.Addr().String())
	err = grpcServer.Serve(listener)
	if err != nil {
		log.Fatalf("Cannot start gRPC server: %v\n", err)
	}
}

func runGinServer(config util.Config, store mydb.Store) {
	server, err := server.NewServer(store, &config)
	if err != nil {
		log.Fatalf("Cannot create the server: %v\n", err)
	}

	// err = server.Start(serverAddress)
	err = server.Start(config.HTTPServerAddress)
	if err != nil {
		log.Fatalf("Cannot start the server: %v\n", err)
	}
}
