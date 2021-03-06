package main

import (
	"fmt"
	"net"
	"os"

	"github.com/kaenova/go-playground/rpc_chunk_file/service/upload"

	"github.com/joho/godotenv"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"google.golang.org/grpc"
)

func main() {
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})

	err := godotenv.Load()
	if err != nil {
		log.Error().Msg("Dotenv is not found, using machine environment")
	}

	// Init static folder if there's none
	if _, err := os.Stat("./static"); os.IsNotExist(err) {
		err := os.Mkdir("./static", 0666)
		if err != nil {
			log.Fatal().Msg("Cannot create static folder")
		}
	}

	// Init RPC Server
	lis, err := net.Listen("tcp", fmt.Sprintf(":%s", os.Getenv("RPC_PORT")))
	if err != nil {
		log.Fatal().Str("failed to listen: %s", err.Error())
	}
	grpcServer := grpc.NewServer()
	upService := upload.Server{}
	upload.RegisterUploadServiceServer(grpcServer, &upService)
	log.Info().Msg("Listening on :" + os.Getenv("RPC_PORT"))
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatal().Str("failed to serve: %s", err.Error())
	}
}
