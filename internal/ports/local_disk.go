package ports

import "io"

type LocalDiskAdapter interface {
	UploadFile(file io.Reader, fileName string) (string, error)
}
