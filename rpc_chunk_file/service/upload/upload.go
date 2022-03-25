package upload

import (
	"io"
	"io/ioutil"
	"net/http"
	"os"

	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
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
			return status.Error(codes.DataLoss, "Dataloss when recieving bytes")
		}
		fullData = append(fullData, buff.Data...)
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
		return status.Error(codes.PermissionDenied, "Format is not permitted")
	}

	err := ioutil.WriteFile("./static/"+finalName, fullData, 0444)
	if err != nil {
		return status.Error(codes.Internal, "Cannot create file")
	}

	log.Info().Msg("Successfull saving a file")
	return stream.SendAndClose(&ResBuffer{
		FileName: os.Getenv("STATIC_URL") + "/" + finalName,
	})
}
