package main

import (
    "context"
    "fmt"
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
        return fmt.Errorf("unable to load SDK config: %w", err)
    }

    client := dynamodb.NewFromConfig(cfg)
    _, err = client.CreateTable(context.TODO(), &dynamodb.CreateTableInput{
        TableName: aws.String("cyderes_api_logs"),
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
        return fmt.Errorf("failed to create table: %w", err)
    }
    fmt.Println("Table created!")
    return nil
}