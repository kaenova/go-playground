package upload

import (
	"errors"
	"io"
	"io/ioutil"
	"net/http"
	"os"

	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
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
		if err != nil {
			return err
		}
		fullData = append(fullData, buff.Data...)
	}
	if _, err := os.Stat("./static"); os.IsNotExist(err) {
		err := os.Mkdir("./static", 0666)
		if err != nil {
			return err
		}
	}

	// Creating image from fullData (byte)
	contentType := http.DetectContentType(fullData)
	finalName := uuid.New().String()

	switch contentType {
	case "image/jpeg":
		finalName += ".jpg"
	case "image/png":
		finalName += ".png"
	default:
		log.Info().Msg("Got into this")
		return errors.New("format is not permitted")
	}

	err := ioutil.WriteFile("./static/"+finalName, fullData, 0444)
	if err != nil {
		return err
	}

	log.Info().Msg("Successfull saving a file")
	return stream.SendAndClose(&ResBuffer{
		FileName: os.Getenv("STATIC_URL") + "/" + finalName,
	})
}
