package upload

import (
    "github.com/gofiber/fiber/v2"
    "github.com/nnn-community/go-upload/upload/codes"
    "github.com/nnn-community/go-upload/upload/utils"
    "mime/multipart"
    "net/http"
    "path/filepath"
    "strings"
)

func (store *Store) uploadFile(c *fiber.Ctx) error {
    form, err := c.MultipartForm()
    siwx := c.Locals("siwx").(SiwxUser)

    /**
      Validate form request and required fields
    */
    if err != nil {
        return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
            "code": codes.UploadFileInvalidMultipart,
        })
    }

    /**
      Get fields and validate
    */
    configValues, configExists := form.Value["config"]

    if !configExists || len(configValues) == 0 {
        return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
            "code": codes.UploadFileInvalidFields,
        })
    }

    file, ferr := c.FormFile("file")

    if ferr != nil {
        return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
            "code": codes.UploadFileMissingFile,
        })
    }

    config, exists := (*store.uploadables)[configValues[0]]

    if !exists || config.GetType() != "file" {
        return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
            "code": codes.UploadFileInvalidConfig,
        })
    }

    fileData, ferr := file.Open()

    if ferr != nil {
        return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
            "code": codes.UploadFileCorrupted,
        })
    }

    ext := filepath.Ext(file.Filename)
    ext = strings.ToLower(ext)
    fileName := utils.GetMD5Hash(file.Filename) + ext

    if len(config.GetAccept()) > 0 {
        mime, err := getMimeType(fileData)

        if err != nil || !contains(config.GetAccept(), mime) {
            return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
                "code": codes.UploadFileInvalidFields,
            })
        }
    }

    result, s3err := store.uploadToS3(fileData, uploadS3Config{
        directory: config.GetDirectory(),
        fileName:  fileName,
        fileSize:  file.Size,
        userId:    siwx.ID,
    })

    if s3err != nil {
        return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
            "code": codes.UploadFileCannotBeUploaded,
        })
    }

    return c.Status(fiber.StatusOK).JSON(fiber.Map{
        "path": result,
    })
}

func contains(arr []string, str string) bool {
    for _, v := range arr {
        if v == str {
            return true
        }
    }

    return false
}

func getMimeType(file multipart.File) (string, error) {
    buf := make([]byte, 512)
    _, err := file.Read(buf)

    if err != nil {
        return "", err
    }

    _, err = file.Seek(0, 0)

    if err != nil {
        return "", err
    }

    return http.DetectContentType(buf), nil
}
