package upload

import (
    "encoding/json"
    "errors"
    "fmt"
    "github.com/gofiber/fiber/v2"
    "github.com/gofiber/storage/redis/v3"
    "net/url"
    "strings"
)

func (store *Store) siwxMiddleware(c *fiber.Ctx) error {
    sid := c.Cookies("connect.sid")
    siwx, err := extractSession(store.redis, sid)

    if err != nil {
        return c.SendStatus(fiber.StatusUnauthorized)
    }

    c.Locals("siwx", siwx)

    return c.Next()
}

type sessionData struct {
    SiwxUser *SiwxUser `json:"siwxUser"`
}

func extractSession(r *redis.Storage, sid string) (SiwxUser, error) {
    decoded, err := url.QueryUnescape(sid)

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

    var data sessionData

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
