package upload

import (
    "database/sql"
    "encoding/json"
    "errors"
    "fmt"
    "github.com/aws/aws-sdk-go-v2/aws"
    "github.com/aws/aws-sdk-go-v2/credentials"
    "github.com/aws/aws-sdk-go-v2/service/s3"
    "github.com/gofiber/fiber/v2"
    "github.com/gofiber/fiber/v2/middleware/cors"
    "github.com/gofiber/fiber/v2/middleware/session"
    "github.com/gofiber/storage/redis/v3"
    _ "github.com/lib/pq"
    "github.com/nnn-community/go-upload/upload/files/defs"
    "log"
    "net/url"
    "strings"
)

// New creates a new Upload named instance.
//
//	app := upload.New(upload.Config{
//	    DatabaseUrl: "host= user= password= dbname= port=",
//      Proxy: "https://cdn.example.com",
//      BodyLimit: 500,
//      S3: upload.S3{
//          Endpoint:  "https://s3.amazonaws.com",
//          Region:    "us-east-1",
//          Bucket:    "goupload",
//          AccessKey: "123",
//          SecretKey: "456",
//      },
//      Redis: upload.Redis{
//          Url: "localhost:6379",
//          DB:  0,
//      },
//	})
func New(config Config) Store {
    /**
      Database
    */
    db, err := getDbClient(config.DatabaseUrl)

    if err != nil {
        log.Fatal("Error connecting to the database:", err)
    }

    /**
      Database
    */
    s3 := getS3Client(config.S3)

    /**
      Create fiber server
    */
    if config.BodyLimit == 0 {
        config.BodyLimit = 12 * 1024 * 1024
    }

    app := fiber.New(fiber.Config{
        BodyLimit: config.BodyLimit,
    })

    /**
      CORS
    */
    app.Use(cors.New(cors.Config{
        AllowOriginsFunc: func(_ string) bool {
            return true
        },
        AllowMethods:     "GET,HEAD,POST,PUT,DELETE,OPTIONS,PATCH",
        AllowCredentials: true,
    }))

    /**
      Redis
    */
    redisStorage := redis.New(redis.Config{
        URL: fmt.Sprintf("%s/%d", config.Redis.Url, config.Redis.DB),
    })

    /**
      Session
    */
    sessionStorage := session.New(session.Config{
        Storage:        redisStorage,
        CookieSecure:   false,
        CookieSameSite: "Lax",
        CookieHTTPOnly: true,
        KeyLookup:      "cookie:connect.sid",
    })

    return Store{
        config:         config,
        app:            app,
        db:             db,
        s3:             s3,
        configs:        &(map[string]defs.UploadConfig{}),
        redisStorage:   redisStorage,
        sessionStorage: sessionStorage,
    }
}

// AddUpload append new config for upload type
//
//	app.AddUpload("picture", upload.NewImageUpload("pictures", upload.Px(500), upload.Px(350)))
//	app.AddUpload("editor", upload.NewImageUpload("customs", upload.Px(500), upload.Auto()))
//	app.AddUpload("document", upload.NewFileUpload("documents", []string{upload.TYPE_PDF}))
func (store *Store) AddUpload(name string, cfg defs.UploadConfig) {
    (*store.configs)[name] = cfg
}

// Listen serves HTTP requests from the given addr.
//
//	app.Listen(":3200")
//	app.Listen("127.0.0.1:3200")
func (store *Store) Listen(addr string) error {
    // Create routes before starting a service
    store.app.Get("/get-config", store.getConfig)
    store.app.Post("/upload/image", store.siwxMiddleware, store.uploadImage)
    store.app.Post("/upload/file", store.siwxMiddleware, store.uploadFile)

    return store.app.Listen(addr)
}

// - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - -
//
// Local functions and structs
//
// - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - -

type Store struct {
    config         Config
    app            *fiber.App
    db             *sql.DB
    s3             *s3.Client
    configs        *map[string]defs.UploadConfig
    redisStorage   *redis.Storage
    sessionStorage *session.Store
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

func (store *Store) siwxMiddleware(c *fiber.Ctx) error {
    sess, err := store.sessionStorage.Get(c)

    if err != nil {
        return err
    }

    siwx, err := extractSession(store.redisStorage, sess)

    if err != nil {
        return c.SendStatus(fiber.StatusUnauthorized)
    }

    c.Locals("siwx", siwx)

    return c.Next()
}

// - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - -
//
// Get SIWX session
//
// - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - -

type SessionData struct {
    SiwxUser *SiwxUser `json:"siwxUser"`
}

func extractSession(r *redis.Storage, s *session.Session) (SiwxUser, error) {
    cookieValue := s.ID()
    decoded, err := url.QueryUnescape(cookieValue)

    if err != nil {
        return SiwxUser{}, err
    }

    parts := strings.SplitN(decoded, ".", 2)

    if len(parts) != 2 {
        return SiwxUser{}, errors.New("invalid cookie format")
    }

    exp := strings.SplitN(parts[0], ":", 2)

    if len(parts) != 2 {
        return SiwxUser{}, errors.New("invalid cookie format")
    }

    sessionId := "sess:" + exp[1]
    siwx, err := r.Get(sessionId)

    if err != nil {
        return SiwxUser{}, err
    }

    if string(siwx) == "" {
        return SiwxUser{}, fmt.Errorf("empty JSON string")
    }

    var data SessionData

    if err := json.Unmarshal(siwx, &data); err != nil {
        return SiwxUser{}, fmt.Errorf("invalid JSON: %w", err)
    }

    if data.SiwxUser == nil {
        return SiwxUser{}, fmt.Errorf("siwxUser field missing")
    }

    if data.SiwxUser.UserData == nil {
        data.SiwxUser.UserData = make(map[string]interface{})
    }

    return *data.SiwxUser, nil
}
