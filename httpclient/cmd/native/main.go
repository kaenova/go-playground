package main

import (
	"bytes"
	"crypto/tls"
	"fmt"
	"io/ioutil"
	"log"
	"mime/multipart"
	"net/http"
)

const (
	URL         = "https://ir-group.ec.tuwien.ac.at/artu_az_identification/identify_az"
	ContentType = "multipart/form-data"
	FileName    = "test.pdf"
	// paper_id string
	// pdf_article file
)

func main() {
	// Read file
	file, err := ioutil.ReadFile(FileName)
	if err != nil {
		log.Fatal(err.Error())
	}

	// Prepare file multipart
	body := new(bytes.Buffer)
	multiPartBody := multipart.NewWriter(body)
	bodyFile, err := multiPartBody.CreateFormFile("pdf_article", FileName)
	if err != nil {
		log.Fatal(err.Error())
	}
	_, err = bodyFile.Write(file)
	if err != nil {
		log.Fatal(err.Error())
	}

	// Prepare string field
	otherField := map[string]string{
		"paper_id": "apapun",
	}
	for param, value := range otherField {
		tempWriter, err := multiPartBody.CreateFormField(param)
		if err != nil {
			log.Fatal(err.Error())
		}
		tempWriter.Write([]byte(value))
	}

	log.Println(multiPartBody.FormDataContentType())
	req, err := http.NewRequest("POST", URL, body)
	req.Header.Add("Content-Length", fmt.Sprint(body.Len()))
	req.Header.Add("Content-Type", multiPartBody.FormDataContentType())
	req.Header.Add("Host", "ir-group.ec.tuwien.ac.at")
	for k, v := range req.Header {
		log.Println(k, v)
	}
	if err != nil {
		log.Fatal(err.Error())
	}
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	client := &http.Client{Transport: tr}
	resp, err := client.Do(req)
	if err != nil {
		log.Fatal(err.Error())
	}

	// read response body
	var a []byte
	a, error := ioutil.ReadAll(resp.Body)
	if error != nil {
		fmt.Println(error)
	}

	// print response body
	fmt.Println(string(a))
}
