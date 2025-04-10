package services

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/bedrockruntime"
)

type BedrockService struct {
	Client  *bedrockruntime.Client
	ModelID string
	Timeout time.Duration
}

func NewBedrockService(client *bedrockruntime.Client, modelID string) *BedrockService {
	return &BedrockService{
		Client:  client,
		ModelID: modelID,
		Timeout: 30 * time.Second,
	}
}

func (s *BedrockService) Invoke(prompt string) (string, error) {
	log.Printf("Bedrock Service: Invoking model %s", s.ModelID)
	log.Printf("Prompt: %s", prompt)

	// Prepare context with timeout
	ctx, cancel := context.WithTimeout(context.Background(), s.Timeout)
	defer cancel()

	// Prepare payload for Nova-Lite model
	payload := map[string]interface{}{
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
	}

	// Convert payload to JSON
	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		log.Printf("Bedrock Service: JSON marshaling error - %v", err)
		return "", fmt.Errorf("failed to prepare payload: %v", err)
	}

	// Log payload for debugging
	log.Printf("Payload: %s", string(payloadBytes))

	// Invoke Bedrock model
	input := &bedrockruntime.InvokeModelInput{
		ModelId:     &s.ModelID,
		Body:        payloadBytes,
		ContentType: aws.String("application/json"),
	}

	output, err := s.Client.InvokeModel(ctx, input)
	if err != nil {
		log.Printf("Bedrock Service: Model invocation error - %v", err)
		return "", fmt.Errorf("model invocation failed: %v", err)
	}

	// Log raw response for debugging
	log.Printf("Raw Response: %s", string(output.Body))

	// Parse response
	var response map[string]interface{}
	err = json.Unmarshal(output.Body, &response)
	if err != nil {
		log.Printf("Bedrock Service: Response parsing error - %v", err)
		return "", fmt.Errorf("failed to parse response: %v", err)
	}

	// Extract text from response
	text, err := extractResponseText(response)
	if err != nil {
		log.Printf("Bedrock Service: Text extraction error - %v", err)
		return "", err
	}

	log.Printf("Bedrock Service: Extracted response - %s", text)
	return text, nil
}

func extractResponseText(response map[string]interface{}) (string, error) {
	// Try different possible response structures for Bedrock
	paths := [][]string{
		{"output", "message", "content", "0", "text"},
		{"output", "message", "content", "text"},
		{"output", "text"},
		{"content", "0", "text"},
		{"generation"},
		{"text"},
		{"body", "text"},
	}

	for _, path := range paths {
		var current interface{} = response
		for _, key := range path {
			if m, ok := current.(map[string]interface{}); ok {
				current = m[key]
			} else if arr, ok := current.([]interface{}); ok && len(arr) > 0 {
				current = arr[0]
			} else {
				break
			}
		}

		if text, ok := current.(string); ok {
			// Clean up the response
			text = strings.TrimSpace(text)
			if text != "" {
				return text, nil
			}
		}
	}

	// Log the full response for debugging
	log.Printf("Could not extract text. Full response: %+v", response)
	return "", fmt.Errorf("could not extract text from response")
}
