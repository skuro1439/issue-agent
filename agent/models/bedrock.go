package models

import (
	"context"
	"encoding/json"
	"fmt"
	"net/url"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/bedrockruntime"

	"github.com/clover0/issue-agent/logger"
)

type BedrockClient struct {
	client  *bedrockruntime.Client
	logger  logger.Logger
	baseURL *url.URL

	// services
	Messages *BedrockMessageService
}

func NewBedrock(logger logger.Logger) (BedrockClient, error) {
	ctx := context.Background()
	var opts []func(*config.LoadOptions) error
	sdkConfig, err := config.LoadDefaultConfig(ctx, opts...)
	if err != nil {
		return BedrockClient{}, fmt.Errorf("failed to load AWS SDK config: %w", err)
	}

	client := bedrockruntime.NewFromConfig(sdkConfig)

	c := BedrockClient{
		logger: logger,
		client: client,
	}

	c.Messages = &BedrockMessageService{client: &c}

	return c, nil
}

type BedrockMessageService struct {
	client *BedrockClient
}

func (s *BedrockMessageService) Create(ctx context.Context, body J) (ResponseMessage, error) {
	b, err := json.Marshal(body)
	if err != nil {
		return ResponseMessage{}, fmt.Errorf("failed to marshal body: %w", err)
	}

	result, err := s.client.client.InvokeModel(ctx, &bedrockruntime.InvokeModelInput{
		// todo: changeable models
		ModelId:     aws.String("anthropic.claude-3-5-sonnet-20240620-v1:0"),
		ContentType: aws.String("application/json"),
		Body:        b,
	})
	if err != nil {
		errMsg := err.Error()
		if strings.Contains(errMsg, "provided model identifier is invalid") {
			return ResponseMessage{}, fmt.Errorf("failed to invoke model: %w: hint - check whether enabled the model and in the AWS region", err)
		}
		return ResponseMessage{}, fmt.Errorf("failed to invoke model: %w", err)
	}

	var message ResponseMessage
	if err := json.Unmarshal(result.Body, &message); err != nil {
		return message, fmt.Errorf("failed to unmarshal response body: %w", err)
	}

	return message, nil
}
