package goupload

import (
    "github.com/nnn-community/go-upload/upload"
    "github.com/nnn-community/go-upload/upload/files"
    "github.com/nnn-community/go-upload/upload/size"
    "github.com/nnn-community/go-utils/env"
    "github.com/nnn-community/go-utils/strings"
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
        Redis: upload.Redis{
            Url: os.Getenv("REDIS_URL"),
            DB:  strings.ToInt(os.Getenv("REDIS_DB"), 0),
        },
    })

    app.AddUpload("article", files.Image("articles", size.Px(1290), size.Px(710)))
    app.AddUpload("editor", files.Image("user-data", size.Px(1290), size.Auto()))
    app.AddUpload("scans", files.File("documents", []string{files.TYPE_PDF, files.TYPE_PNG}))

    app.Listen(":" + os.Getenv("PORT"))
}
