package store

import (
    "context"
    "log"
    "fmt"
    "time"
    "github.com/aws/aws-sdk-go-v2/aws"
    "github.com/aws/aws-sdk-go-v2/config"
    "github.com/aws/aws-sdk-go-v2/service/dynamodb"
    "github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	
)

func CreateDynamoDBTable() error {
    cfg, err := config.LoadDefaultConfig(context.TODO(),
        config.WithRegion("us-west-2"),
        // config.WithEndpointResolver(aws.EndpointResolverFunc(
        //     func(service, region string) (aws.Endpoint, error) {
        //         return aws.Endpoint{
        //             URL:           "http://dynamodb-local:8000", // DynamoDB Local endpoint
        //             SigningRegion: "us-west-2",
        //         }, nil
        //     },
        // )),
    )
    if err != nil {
        return fmt.Errorf("Unable to load SDK config: %w", err)
    }
    tableName := "cyderes_api_logs"
    client := dynamodb.NewFromConfig(cfg)

     _, err = client.DescribeTable(context.TODO(), &dynamodb.DescribeTableInput{
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
    time.Sleep(5 * time.Second) 
    log.Println("Table Successfully created in Database: ", tableName) 
    return nil
}