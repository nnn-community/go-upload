package uploadable

type UploadableItem struct {
    directory  string
    uploadType string
    width      *int
    height     *int
    accept     []string
}

func (cfg UploadableItem) GetDirectory() string {
    return cfg.directory
}

func (cfg UploadableItem) GetType() string {
    return cfg.uploadType
}

func (cfg UploadableItem) GetWidth() *int {
    return cfg.width
}

func (cfg UploadableItem) GetHeight() *int {
    return cfg.height
}

func (cfg UploadableItem) GetAccept() []string {
    return cfg.accept
}

func (cfg UploadableItem) ToJson() map[string]any {
    result := map[string]any{
        "type": cfg.uploadType,
    }

    if cfg.uploadType == "file" {
        result["accept"] = cfg.accept
    }

    if cfg.uploadType == "image" {
        result["width"] = cfg.width
        result["height"] = cfg.height
    }

    return result
}
