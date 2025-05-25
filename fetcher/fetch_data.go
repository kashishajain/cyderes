package fetcher

import (
	"fmt"
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
			fmt.Println("Successfully fetched data with status:", response.StatusCode)
			return io.ReadAll(response.Body)
		}
		time.Sleep(2 * time.Second) 
	}

	if err != nil {
		return nil, err
	}
	return nil, errors.New("failed to fetch data after retries")
}