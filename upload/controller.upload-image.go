package upload

import (
    "github.com/gofiber/fiber/v2"
    "github.com/nnn-community/go-upload/upload/codes"
    "github.com/nnn-community/go-upload/upload/utils"
    "github.com/nnn-community/go-utils/strings"
    "math"
)

func (store *Store) uploadImage(c *fiber.Ctx) error {
    form, err := c.MultipartForm()
    siwx := c.Locals("siwx").(SiwxUser)

    /**
      Validate form request and required fields
    */
    if err != nil {
        return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
            "code": codes.UploadImageInvalidMultipart,
        })
    }

    /**
      Get fields and validate
    */
    configValues, configExists := form.Value["config"]
    cropWidthValues, cropWidthExists := form.Value["crop_width"]
    cropHeightValues, cropHeightExists := form.Value["crop_height"]
    cropXValues, cropXExists := form.Value["crop_x"]
    cropYValues, cropYExists := form.Value["crop_y"]

    if !configExists || len(configValues) == 0 ||
        !cropWidthExists || len(cropWidthValues) == 0 ||
        !cropHeightExists || len(cropHeightValues) == 0 ||
        !cropXExists || len(cropXValues) == 0 ||
        !cropYExists || len(cropYValues) == 0 {
        return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
            "code": codes.UploadImageInvalidFields,
        })
    }

    file, ferr := c.FormFile("file")

    if ferr != nil {
        return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
            "code": codes.UploadImageMissingFile,
        })
    }

    config, exists := (*store.uploadables)[configValues[0]]

    if !exists || config.GetType() != "image" {
        return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
            "code": codes.UploadImageInvalidConfig,
        })
    }

    /**
      Parse crop values and validate them
    */
    cropWidth := strings.ToInt(cropWidthValues[0], 0)
    cropHeight := strings.ToInt(cropHeightValues[0], 0)
    cropX := strings.ToInt(cropXValues[0], 0)
    cropY := strings.ToInt(cropYValues[0], 0)

    if cropWidth == 0 || cropHeight == 0 {
        return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
            "code": codes.UploadImageInvalidCropSize,
        })
    }

    resizeWidth := cropWidth
    resizeHeight := cropHeight

    if config.GetWidth() != nil && config.GetHeight() != nil {
        // Check if the cropped image has same ratio as the desired resolution
        if *config.GetWidth()/cropWidth != *config.GetHeight()/cropHeight {
            return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
                "code": codes.UploadImageInvalidCropSize,
            })
        }

        resizeWidth = *config.GetWidth()
        resizeHeight = *config.GetHeight()
    } else if config.GetWidth() == nil {
        if cropHeight > *config.GetHeight() {
            resizeWidth = int(math.Round(float64(cropWidth) / (float64(cropHeight) / float64(*config.GetHeight()))))
            resizeHeight = *config.GetHeight()
        }
    } else if config.GetHeight() == nil {
        if cropWidth > *config.GetWidth() {
            resizeWidth = *config.GetWidth()
            resizeHeight = int(math.Round(float64(cropHeight) / (float64(cropWidth) / float64(*config.GetWidth()))))
        }
    }

    resizedData, fileName, fileSize, cerr := utils.CropImage(file, utils.CropImageConfig{
        Width:      resizeWidth,
        Height:     resizeHeight,
        CropWidth:  cropWidth,
        CropHeight: cropHeight,
        CropX:      cropX,
        CropY:      cropY,
    })

    if cerr != nil {
        return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
            "code": codes.UploadImageCannotBeCropped,
        })
    }

    result, s3err := store.uploadToS3(resizedData, uploadS3Config{
        directory: config.GetDirectory(),
        fileName:  fileName,
        fileSize:  fileSize,
        userId:    siwx.ID,
    })

    if s3err != nil {
        return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
            "code": codes.UploadImageCannotBeUploaded,
        })
    }

    return c.Status(fiber.StatusOK).JSON(fiber.Map{
        "path": result,
    })
}
