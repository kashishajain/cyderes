package integration

import (
    "context"
    "os"
    "testing"
    "time"
    "fmt"
    "log"
    "github.com/aws/aws-sdk-go-v2/aws"
    "github.com/aws/aws-sdk-go-v2/config"
    "github.com/aws/aws-sdk-go-v2/service/dynamodb"
    "github.com/aws/aws-sdk-go-v2/service/dynamodb/types"

    "github.com/stretchr/testify/assert"
    "github.com/kashishajain/cyderes-app/transformer"
    "github.com/kashishajain/cyderes-app/fetcher"
)

func setupTestTable(t *testing.T, client *dynamodb.Client, tableName string) error {
    _, err := client.DescribeTable(context.TODO(), &dynamodb.DescribeTableInput{
        TableName: aws.String(tableName),
    })
    if err == nil {
        log.Println("Table already exists:", tableName)
        return nil
    }

    _, err = client.CreateTable(context.TODO(), &dynamodb.CreateTableInput{
        TableName: aws.String(tableName),
        AttributeDefinitions: []types.AttributeDefinition{
            {
                AttributeName: aws.String("userId"),
                AttributeType: types.ScalarAttributeTypeN,
            },
            {
                AttributeName: aws.String("id"),
                AttributeType: types.ScalarAttributeTypeN,
            },
        },
        KeySchema: []types.KeySchemaElement{
            {
                AttributeName: aws.String("userId"),
                KeyType:       types.KeyTypeHash, // Partition key
            },
            {
                AttributeName: aws.String("id"),
                KeyType:       types.KeyTypeRange, // Sort key
            },
        },
        ProvisionedThroughput: &types.ProvisionedThroughput{
            ReadCapacityUnits:  aws.Int64(1),
            WriteCapacityUnits: aws.Int64(1),
        },
    })
    if err != nil {
        return fmt.Errorf("Failed to create dynamodb table: %w", err)
    }
    log.Println("Waiting for table to be created...")
    time.Sleep(3 * time.Second) 
    log.Println("Table Successfully created in Database:", tableName) 
    return nil
}

func StoreToDynamoDB(t *testing.T, data []transformer.Transformed_data, client *dynamodb.Client, tableName string) error {
	for _, item := range data {
		_, err := client.PutItem(context.TODO(), &dynamodb.PutItemInput{
			TableName: aws.String(tableName),
			Item: map[string]types.AttributeValue{
				"id":          &types.AttributeValueMemberN{Value: fmt.Sprintf("%d", item.ID)},
				"userId":      &types.AttributeValueMemberN{Value: fmt.Sprintf("%d", item.UserID)},
				"title":       &types.AttributeValueMemberS{Value: item.Title},
				"body":        &types.AttributeValueMemberS{Value: item.Body},
				"ingested_at": &types.AttributeValueMemberS{Value: item.IngestedAt},
				"source":      &types.AttributeValueMemberS{Value: item.Source},
			},
		})

		if err != nil {
			return fmt.Errorf("Failed to put item: %w", err)
		}
	}
	log.Println("Data stored in DynamoDB successfully.")
	return nil
}

func TestIngestDataToDynamoDB(t *testing.T) {
    os.Setenv("AWS_ACCESS_KEY_ID", "dummy")
    os.Setenv("AWS_SECRET_ACCESS_KEY", "dummy")
    os.Setenv("AWS_REGION", "us-west-2")
    os.Setenv("DYNAMO_ENDPOINT", "http://localhost:8000")
    testTableName := "test_cyderes_api_logs"

    cfg, err := config.LoadDefaultConfig(context.TODO(),
        config.WithEndpointResolverWithOptions(
            aws.EndpointResolverWithOptionsFunc(func(service, region string, _ ...interface{}) (aws.Endpoint, error) {
                if service == dynamodb.ServiceID {
                    return aws.Endpoint{URL: "http://localhost:8000", SigningRegion: "us-west-2"}, nil
                }
                return aws.Endpoint{}, nil
            }),
        ),
    )
    assert.NoError(t, err)

    dynamoClient := dynamodb.NewFromConfig(cfg)

    

    // inject test table into ingestion logic
    data , err := fetcher.FetchData()
	if err != nil {
		log.Fatalf("Failed to fetch data even after retries, ERROR: %v", err)
	}
	transformed_data, err := transformer.TransformData(data)
	if err != nil {
		log.Fatalf("Failed to transform data: %v", err)
	}
	err = setupTestTable(t, dynamoClient, testTableName)
    if err != nil {
		log.Fatalf("Failed to create DynamoDB table: %v", err)
	}
	if err := StoreToDynamoDB(t, transformed_data, dynamoClient, testTableName); err != nil {
		log.Fatalf("DynamoDB error: %v", err)
	}
    assert.NoError(t, err)

    // check that at least one item was ingested
    out, err := dynamoClient.Scan(context.TODO(), &dynamodb.ScanInput{
        TableName: aws.String(testTableName),
    })
    assert.NoError(t, err)
    assert.Greater(t, len(out.Items), 0)

    // Optionally, validate transformed fields
    assert.Contains(t, out.Items[0], "ingested_at")
    assert.Equal(t, "placeholder_api" , out.Items[0]["source"].(*types.AttributeValueMemberS).Value)
}
