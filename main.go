package main

import (
	"io"
	"log"
	"net/http"
	"os"
)

// given a url can I download the HTML content of this page

var sampleURl = "http://kenyaseed.com"

func main() {
	print("seed crawler init")

	file, err := os.Create("out.html")
	if err != nil {
		log.Fatal(err)
	}

	defer file.Close()

	c := client()
	res, err := c.Get(sampleURl)
	if err != nil {
		log.Fatal(err)
	}

	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		log.Fatal(err)
	}

	_, err = file.Write(body)
	if err != nil {
		log.Fatal(err)
	}

	// fmt.Println(fn)
}

func client() *http.Client {
	tr := &http.Transport{
		IdleConnTimeout: 10,
	}

	return &http.Client{
		Transport: tr,
	}
}
