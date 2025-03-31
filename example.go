package goupload

import (
    "github.com/nnn-community/go-siwx/siwx"
    "github.com/nnn-community/go-upload/upload"
    "github.com/nnn-community/go-upload/upload/mime"
    "github.com/nnn-community/go-upload/upload/size"
    "github.com/nnn-community/go-upload/upload/uploadable"
    "github.com/nnn-community/go-utils/env"
    "os"
)

func main() {
    env.Load()

    app := upload.New(upload.Config{
        Proxy: os.Getenv("CDN_PROXY"),
        S3: upload.S3{
            Endpoint:  os.Getenv("S3_ENDPOINT"),
            Region:    os.Getenv("S3_REGION"),
            Bucket:    os.Getenv("S3_BUCKET"),
            AccessKey: os.Getenv("S3_ACCESS_KEY"),
            SecretKey: os.Getenv("S3_SECRET_KEY"),
        },
        Redis: siwx.Redis{
            Url: os.Getenv("REDIS_URL"),
            DB:  os.Getenv("REDIS_DB"),
        },
    })

    app.AddUpload("article", uploadable.Image("articles", size.Px(1290), size.Px(710)))
    app.AddUpload("editor", uploadable.Image("user-data", size.Px(1290), size.Auto()))
    app.AddUpload("scans", uploadable.File("documents", []string{mime.TYPE_PDF, mime.TYPE_PNG}))

    app.Listen(":" + os.Getenv("PORT"))
}
