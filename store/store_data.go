package store

import (
	"context"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/kashishajain/cyderes-app/transformer"
)

func StoreToDynamoDB(data []transformer.Transformed_data) error {
	cfg, err := config.LoadDefaultConfig(context.TODO(),
		config.WithRegion("us-west-2"),
		// config.WithEndpointResolver(aws.EndpointResolverFunc(
		// 	func(service, region string) (aws.Endpoint, error) {
		// 		return aws.Endpoint{
		// 			URL:           "http://dynamodb-local:8000", // DynamoDB Local endpoint
		// 			SigningRegion: "us-west-2",
		// 		}, nil
		// 	},
		// )),
	)
	if err != nil {
		return fmt.Errorf("unable to load SDK config: %w", err)
	}

	client := dynamodb.NewFromConfig(cfg)
	tableName := "cyderes_api_logs"

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
			return fmt.Errorf("failed to put item: %w", err)
		}
	}

	return nil
}
