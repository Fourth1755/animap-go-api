package aws

import (
	"github.com/Fourth1755/animap-go-api/internal/core/config"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
)

type AwsAdapter struct {
	Session *session.Session
}

func NewAwsAdapter(cfg *config.AWS) (*AwsAdapter, error) {
	creds := credentials.NewStaticCredentials(cfg.AccessKey, cfg.SecretKey, "")
	sess, err := session.NewSession(&aws.Config{
		Region:      aws.String(cfg.Region),
		Credentials: creds,
	})
	if err != nil {
		return nil, err
	}

	return &AwsAdapter{Session: sess}, nil
}
