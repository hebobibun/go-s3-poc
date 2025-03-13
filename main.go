package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/feature/s3/manager"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"gopkg.in/yaml.v3"
)

// Config struct for YAML
type Config struct {
	AWS struct {
		AccessKeyID     string `yaml:"access_key_id"`
		SecretAccessKey string `yaml:"secret_access_key"`
		Region          string `yaml:"region"`
		Bucket          string `yaml:"bucket"`
	} `yaml:"aws"`
}

// Load YAML config
func loadConfig(filename string) (*Config, error) {
	data, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	var cfg Config
	err = yaml.Unmarshal(data, &cfg)
	if err != nil {
		return nil, err
	}

	return &cfg, nil
}

func main() {
	var awsCfg aws.Config
	var err error
	var bucketName string

	// Try to load YAML config
	cfg, yamlErr := loadConfig("config.yaml")
	if yamlErr == nil {
		// If YAML exists, use it
		os.Setenv("AWS_ACCESS_KEY_ID", cfg.AWS.AccessKeyID)
		os.Setenv("AWS_SECRET_ACCESS_KEY", cfg.AWS.SecretAccessKey)
		os.Setenv("AWS_REGION", cfg.AWS.Region)
		bucketName = cfg.AWS.Bucket
		fmt.Println("✅ Loaded AWS credentials from config.yaml")
	}

	// Load AWS SDK config (fallback to CLI credentials if YAML is missing)
	awsCfg, err = config.LoadDefaultConfig(context.TODO(), config.WithRegion(cfg.AWS.Region))
	if err != nil {
		log.Fatalf("❌ Failed to load AWS credentials: %v", err)
	}

	// Create S3 client
	s3Client := s3.NewFromConfig(awsCfg)

	// Upload a file
	filePath := "test.txt"
	err = uploadFile(context.TODO(), s3Client, bucketName, filePath)
	if err != nil {
		log.Fatalf("❌ Upload failed: %v", err)
	}

	fmt.Println("✅ File uploaded successfully to AWS S3!")
}

// Upload file to S3
func uploadFile(ctx context.Context, client *s3.Client, bucket, filePath string) error {
	uploader := manager.NewUploader(client)

	// Open the file
	file, err := os.Open(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	// Upload
	_, err = uploader.Upload(ctx, &s3.PutObjectInput{
		Bucket: &bucket,
		Key:    &filePath,
		Body:   file,
	})
	return err
}
