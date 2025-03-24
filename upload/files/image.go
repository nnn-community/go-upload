package files

import (
    "github.com/nnn-community/go-upload/upload/files/defs"
)

type ImageUpload struct {
    Directory string
    Type      string
    Width     *int
    Height    *int
}

func Image(directory string, width *int, height *int) defs.UploadConfig {
    return ConfigValues{
        Type:      "image",
        Directory: directory,
        Width:     width,
        Height:    height,
    }
}
