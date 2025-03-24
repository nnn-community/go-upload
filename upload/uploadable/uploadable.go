package uploadable

type Uploadable interface {
    GetDirectory() string
    GetType() string
    ToJson() map[string]any
    GetWidth() *int
    GetHeight() *int
    GetAccept() []string
}

func File(directory string, accept []string) Uploadable {
    return UploadableItem{
        uploadType: "file",
        directory:  directory,
        accept:     accept,
    }
}

func Image(directory string, width *int, height *int) Uploadable {
    return UploadableItem{
        uploadType: "image",
        directory:  directory,
        width:      width,
        height:     height,
    }
}
