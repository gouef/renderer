package renderer

type File struct {
	Path     string
	Layout   string
	Includes []File
}
