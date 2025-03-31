package upload

import (
    "github.com/gofiber/fiber/v2"
    "github.com/nnn-community/go-siwx/siwx"
    "github.com/nnn-community/go-upload/upload/uploadable"
    "log"
    "os"
)

type App interface {
    // AddUpload appends new upload config for the app, use uploadable.File or uploadable.Image to set uploadables.
    AddUpload(name string, cfg uploadable.Uploadable)

    // Listen serves HTTP requests from the given addr.
    Listen(addr string) error
}

func New(config Config) App {
    db, err := getDbClient(config.DatabaseUrl)

    if err != nil {
        log.Fatal("Error connecting to the database:", err)
    }

    if config.S3.Endpoint == "" {
        config.S3.Endpoint = os.Getenv("S3_ENDPOINT")
    }

    if config.S3.Region == "" {
        config.S3.Region = os.Getenv("S3_REGION")
    }

    if config.S3.Bucket == "" {
        config.S3.Bucket = os.Getenv("S3_BUCKET")
    }

    if config.S3.AccessKey == "" {
        config.S3.AccessKey = os.Getenv("S3_ACCESS_KEY")
    }

    if config.S3.SecretKey == "" {
        config.S3.SecretKey = os.Getenv("S3_SECRET_KEY")
    }

    s3 := getS3Client(config.S3)

    if config.BodyLimit == 0 {
        config.BodyLimit = 12 * 1024 * 1024
    }

    app := siwx.New(siwx.Config{
        Fiber: fiber.Config{
            BodyLimit: config.BodyLimit,
        },
        Redis: config.Redis,
    })

    return &Store{
        config:      config,
        app:         app,
        db:          db,
        s3:          s3,
        uploadables: &(map[string]uploadable.Uploadable{}),
    }
}

func (store *Store) AddUpload(name string, cfg uploadable.Uploadable) {
    (*store.uploadables)[name] = cfg
}

func (store *Store) Listen(addr string) error {
    // Create routes before starting a service
    store.app.Get("/get-config", store.getConfig)
    store.app.Post("/upload/image", siwx.Middleware, store.uploadImage)
    store.app.Post("/upload/file", siwx.Middleware, store.uploadFile)

    return store.app.Listen(addr)
}
