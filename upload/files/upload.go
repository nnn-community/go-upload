package files

type ConfigValues struct {
    Directory string
    Type      string
    Width     *int
    Height    *int
    Accept    []string
}

func (cfg ConfigValues) GetDirectory() string {
    return cfg.Directory
}

func (cfg ConfigValues) GetType() string {
    return cfg.Type
}

func (cfg ConfigValues) GetWidth() *int {
    return cfg.Width
}

func (cfg ConfigValues) GetHeight() *int {
    return cfg.Height
}

func (cfg ConfigValues) GetAccept() []string {
    return cfg.Accept
}

func (cfg ConfigValues) ToJson() map[string]any {
    result := map[string]any{
        "type": cfg.Type,
    }

    if cfg.Type == "file" {
        result["accept"] = cfg.Accept
    }

    if cfg.Type == "image" {
        result["width"] = cfg.Width
        result["height"] = cfg.Height
    }

    return result
}
