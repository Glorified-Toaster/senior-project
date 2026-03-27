package storage

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
)

type LocalDiskAdapter struct {
}

func NewLocalDiskAdapter() *LocalDiskAdapter {
	return &LocalDiskAdapter{}
}

func (ld *LocalDiskAdapter) UploadFile(file io.Reader, fileName string) (string, error) {
	uploadDir := "./web/static/uploads"

	if _, err := os.Stat(uploadDir); os.IsNotExist(err) {
		if err := os.MkdirAll(uploadDir, 0755); err != nil {
			return "", err
		}
	}

	path := filepath.Join(uploadDir, fileName)
	dst, err := os.Create(path)
	if err != nil {
		return "", err
	}
	defer dst.Close()

	if _, err := io.Copy(dst, file); err != nil {
		return "", err
	}

	return fmt.Sprintf("/static/uploads/%s", fileName), nil
}
