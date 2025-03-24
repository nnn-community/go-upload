package upload

import (
    "fmt"
    "github.com/gofiber/fiber/v2"
    "github.com/gofiber/fiber/v2/middleware/cors"
    "github.com/gofiber/storage/redis/v3"
    _ "github.com/lib/pq"
    "github.com/nnn-community/go-upload/upload/uploadable"
    "log"
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

    s3 := getS3Client(config.S3)

    if config.BodyLimit == 0 {
        config.BodyLimit = 12 * 1024 * 1024
    }

    app := fiber.New(fiber.Config{
        BodyLimit: config.BodyLimit,
    })

    app.Use(cors.New(cors.Config{
        AllowOriginsFunc: func(_ string) bool {
            return true
        },
        AllowMethods:     "GET,HEAD,POST,PUT,DELETE,OPTIONS,PATCH",
        AllowCredentials: true,
    }))

    redisStorage := redis.New(redis.Config{
        URL: fmt.Sprintf("%s/%d", config.Redis.Url, config.Redis.DB),
    })

    return &Store{
        config:      config,
        app:         app,
        db:          db,
        s3:          s3,
        uploadables: &(map[string]uploadable.Uploadable{}),
        redis:       redisStorage,
    }
}

func (store *Store) AddUpload(name string, cfg uploadable.Uploadable) {
    (*store.uploadables)[name] = cfg
}

func (store *Store) Listen(addr string) error {
    // Create routes before starting a service
    store.app.Get("/get-config", store.getConfig)
    store.app.Post("/upload/image", store.siwxMiddleware, store.uploadImage)
    store.app.Post("/upload/file", store.siwxMiddleware, store.uploadFile)

    return store.app.Listen(addr)
}
