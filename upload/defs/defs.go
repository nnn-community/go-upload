package defs

type UploadConfig interface {
    GetDirectory() string
    GetType() string
    ToJson() map[string]any
    GetWidth() *int
    GetHeight() *int
    GetAccept() []string
}
