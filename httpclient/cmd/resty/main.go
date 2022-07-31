package main

import (
	"bytes"
	"crypto/tls"
	"io/ioutil"
	"log"
	"net/http"
	"strings"

	"github.com/go-resty/resty/v2"
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
	fileBytes, err := ioutil.ReadFile(FileName)
	if err != nil {
		log.Fatal(err.Error())
	}

	client := resty.New()
	client.SetDisableWarn(true)
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	client.SetTransport(tr)
	resp, err := client.R().
		SetMultipartFields(
			&resty.MultipartField{
				Param:       "pdf_article",
				FileName:    "test.pdf",
				ContentType: "application/pdf",
				Reader:      bytes.NewReader(fileBytes),
			},
			&resty.MultipartField{
				Param:       "paper_id",
				FileName:    "",
				ContentType: "text/plain",
				Reader:      strings.NewReader(`apapun`),
			}).
		SetContentLength(true).
		Post(URL)

	if err != nil {
		log.Fatal(err.Error())
	}
	log.Println(string(resp.Body()))
}
