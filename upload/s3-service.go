package upload

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

type uploadS3Config struct {
    directory string
    fileName  string
    fileSize  int64
    userId    string
}

func (store *Store) uploadToS3(fileData io.Reader, config uploadS3Config) (string, error) {
    if config.userId != "" {
        config.fileName = config.userId + "-" + config.fileName
    }

    if config.directory != "" {
        config.fileName = config.directory + "/" + config.fileName
    }

    uploader := manager.NewUploader(store.s3)
    result, err := uploader.Upload(context.Background(), &s3.PutObjectInput{
        Bucket: aws.String(store.config.S3.Bucket),
        Key:    aws.String(config.fileName),
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
            "INSERT INTO uploaded_files(id, user_id, filename, size, uploaded_url) VALUES('%s', '%s', '%s', %d, '%s')",
            id, config.userId, config.fileName, int(config.fileSize), uploadedPath,
        )

        store.db.Exec(query)
    }

    return uploadedPath, nil
}
