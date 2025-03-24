package upload

type Config struct {
    // S3 config to your server.
    //
    // Required.
    S3 S3 `json:"s3"`

    // Redis where your session storage is located.
    // Expects you have configured SIWX (`https://docs.nnn-community.com/react/siwx/install`).
    //
    // Required.
    Redis Redis `json:"redis"`

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

    // DisableHttps disables HTTPS for session.
    // !!! NOT RECOMMENDED !!!
    // Set to `false` ONLY when your localhost doesn't support HTTPS connection.
    //
    // Optional. Default: false
    DisableHttps bool `json:"disable_https,omitempty"`
}

type S3 struct {
    // Required.
    Endpoint string `json:"endpoint"`

    // Required.
    Region string `json:"region"`

    // Required.
    Bucket string `json:"bucket"`

    // Required.
    AccessKey string `json:"access_key"`

    // Required.
    SecretKey string `json:"secret_key"`
}

type Redis struct {
    // Provide Url string to the Redis Server, omit DB from the string as it will be added from the DB option.
    //
    // Required.
    Url string `json:"url"`

    // DB number with session storage.
    //
    // Required.
    DB int `json:"db"`
}

// SiwxUser is available to get when the route is authenticated via `siwx := c.Locals("siwx").(SiwxUser)`
type SiwxUser struct {
    ID          string                 `json:"id"`
    Address     string                 `json:"address"`
    Permissions []string               `json:"permissions"`
    UserData    map[string]interface{} `json:"userData"`
}
