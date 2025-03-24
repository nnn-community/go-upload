package upload

import "github.com/nnn-community/go-upload/upload/defs"

type FileUpload struct {
    Directory string
    Type      string
    Accept    []string
}

func NewFileUpload(directory string, accept []string) defs.UploadConfig {
    return ConfigValues{
        Type:      "file",
        Directory: directory,
        Accept:    accept,
    }
}
