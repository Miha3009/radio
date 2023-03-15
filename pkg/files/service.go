package files

import (
	"io"
	"net/http"
	"os"
	"path/filepath"

	"github.com/google/uuid"
)

var basePath = filepath.Join(".", "files")

func Save(r *http.Request, ext string) (string, error) {
	err := r.ParseMultipartForm(32 << 20)
	if err != nil {
		return "", err
	}

	f, _, err := r.FormFile("file")
	if err != nil {
		return "", err
	}

	id := uuid.New()

	defer f.Close()
	_ = os.MkdirAll(basePath, os.ModePerm)
	fullPath := basePath + "/" + id.String() + ext
	file, err := os.OpenFile(fullPath, os.O_WRONLY|os.O_CREATE, os.ModePerm)
	if err != nil {
		return "", err
	}
	defer file.Close()

	_, err = io.Copy(file, f)
	if err != nil {
		return "", err
	}

	return fullPath, nil
}
