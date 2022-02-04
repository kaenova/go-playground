package main

import (
	"fmt"
	"net"
	"os"

	"github.com/kaenova/go-playground/rpc_chunk_file/service/upload"

	"github.com/joho/godotenv"
	"github.com/rs/zerolog/log"
	"google.golang.org/grpc"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Error().Msg("Dotenv is not found, using machine environment")
	}

	// Init RPC Server
	lis, err := net.Listen("tcp", fmt.Sprintf(":%s", os.Getenv("RPC_PORT")))
	if err != nil {
		log.Fatal().Str("failed to listen: %s", err.Error())
	}
	grpcServer := grpc.NewServer()
	upService := upload.Server{}
	upload.RegisterUploadServiceServer(grpcServer, &upService)
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatal().Str("failed to serve: %s", err.Error())
	}
}
