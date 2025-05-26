package fetcher

import (
	"log"
	"net/http"
	"time"
	"io"
	"errors"
)

func FetchData()([]byte, error){
	retry := 5
	api_url := "https://jsonplaceholder.typicode.com/posts"
	var response *http.Response
	var err error
	for attempts := 0; attempts < retry; attempts++ {
		response, err = http.Get(api_url)
		if err == nil && response.StatusCode == http.StatusOK {
			defer response.Body.Close()
			log.Println("Successfully fetched data with status:", response.StatusCode)
			return io.ReadAll(response.Body)
		}
		time.Sleep(2 * time.Second) 
		log.Println("Retrying after 2 seconds!! API failed with error:", err)
	}
	if err != nil {
		return nil, err
	}
	return nil, errors.New("Failed to fetch data after retries")
}