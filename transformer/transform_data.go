package transformer

import (
	"encoding/json"
	"time"
	"log"
)

type Original_data struct {
    UserID int    `json:"userId" dynamodbav:"userId"`
    ID     int    `json:"id" dynamodbav:"id"`
    Title  string `json:"title" dynamodbav:"title"`
    Body   string `json:"body" dynamodbav:"body"`
}

type Transformed_data struct {
    UserID int    `json:"userId" dynamodbav:"userId"`
    ID     int    `json:"id" dynamodbav:"id"`
    Title  string `json:"title" dynamodbav:"title"`
    Body   string `json:"body" dynamodbav:"body"`
	IngestedAt string `json:"ingested_at" dynamodbav:"ingested_at"`
	Source string `json:"source" dynamodbav:"source"`
}

func TransformData(data []byte) ([]Transformed_data, error){

	var orig_data []Original_data
	if err := json.Unmarshal(data, &orig_data); err != nil {
    return nil, err
	}

	var transformed_data []Transformed_data
	source := "placeholder_api"
	ingested_time := time.Now().UTC().Format(time.RFC3339)
	for _, d := range orig_data{
		td := Transformed_data{
			UserID:     d.UserID,
            ID:         d.ID,
            Title:      d.Title,
            Body:       d.Body,
			IngestedAt: ingested_time,
			Source: source,
		}
		transformed_data = append(transformed_data, td)
	}
	log.Println("Successfully transformed data. Added UTC stamp, source string in Original Data.")
	return transformed_data, nil

	
}