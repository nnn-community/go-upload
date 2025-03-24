package upload

import (
    "database/sql"
    "github.com/aws/aws-sdk-go-v2/aws"
    "github.com/aws/aws-sdk-go-v2/credentials"
    "github.com/aws/aws-sdk-go-v2/service/s3"
    "github.com/gofiber/fiber/v2"
    "github.com/gofiber/storage/redis/v3"
    "github.com/nnn-community/go-upload/upload/uploadable"
)

type Store struct {
    config      Config
    app         *fiber.App
    db          *sql.DB
    s3          *s3.Client
    uploadables *map[string]uploadable.Uploadable
    redis       *redis.Storage
}

func getDbClient(connStr string) (*sql.DB, error) {
    if connStr == "" {
        return nil, nil
    }

    conn, err := sql.Open("postgres", connStr)

    if err != nil {
        return nil, err
    }

    err = conn.Ping()

    if err != nil {
        return nil, err
    }

    defer conn.Close()

    return conn, nil
}

func getS3Client(setup S3) *s3.Client {
    cfg := aws.Config{
        Region:       setup.Region,
        Credentials:  credentials.NewStaticCredentialsProvider(setup.AccessKey, setup.SecretKey, ""),
        BaseEndpoint: aws.String(setup.Endpoint),
    }

    return s3.NewFromConfig(cfg)
}
