package main
import (
"log"
"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
"github.com/kashishajain/cyderes-app/transformer"
"github.com/kashishajain/cyderes-app/fetcher"
"github.com/kashishajain/cyderes-app/store"
)



func main() {

	data , err := fetcher.FetchData()
	if err != nil {
		log.Fatalf("Failed to fetch data even after retries, ERROR: %v", err)
	}
	transformed_data, err := transformer.TransformData(data)
	if err != nil {
		log.Fatalf("Failed to transform data: %v", err)
	}
	err = store.CreateDynamoDBTable()
	if err != nil {
		log.Fatalf("Failed to create DynamoDB table: %v", err)
	}
	if err := store.StoreToDynamoDB(transformed_data); err != nil {
		log.Fatalf("DynamoDB error: %v", err)
	}
	db_data, err := store.FetchStoredData()
    if err != nil {
        log.Fatalf("Failed to fetch stored data: %v", err)
    }
	if len(db_data) == 0 {
		log.Println("No data found in DynamoDB table.")
		return
	}
	log.Printf("Fetched %d items from DynamoDB\n", len(db_data))
	log.Println("Displaying first 10 items:")
    for i, item := range db_data {
	if i> 9 {
		break // Limit output to first 10 items
	}
    var decoded transformer.Transformed_data
    err := attributevalue.UnmarshalMap(item, &decoded)
    if err != nil {
        log.Printf("Failed to decode item %d: %v", i+1, err)
        continue
    }
    log.Printf("Item %d: %+v\n", i+1, decoded)
	}
	
}

