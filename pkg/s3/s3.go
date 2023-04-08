package s3

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"wildfire-backend/internal/config"
)

type Client struct {
	Client *s3.S3
	Bucket string
}

func NewClient(cfg config.S3) Client {
	s, err := session.NewSession()
	if err != nil {
		panic(err)
	}
	creds := credentials.NewStaticCredentials(cfg.AccessKey, cfg.SecretKey, "")
	c := s3.New(s, aws.NewConfig().
		WithEndpoint("https://s3.storage.selcloud.ru").
		WithCredentials(creds).
		WithRegion("ru-1"))
	return Client{
		Client: c,
		Bucket: cfg.Bucket,
	}
}
