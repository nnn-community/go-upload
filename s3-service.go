package goupload

import (
    "context"
    "errors"
    "fmt"
    "github.com/aws/aws-sdk-go-v2/aws"
    "github.com/aws/aws-sdk-go-v2/feature/s3/manager"
    "github.com/aws/aws-sdk-go-v2/service/s3"
    "github.com/google/uuid"
    "io"
    "strings"
)

type UploadS3Config struct {
    Directory string
    FileName  string
    FileSize  int64
    UserId    string
}

func (store *Store) UploadS3(fileData io.Reader, config UploadS3Config) (string, error) {
    if config.UserId != "" {
        config.FileName = config.UserId + "-" + config.FileName
    }

    if config.Directory != "" {
        config.FileName = config.Directory + "/" + config.FileName
    }

    uploader := manager.NewUploader(store.s3)
    result, err := uploader.Upload(context.Background(), &s3.PutObjectInput{
        Bucket: aws.String(store.config.S3.Bucket),
        Key:    aws.String(config.FileName),
        Body:   fileData,
        ACL:    "public-read",
    })

    if err != nil {
        return "", errors.New("file was not uploaded to S3")
    }

    uploadedPath := result.Location
    id := uuid.New().String()

    if store.config.Proxy != "" {
        if !strings.HasSuffix(store.config.Proxy, "/") {
            store.config.Proxy = store.config.Proxy + "/"
        }

        uploadedPath = store.config.Proxy + *result.Key
    }

    if store.db != nil {
        query := fmt.Sprintf(
            "INSERT INTO uploaded_files(id, user_id, filename, size, source_url) VALUES('%s', '%s', '%s', %d, '%s')",
            id, config.UserId, config.FileName, int(config.FileSize), uploadedPath,
        )

        store.db.Exec(query)
    }

    return uploadedPath, nil
}
