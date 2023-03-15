package files

import (
	"io"
	"net/http"
	"os"
	"path/filepath"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

var basePath = filepath.Join(".", "files")

func StartFileServer(router chi.Router) {
	router.Handle("/files/*", http.StripPrefix("/files/", http.FileServer(http.Dir(basePath))))
}

func Save(r *http.Request) (string, error) {
	err := r.ParseMultipartForm(32 << 20)
	if err != nil {
		return "", err
	}

	f, h, err := r.FormFile("file")
	if err != nil {
		return "", err
	}

	id := uuid.New()

	defer f.Close()
	_ = os.MkdirAll(basePath, os.ModePerm)
	fullPath := basePath + "/" + id.String() + filepath.Ext(h.Filename)
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
