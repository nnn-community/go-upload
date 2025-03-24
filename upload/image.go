package upload

import "github.com/nnn-community/go-upload/upload/defs"

type ImageUpload struct {
    Directory string
    Type      string
    Width     *int
    Height    *int
}

func NewImageUpload(directory string, width *int, height *int) defs.UploadConfig {
    return ConfigValues{
        Type:      "image",
        Directory: directory,
        Width:     width,
        Height:    height,
    }
}
