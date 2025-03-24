package utils

import (
    "bytes"
    "errors"
    "image"
    "image/jpeg"
    "image/png"
    "io"
    "mime/multipart"
    "path/filepath"
    "strconv"
    "strings"
    "time"

    "github.com/chai2010/webp"
    "github.com/nfnt/resize"
)

type CropImageConfig struct {
    Width      int
    Height     int
    CropWidth  int
    CropHeight int
    CropX      int
    CropY      int
}

func CropImage(file *multipart.FileHeader, config CropImageConfig) (io.Reader, string, int64, error) {
    /**
      Read file's data
    */
    fileData, ferr := file.Open()

    if ferr != nil {
        return nil, "", 0, errors.New("the file is corrupted")
    }

    defer fileData.Close()

    /**
      Generating file name
    */
    ext := filepath.Ext(file.Filename)
    ext = strings.ToLower(ext)
    encryptedFileName := strconv.Itoa(int(time.Now().Unix())) + GetMD5Hash(file.Filename)
    fileName := encryptedFileName + ext

    /**
      Decode image
    */
    var img image.Image
    var err error

    switch ext {
    case ".jpg", ".jpeg":
        img, err = jpeg.Decode(fileData)
    case ".png":
        img, err = png.Decode(fileData)
    case ".webp":
        img, err = webp.Decode(fileData)
    default:
        err = errors.New("unsupported image format")
    }

    if err != nil {
        return nil, "", 0, errors.New("the file is not a valid image")
    }

    /**
      Crop and resize
    */
    croppedImage := img.(interface {
        SubImage(r image.Rectangle) image.Image
    }).SubImage(image.Rect(config.CropX, config.CropY, config.CropX+config.CropWidth, config.CropY+config.CropHeight))

    /**
      Resizing to maintain aspect ratio
    */
    resizedImage := resize.Resize(uint(config.Width), uint(config.Height), croppedImage, resize.Lanczos3)
    var resizedImageBuffer bytes.Buffer

    switch ext {
    case ".jpg", ".jpeg":
        err = jpeg.Encode(&resizedImageBuffer, resizedImage, nil)
    case ".png":
        err = png.Encode(&resizedImageBuffer, resizedImage)
    case ".webp":
        err = webp.Encode(&resizedImageBuffer, resizedImage, nil)
    }

    if err != nil {
        return nil, "", 0, errors.New("the image cannot be cropped or resized")
    }

    resizedData := bytes.NewReader(resizedImageBuffer.Bytes())

    return resizedData, fileName, file.Size, nil
}
