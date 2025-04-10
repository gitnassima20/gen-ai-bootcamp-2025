package services

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/bedrockruntime"
)

type BedrockService struct {
	Client *bedrockruntime.Client
	Model  string
}

func NewBedrockService(client *bedrockruntime.Client, model string) *BedrockService {
	return &BedrockService{
		Client: client,
		Model:  model,
	}
}

func (b *BedrockService) Invoke(prompt string) (string, error) {
	payload, err := json.Marshal(map[string]interface{}{
		"inferenceConfig": map[string]int{
			"max_new_tokens": 1000,
		},
		"messages": []map[string]interface{}{
			{
				"role": "user",
				"content": []map[string]string{
					{
						"text": prompt,
					},
				},
			},
		},
	})
	if err != nil {
		return "", fmt.Errorf("failed to marshal request: %v", err)
	}

	resp, err := b.Client.InvokeModel(context.Background(), &bedrockruntime.InvokeModelInput{
		ModelId:     aws.String(b.Model),
		ContentType: aws.String("application/json"),
		Accept:      aws.String("application/json"),
		Body:        payload,
	})
	if err != nil {
		return "", fmt.Errorf("failed to invoke model: %v", err)
	}

	var result map[string]interface{}
	if err := json.Unmarshal(resp.Body, &result); err != nil {
		return "", fmt.Errorf("failed to decode response: %v", err)
	}

	var fullResponse string
	if body, ok := result["body"].(map[string]interface{}); ok {
		if text, ok := body["text"].(string); ok {
			fullResponse = text
		}
	}

	if fullResponse == "" {
		return "", fmt.Errorf("no response content received")
	}

	return fullResponse, nil
}
