/*
Di file ini ditunjukkan kalau kita punya sebuah variable byte penuh yang berisi
suatu gambar.

Pada kasus ini saya harus membagi-bagi array tersebut sebesar maxSizeChunk, lalu
mengirimkan ke server
*/

package main

import (
	"context"
	"fmt"
	"io/ioutil"
	"math"
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

	// Read full image into a buffer
	buff, err := ioutil.ReadFile("./test.png")
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	// Sending to Server
	counter := 0
	var minIdx, maxIdx int
	for {
		minIdx = maxIdx
		maxIdx = int(math.Min(float64(minIdx+maxSizeChunk), float64(len(buff))))

		// Get Byte Data from minimumChunkIdx to maximumChunkIdx
		buf := buff[minIdx:maxIdx]
		err = svc.Send(&upload.ReqBuffer{
			FileName: "Testing",
			Data:     buf,
		})
		if err != nil {
			log.Fatal().Err(err)
			break
		}
		counter++

		// If there's no next byte to send
		if maxIdx+1 >= len(buff)-1 {
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
