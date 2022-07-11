package s3object

import (
	"context"
	"crypto/rand"
	"errors"
	"io/ioutil"
	"math/big"
	"mime/multipart"
	"net/http"
	"os"
	"strings"

	"github.com/google/uuid"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

type S3Object struct {
	client     *minio.Client
	endpoint   string
	bucketName string
	useSSL     bool
	ctx        context.Context
}

type S3ObjectI interface {
	UploadFileMultipart(file multipart.File) (objectOutput, error)
	UploadFileFromPath(filePath string) (objectOutput, error)
	DeleteObject(objectPath string) error
	GetObjectPath(fullPathEndpoint string) string
}

type objectOutput struct {
	EndpointPath string
	Path         string
	Endpoint     string
}

// Create new instace of S3 Object with MiniIO APIs
// This also will create a "./tmp" folder for uploading from memory file (multipart)
func NewS3Object(endpoint, accessKeyID, secretAcessKey, bucketName, location string, useSSL bool) (S3ObjectI, error) {
	ctx := context.Background()

	err := os.Mkdir("./tmp", 0644)
	if err != nil {
		return nil, err
	}

	client, err := minio.New(endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(accessKeyID, secretAcessKey, ""),
		Secure: useSSL,
	})
	if err != nil {
		return &S3Object{}, err
	}

	exists, err := client.BucketExists(ctx, bucketName)
	if err != nil {
		return &S3Object{}, err
	}
	if !exists {
		return &S3Object{}, errors.New("bucket not exists, please contact admin")
	}

	return &S3Object{
		client:     client,
		endpoint:   endpoint,
		bucketName: bucketName,
		useSSL:     useSSL,
		ctx:        ctx,
	}, nil
}

// Upload file using local file instances. This will generate random path to the file.
func (s *S3Object) UploadFileMultipart(file multipart.File) (objectOutput, error) {
	data, err := s.multiPartToByte(file)
	if err != nil {
		return objectOutput{}, err
	}
	tempFile, err := s.createTempFile(data)
	if err != nil {
		return objectOutput{}, err
	}

	object, err := s.UploadFileFromPath(tempFile)
	if err != nil {
		return objectOutput{}, err
	}
	err = s.deleteTempFile(tempFile)
	if err != nil {
		return objectOutput{}, err
	}
	return object, nil
}

// Upload file using local file paths. This will generate random path to the file.
func (s *S3Object) UploadFileFromPath(filePath string) (objectOutput, error) {
	data, err := ioutil.ReadFile(filePath)
	if err != nil {
		return objectOutput{}, err
	}
	objectPathName, err := s.generateObjectPathName(data)
	if err != nil {
		return objectOutput{}, err
	}
	_, err = s.client.FPutObject(s.ctx, s.bucketName, objectPathName, filePath, minio.PutObjectOptions{})
	if err != nil {
		return objectOutput{}, err
	}
	return s.createObjectOutput(objectPathName), nil
}

// Convert from "https://example.com/bucket/file.jpg" to "file.jpg"
// or from "bucket/file.jpg" to "file.jpg"
func (s *S3Object) GetObjectPath(fullPathEndpoint string) string {
	targetString := s.bucketName + "/"
	idx := strings.Index(fullPathEndpoint, targetString)
	if idx == -1 {
		return fullPathEndpoint
	}
	return fullPathEndpoint[idx+len(targetString):]
}

// Delete based on object path or full path
func (s *S3Object) DeleteObject(objectPath string) error {
	return s.client.RemoveObject(s.ctx, s.bucketName, objectPath, minio.RemoveObjectOptions{
		ForceDelete: true,
	})
}

func (s *S3Object) createObjectOutput(objectPath string) objectOutput {
	fullPath := "http"
	if s.useSSL {
		fullPath = fullPath + "s"
	}
	fullPath = fullPath + "://" + s.endpoint + "/" + s.bucketName + "/" + objectPath
	return objectOutput{
		EndpointPath: fullPath,
		Path:         objectPath,
		Endpoint:     s.endpoint,
	}
}

func (s *S3Object) generateObjectPathName(data []byte) (string, error) {
	fileName, err := s.generateFileName(data)
	if err != nil {
		return "", err
	}
	path, err := s.generatePath()
	if err != nil {
		return "", err
	}
	return path + fileName, nil
}

func (s *S3Object) createTempFile(data []byte) (string, error) {
	objectPathName, err := s.generateFileName(data)
	if err != nil {
		return "", err
	}
	tempFile := "./tmp/" + objectPathName
	err = ioutil.WriteFile(tempFile, data, 0644)
	if err != nil {
		return "", err
	}
	return tempFile, nil
}

func (s *S3Object) deleteTempFile(tempFile string) error {
	// NOTE: Probably will cause bug if there's a concurent connection on uploading
	return os.RemoveAll(tempFile)
}

func (s *S3Object) generatePath() (string, error) {
	finalPath := ""
	for i := 0; i < 4; i++ {
		path, err := s.randomString(10)
		if err != nil {
			return "", err
		}
		finalPath = finalPath + path + "/"
	}
	return finalPath, nil
}

func (s *S3Object) generateFileName(data []byte) (string, error) {
	mimeType := http.DetectContentType(data)
	fileExtension := ""
	switch mimeType {
	case "image/jpeg":
		fileExtension = fileExtension + ".jpg"
		break
	case "image/png":
		fileExtension = fileExtension + ".png"
		break
	case "application/pdf":
		fileExtension = fileExtension + ".pdf"
		break
	case "video/mp4":
		fileExtension = fileExtension + ".mp4"
		break
	case "application/octet-stream":
		return "", errors.New("unsupported file type")
	}

	return uuid.New().String() + fileExtension, nil
}

func (s *S3Object) multiPartToByte(file multipart.File) ([]byte, error) {
	var finalByte []byte
	b := make([]byte, 100)
	for {
		n, err := file.Read(b)
		finalByte = append(finalByte, b...)
		if n == 0 {
			break
		}
		if err != nil {
			return []byte{}, err
		}
	}
	return finalByte, nil
}

func (s *S3Object) randomString(n int) (string, error) {
	const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	b := make([]byte, n)
	for i := range b {
		idx, err := rand.Int(rand.Reader, big.NewInt(int64(len(letterBytes))))
		if err != nil {
			return "", err
		}
		b[i] = letterBytes[idx.Int64()]
	}
	return string(b), nil
}
