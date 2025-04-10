package config

import (
	"context"
	"log"
	"os"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/joho/godotenv"
)

type AWSConfig struct {
	AccessKey string
	SecretKey string
	Region    string
}

func LoadAWSConfig() aws.Config {
	// Load .env file
	err := godotenv.Load()
	if err != nil {
		log.Println("No .env file found or failed to load it")
	}

	// Get AWS credentials from environment
	accessKey := os.Getenv("AWS_ACCESS_KEY_ID")
	secretKey := os.Getenv("AWS_SECRET_ACCESS_KEY")
	region := os.Getenv("AWS_REGION")

	// Validate required fields
	if accessKey == "" || secretKey == "" || region == "" {
		log.Fatal("Missing required AWS configuration")
	}

	// Create AWS configuration
	cfg, err := config.LoadDefaultConfig(context.TODO(),
		config.WithRegion(region),
		config.WithCredentialsProvider(
			credentials.NewStaticCredentialsProvider(accessKey, secretKey, ""),
		),
	)
	if err != nil {
		log.Fatalf("Unable to load AWS SDK config: %v", err)
	}

	return cfg
}
