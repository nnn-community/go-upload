package upload

import "github.com/nnn-community/go-siwx/siwx"

type Config struct {
    // S3 config to your server.
    //
    // Required.
    S3 S3 `json:"s3"`

    // Redis where your session storage is located.
    // Expects you have configured SIWX (`https://docs.nnn-community.com/react/siwx/install`).
    //
    // Required.
    Redis siwx.Redis `json:"redis"`

    // DatabaseUrl to you database, where the file log will be stored. If you do not need to log files, do not define
    // this option. Use `schema.sql` in the root to create a log table.
    // Supports only postgres database for now.
    //
    // Optional.
    DatabaseUrl string `json:"database_url"`

    // Proxy defines your custom proxy path for CDN (ie: https://cdn.example.com) when file is uploaded.
    // If not provided, will use S3.Endpoint instead.
    //
    // Optional. Default: nil
    Proxy string `json:"proxy"`

    // BodyLimit sets max upload limit.
    //
    // Optional. Default: 12 * 1024 * 1024
    BodyLimit int `json:"body_limit,omitempty"`
}

type S3 struct {
    // Optional, default: os.Getenv("S3_ENDPOINT")
    Endpoint string `json:"endpoint"`

    // Optional, default: os.Getenv("S3_REGION")
    Region string `json:"region"`

    // Optional, default: os.Getenv("S3_BUCKET")
    Bucket string `json:"bucket"`

    // Optional, default: os.Getenv("S3_ACCESS_KEY")
    AccessKey string `json:"access_key"`

    // Optional, default: os.Getenv("S3_SECRET_KEY")
    SecretKey string `json:"secret_key"`
}
