package upload

import (
	"bytes"
	"fmt"
	"image"
	"image/jpeg"
	"image/png"
	"io"
	"log"
	"net/http"
	"os"

	"github.com/google/uuid"
)

type Server struct {
}

func (s *Server) Upload(stream UploadService_UploadServer) error {
	var fullData []byte

	// Recieve data
	for {
		buff, err := stream.Recv()
		if err == io.EOF {
			break
		}
		fmt.Println("Recieving")
		if err != nil {
			return stream.SendAndClose(&ResBuffer{
				Status:  0,
				Message: "Fail to recieve file",
			})
		}
		fullData = append(fullData, buff.Data...)
	}
	if _, err := os.Stat("./static"); os.IsNotExist(err) {
		err := os.Mkdir("./static", 0666)
		if err != nil {
			return stream.SendAndClose(&ResBuffer{
				Status:  0,
				Message: "Cannot create static folder",
			})
		}
	}

	// Creating image from fullData (byte)
	contentType := http.DetectContentType(fullData)
	finalName := uuid.New().String()

	switch contentType {
	case "image/jpeg":
		outFile, err := os.Create("./static/" + finalName + ".jpg")
		finalName += ".jpg"
		if err != nil {
			log.Fatal(err)
		}
		img, _, err := image.Decode(bytes.NewReader(fullData))
		if err != nil {
			log.Fatal(err)
		}
		err = jpeg.Encode(outFile, img, &jpeg.Options{Quality: 100})
		if err != nil {
			log.Fatal(err)
		}
		outFile.Close()
	case "image/png":
		outFile, err := os.Create("./static/" + finalName + ".png")
		finalName += ".png"
		if err != nil {
			log.Fatal(err)
		}
		img, _, err := image.Decode(bytes.NewReader(fullData))
		if err != nil {
			log.Fatal(err)
		}
		err = png.Encode(outFile, img)
		outFile.Close()
	default:
		return stream.SendAndClose(&ResBuffer{
			Status:  0,
			Message: "Format is restricted",
		})
	}

	return stream.SendAndClose(&ResBuffer{
		Status:   1,
		Message:  "Success",
		FileName: os.Getenv("STATIC_URL") + "/" + finalName,
	})
}
