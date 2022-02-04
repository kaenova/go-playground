package main

import (
	"bufio"
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

func main() {
	var finalData []byte

	f, err := os.Open("./test2.png")
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	nBytes, nChunks := int64(0), int64(0)
	r := bufio.NewReader(f)
	buf := make([]byte, 0, 4*1024)
	for {
		n, err := r.Read(buf[:cap(buf)])
		buf = buf[:n]
		if n == 0 {
			if err == nil {
				continue
			}
			if err == io.EOF {
				break
			}
			log.Fatal(err)
		}
		nChunks++
		nBytes += int64(len(buf))
		finalData = append(finalData, buf...)
		if err != nil && err != io.EOF {
			log.Fatal(err)
		}
	}

	fmt.Println("Bytes:", nBytes, "Chuncks:", nChunks)

	// Creating image from finalData (byte)
	contentType := http.DetectContentType(finalData)
	finalName := uuid.New().String()

	if _, err := os.Stat("./static"); os.IsNotExist(err) {
		err := os.Mkdir("./static", 0666)
		if err != nil {
			log.Fatal(err)
		}
	}

	switch contentType {
	case "image/jpeg":
		outFile, err := os.Create("./static/" + finalName + ".jpg")
		if err != nil {
			log.Fatal(err)
		}
		img, _, err := image.Decode(bytes.NewReader(finalData))
		if err != nil {
			log.Fatal(err)
		}
		err = jpeg.Encode(outFile, img, &jpeg.Options{Quality: 100})
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println("Successfully written a file")
	case "image/png":
		outFile, err := os.Create("./static/" + finalName + ".png")
		if err != nil {
			log.Fatal(err)
		}
		img, _, err := image.Decode(bytes.NewReader(finalData))
		if err != nil {
			log.Fatal(err)
		}
		err = png.Encode(outFile, img)
		fmt.Println("Successfully written a file")
		break
	}
}
