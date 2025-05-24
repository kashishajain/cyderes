package main

import (
    "context"
    "fmt"
    "github.com/aws/aws-sdk-go-v2/aws"
    "github.com/aws/aws-sdk-go-v2/config"
    "github.com/aws/aws-sdk-go-v2/service/dynamodb"
    "github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	
)

func FetchStoredData() ([]map[string]types.AttributeValue, error) {
    cfg, err := config.LoadDefaultConfig(context.TODO(),
        config.WithRegion("us-west-2"),
        config.WithEndpointResolver(aws.EndpointResolverFunc(
            func(service, region string) (aws.Endpoint, error) {
                return aws.Endpoint{
                    URL:           "http://dynamodb-local:8000", // DynamoDB Local endpoint
                    SigningRegion: "us-west-2",
                }, nil
            },
        )),
    )
    if err != nil {
        return nil, fmt.Errorf("unable to load SDK config: %w", err)
    }

    client := dynamodb.NewFromConfig(cfg)
    tableName := "cyderes_api_logs"

    out, err := client.Scan(context.TODO(), &dynamodb.ScanInput{
        TableName: aws.String(tableName),
    })
    if err != nil {
        return nil, fmt.Errorf("failed to scan table: %w", err)
    }

    return out.Items, nil
}