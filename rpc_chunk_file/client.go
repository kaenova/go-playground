/*
Di file ini ditunjukkan kalau kita mempunyai file reader.

Pada kasus ini saya harus membaca suatu array tersebut sebesar maxSizeChunk,
mengirimkan ke server, dan membaca lagi sebesar maxSizeChunk, hingga yang dibaca
lagi tidak ada.
*/

package main

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"os"
	"test-go/rpc_chunk_file/service/upload"

	"github.com/joho/godotenv"
	"github.com/rs/zerolog/log"
	"google.golang.org/grpc"
)

func main() {
	const maxSizeChunk = 2 * 1024 // 2kB

	// Setting up gRPC connection
	err := godotenv.Load()
	if err != nil {
		log.Error().Msg("Dotenv is not found, using machine environment")
	}
	var conn *grpc.ClientConn
	conn, err = grpc.Dial(":"+os.Getenv("RPC_PORT"), grpc.WithInsecure())
	if err != nil {
		log.Fatal().Err(err)
	}
	defer conn.Close()
	c := upload.NewUploadServiceClient(conn)
	svc, err := c.Upload(context.Background())
	if err != nil {
		log.Fatal().Err(err)
	}

	// Setup file reader
	f, err := os.Open("./test.png")
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	// Setup chunk buffer
	r := bufio.NewReader(f)
	buf := make([]byte, 0, maxSizeChunk)

	for {
		// Read file with a chunksize
		n, err := r.Read(buf[:cap(buf)])
		buf = buf[:n]
		if n == 0 {
			if err == nil {
				continue
			}
			if err == io.EOF {
				break
			}
			log.Fatal().Err(err)
		}
		err = svc.Send(&upload.ReqBuffer{
			FileName: "Testing",
			Data:     buf,
		})
		if err != nil && err != io.EOF {
			log.Fatal().Err(err)
			break
		}
	}

	// Get final Data from Server
	finalRes, err := svc.CloseAndRecv()
	if err != nil {
		log.Fatal().Err(err)
	}

	fmt.Println(finalRes.Status)
	fmt.Println(finalRes.Message)
	fmt.Println(finalRes.FileName)
}
