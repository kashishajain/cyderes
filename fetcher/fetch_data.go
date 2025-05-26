package fetcher

import (
	"log"
	"net/http"
	"time"
	"io"
	"errors"
)

func FetchData()([]byte, error){
	const (
		retry  = 5
		timeout = 10 * time.Second
		apiURL = "https://jsonplaceholder.typicode.com/posts"
	)

	client := &http.Client{
		Timeout: timeout,
	}
	
	var err error
	for attempts := 0; attempts < retry; attempts++ {
		req, reqErr := http.NewRequest(http.MethodGet, apiURL, nil)
		if reqErr != nil {
			log.Println("Failed to create request:", reqErr)
			return nil, reqErr
		}

		response, err := client.Do(req)
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