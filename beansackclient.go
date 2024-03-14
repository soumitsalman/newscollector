package main

import (
	"log"
	"os"
	"time"

	"github.com/go-resty/resty/v2"
	"github.com/soumitsalman/document-loader/document"
)

const (
	_JSON_BODY   = "application/json"
	_MAX_TIMEOUT = 15 * time.Second
)

func storeNewBeans(contents []*document.Document) {
	// debug_writeJsonFile(contents)
	_, err := getMediaStoreClient().R().
		SetHeader("Content-Type", _JSON_BODY).
		SetBody(contents).
		Post("/contents")
	if err != nil {
		log.Println("FAILED storing new contents", err)
	}
}

var bean_sack_client *resty.Client

func getMediaStoreClient() *resty.Client {
	if bean_sack_client == nil {
		bean_sack_client = resty.New().
			SetTimeout(_MAX_TIMEOUT).
			SetBaseURL(os.Getenv("BEAN_SACK_URL")).
			SetHeader("User-Agent", os.Getenv("Web Beans")).
			SetHeader("X-API-Key", os.Getenv("INTERNAL_AUTH_TOKEN"))
	}
	return bean_sack_client
}
