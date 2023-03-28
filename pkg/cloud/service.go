package cloud

import (
	"context"
	"fmt"
	"log"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

var (
	client *s3.Client
	cfg    Config
)

func SetConfig(newCfg Config) {
	cfg = newCfg

	customResolver := aws.EndpointResolverWithOptionsFunc(func(service, region string, options ...interface{}) (aws.Endpoint, error) {
		if service == s3.ServiceID && region == "ru-central1" {
			return aws.Endpoint{
				PartitionID:   "yc",
				URL:           "https://storage.yandexcloud.net",
				SigningRegion: "ru-central1",
			}, nil
		}
		return aws.Endpoint{}, fmt.Errorf("unknown endpoint requested")
	})

	yandexCfg, err := config.LoadDefaultConfig(context.TODO(), config.WithEndpointResolverWithOptions(customResolver))
	if err != nil {
		log.Fatal(err)
	}

	client = s3.NewFromConfig(yandexCfg)
}

func GetClient() *s3.Client {
	return client
}

func GetBucket() string {
	return cfg.Bucket
}
