package main

import (
	"context"
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

func main() {
	/*
		By default S3 object will not be available to public. If you want to default
		to public. Please create policy generation script and use it to create policy.
		https://awspolicygen.s3.amazonaws.com/policygen.html
		https://stackoverflow.com/questions/19176926/how-to-make-all-objects-in-aws-s3-bucket-public-by-default
	*/

	err := godotenv.Load()
	endpoint := os.Getenv("ENDPOINT")
	accessKeyID := os.Getenv("ACESS_KEY_ID")
	secretAccessKey := os.Getenv("SECRET_ACCESS_KEY")
	useSSL := true

	ctx := context.Background()

	// Initialize minio client object.
	minioClient, err := minio.New(endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(accessKeyID, secretAccessKey, ""),
		Secure: useSSL,
	})
	if err != nil {
		log.Fatalln(err)
	}

	// Make a new bucket called mymusic.
	bucketName := "kaenova-test"
	location := "default"

	err = minioClient.MakeBucket(ctx, bucketName, minio.MakeBucketOptions{Region: location})
	if err != nil {
		// Check to see if we already own this bucket (which happens if you run this twice)
		exists, errBucketExists := minioClient.BucketExists(ctx, bucketName)
		if errBucketExists == nil && exists {
			log.Printf("We already own %s\n", bucketName)
		} else {
			log.Fatalln(err)
		}
	} else {
		log.Printf("Successfully created %s\n", bucketName)
	}

	objectName := "test/6b1b7f04-6424-450e-b982-8bcddba3818g.png"
	filePath := "6b1b7f04-6424-450e-b982-8bcddba3818b.png"
	contentType := "image/png"

	// Upload the zip file with FPutObject
	info, err := minioClient.FPutObject(ctx, bucketName, objectName, filePath, minio.PutObjectOptions{ContentType: contentType})
	if err != nil {
		log.Fatalln(err)
	}

	log.Printf("Successfully uploaded %s of size %d\n", objectName, info.Size)
}
