package main
import (
"log"
"fmt"
"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
"github.com/kashishajain/cyderes-app/transformer"
"github.com/kashishajain/cyderes-app/fetcher"
"github.com/kashishajain/cyderes-app/store"
)



func main() {

	data , err := fetcher.FetchData()
	if err != nil {
		log.Fatalf("Failed to fetch data: %v", err)
	}
	transformed_data, err := transformer.TransformData(data)
	if err != nil {
		log.Fatalf("Failed to transform data: %v", err)
	}
	store.CreateDynamoDBTable()
	if err := store.StoreToDynamoDB(transformed_data); err != nil {
		log.Fatalf("DynamoDB error: %v", err)
	}

	log.Println("Data stored in DynamoDB successfully.")
	db_data, err := store.FetchStoredData()
    if err != nil {
        log.Fatalf("Failed to fetch stored data: %v", err)
    }
    for i, item := range db_data {
    var decoded transformer.Transformed_data
    err := attributevalue.UnmarshalMap(item, &decoded)
    if err != nil {
        log.Printf("Failed to decode item %d: %v", i+1, err)
        continue
    }
    fmt.Printf("Item %d: %+v\n", i+1, decoded)
	}
	
}

