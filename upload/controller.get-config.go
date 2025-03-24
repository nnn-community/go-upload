package upload

import (
    "github.com/gofiber/fiber/v2"
)

func (store *Store) getConfig(c *fiber.Ctx) error {
    configs := make(map[string]any, len(*store.uploadables))

    for i, cfg := range *store.uploadables {
        configs[i] = cfg.ToJson()
    }

    return c.Status(fiber.StatusOK).JSON(fiber.Map{
        "max_file_size": store.config.BodyLimit,
        "configs":       configs,
    })
}
