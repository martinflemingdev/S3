package main

import (
    "context"
    "fmt"
    "io"
    "log"
	"bytes"
	"os"

    "github.com/aws/aws-sdk-go-v2/aws"
    "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
    "github.com/aws/aws-sdk-go-v2/service/s3"
)

// GetObjectFromS3 retrieves an object from an S3 bucket using AWS SDK for Go v2.
func GetObjectFromS3(ctx context.Context, s3Client *s3.Client, bucket, key string) ([]byte, error) {
    // Create the input configuration instance
    input := &s3.GetObjectInput{
        Bucket: aws.String(bucket),
        Key:    aws.String(key),
    }

    // Perform the GetObject operation
    result, err := s3Client.GetObject(ctx, input)
    if err != nil {
        return nil, fmt.Errorf("failed to get object: %v", err)
    }
    defer result.Body.Close()

    // Read the body of the response
    body, err := io.ReadAll(result.Body)
    if err != nil {
        return nil, fmt.Errorf("failed to read object body: %v", err)
    }

    return body, nil
}

// PutObjectOnS3 uploads an object to an S3 bucket using AWS SDK for Go v2.
func PutObjectOnS3(ctx context.Context, s3Client *s3.Client, bucket, key string, content []byte) error {
    // Create the input configuration instance
    input := &s3.PutObjectInput{
        Bucket: aws.String(bucket),
        Key:    aws.String(key),
        Body:   bytes.NewReader(content),
    }

    // Perform the PutObject operation
    _, err := s3Client.PutObject(ctx, input)
    if err != nil {
        return fmt.Errorf("failed to put object: %v", err)
    }

    return nil
}

// NewS3ClientWithIAMUserCreds creates an S3 client using IAM user credentials from environment variables.
func NewS3ClientWithIAMUserCreds(ctx context.Context) (*s3.Client, error) {
    accessKeyID := os.Getenv("AWS_ACCESS_KEY_ID")
    secretAccessKey := os.Getenv("AWS_SECRET_ACCESS_KEY")
    region := os.Getenv("AWS_REGION") // Region is also fetched from an environment variable

    cfg, err := config.LoadDefaultConfig(ctx,
        config.WithRegion(region),
        config.WithCredentialsProvider(credentials.NewStaticCredentialsProvider(accessKeyID, secretAccessKey, "")),
    )
    if err != nil {
        return nil, err
    }

    return s3.NewFromConfig(cfg), nil
}

func main() {
    // For testing with IAM User
	ctx := context.TODO()

    // Create a new S3 client with IAM user credentials
    s3Client, err := NewS3ClientWithIAMUserCreds(ctx)
    if err != nil {
        log.Fatalf("Failed to create S3 client: %v", err)
    }

    // Example usage of the S3 client
    // List the buckets
    result, err := s3Client.ListBuckets(ctx, nil)
    if err != nil {
        log.Fatalf("Failed to list buckets: %v", err)
    }

    for _, bucket := range result.Buckets {
        log.Printf("Bucket: %s", aws.ToString(bucket.Name))
    }
	
	// WHEN IN LAMBDA DO IT THIS WAY
	
	// // Load the AWS default configuration
	// cfg, err := config.LoadDefaultConfig(context.TODO())
    // if err != nil {
    //     log.Fatalf("unable to load SDK config, %v", err)
    // }

    // // Create a new S3 client
    // s3Client := s3.NewFromConfig(cfg)

    // // Define the bucket and key
    // bucket := "your-bucket-name"
    // key := "your-object-key"

    // // Define the content to upload
    // content := []byte("Hello, world!")

    // // Call the PutObjectOnS3 function
    // err = PutObjectOnS3(context.TODO(), s3Client, bucket, key, content)
    // if err != nil {
    //     log.Fatalf("Failed to put object on S3: %v", err)
    // }

    // fmt.Println("Object uploaded successfully.")
}