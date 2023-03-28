package files

import (
	"context"
	"fmt"
	"net/http"
	"netradio/pkg/cloud"
	"path/filepath"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
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

	id := uuid.New().String() + filepath.Ext(h.Filename)

	client := cloud.GetClient()
	_, err = client.PutObject(context.TODO(), &s3.PutObjectInput{
		Bucket: aws.String(cloud.GetBucket()),
		Key:    aws.String(id),
		Body:   f,
	})
	if err != nil {
		return "", err
	}

	return id, nil
}

func ToURL(id string) string {
	return fmt.Sprintf("https://storage.yandexcloud.net/%s/%s", cloud.GetBucket(), id)
}
