package files

import (
    "github.com/nnn-community/go-upload/upload/files/defs"
)

type FileUpload struct {
    Directory string
    Type      string
    Accept    []string
}

func File(directory string, accept []string) defs.UploadConfig {
    return ConfigValues{
        Type:      "file",
        Directory: directory,
        Accept:    accept,
    }
}
