package examples

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"time"

	ds "github.com/soumitsalman/beansack/sdk"
	"github.com/soumitsalman/newscollector/collector"
)

const _SITEMAPS = "./examples/sitemaps.csv"

func StoreLocal() {
	start_time := time.Now()
	// initialize to save locally
	collector := collector.NewCollector(_SITEMAPS, localFileStore)
	collector.Collect()
	log.Println("Collection took", time.Since(start_time))
}

func localFileStore(contents []ds.Bean) {
	if len(contents) > 0 {
		data, _ := json.MarshalIndent(contents, "", "\t")
		filename := fmt.Sprintf("test_%s_%s", contents[0].Source, time.Now().Format("2006-01-02-15-04-05.json"))
		os.WriteFile(filename, data, 0644)
	}
}

// func StoreRemote() {
// 	start_time := time.Now()
// 	// initialize to save locally
// 	collector := collector.NewCollector(_SITEMAPS, remoteStoreBeans)
// 	collector.Collect()
// 	log.Println("Collection took", time.Since(start_time))
// }

// const (
// 	_JSON_BODY   = "application/json"
// 	_MAX_TIMEOUT = 10 * time.Minute
// 	_USER_AGENT  = "Web Beans"
// )

// func remoteStoreBeans(contents []ds.Bean) {
// 	// debug_writeJsonFile(contents)
// 	_, err := getMediaStoreClient().R().
// 		SetHeader("Content-Type", _JSON_BODY).
// 		SetBody(contents).
// 		Put("/beans")
// 	if err != nil {
// 		log.Println("FAILED storing new contents", err)
// 	}
// }

// var bean_sack_client *resty.Client

// func getMediaStoreClient() *resty.Client {
// 	if bean_sack_client == nil {
// 		bean_sack_client = resty.New().
// 			SetTimeout(_MAX_TIMEOUT).
// 			SetBaseURL(os.Getenv("BEAN_SACK_URL")).
// 			SetHeader("User-Agent", _USER_AGENT).
// 			SetHeader("X-API-Key", os.Getenv("INTERNAL_AUTH_TOKEN"))
// 	}
// 	return bean_sack_client
// }
