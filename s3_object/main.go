package main

import (
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/kaenova/go-playground/s3_object/s3object"
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

	// Make a new bucket called mymusic.
	bucketName := "s3.kaenova.my.id"
	// location := "default"

	s3, err := s3object.NewS3Object(endpoint, accessKeyID, secretAccessKey, bucketName, useSSL)
	if err != nil {
		log.Fatal(err)
	}

	// app := fiber.New()

	// app.Post("/", func(c *fiber.Ctx) error {

	// 	if form, err := c.MultipartForm(); err == nil {

	// 		// Get all files from "documents" key:
	// 		files := form.File["photo"]
	// 		// => []*multipart.FileHeader

	// 		// Loop through files:
	// 		for _, file := range files {
	// 			file, err := file.Open()
	// 			if err != nil {
	// 				return err
	// 			}
	// 			object, err := s3.UploadFileMultipart(file)
	// 			if err != nil {
	// 				return err
	// 			}
	// 			return c.JSON(fiber.Map{
	// 				"path": object.EndpointPath,
	// 			})
	// 		}
	// 	}

	// 	return err
	// })

	// app.Listen(":3000")

	// Testing for local upload

	// a, err := s3.UploadFileFromPath("./6b1b7f04-6424-450e-b982-8bcddba3818b.png")
	// if err != nil {
	// 	log.Fatal(err)
	// }

	b, err := s3.GetObjectPresigned("Datavidia/kaenova_lstms.csv")

	log.Println(b)

	// objs := s3.ListObjectParentDir()
	// for _, obj := range objs {
	// 	s3.DeleteObject(obj)
	// }

}
