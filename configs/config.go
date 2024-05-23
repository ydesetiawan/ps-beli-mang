package configs

import (
	"fmt"
	"ps-beli-mang/pkg/helper"
	"strings"

	"github.com/joho/godotenv"
	"github.com/kelseyhightower/envconfig"
)

type MainConfig struct {
	Database struct {
		Host         string `envconfig:"DB_HOST"`
		Port         string `envconfig:"DB_PORT"`
		Username     string `envconfig:"DB_USERNAME"`
		Password     string `envconfig:"DB_PASSWORD"`
		DatabaseName string `envconfig:"DB_NAME"`
		Params       string `envconfig:"DB_PARAMS"`
	}
	S3 struct {
		ID         string `envconfig:"AWS_ACCESS_KEY_ID"`
		SecretKey  string `envconfig:"AWS_SECRET_ACCESS_KEY"`
		BucketName string `envconfig:"AWS_S3_BUCKET_NAME"`
		Region     string `envconfig:"AWS_REGION"`
	}
}

func (mc *MainConfig) GetDsnString() string {
	dc := mc.Database
	return fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s?%s",
		dc.Username,
		dc.Password,
		dc.Host,
		dc.Port,
		dc.DatabaseName,
		strings.Replace(dc.Params, "\"", "", -1),
	)
}

func Init() *MainConfig {
	godotenv.Load(".env")

	var cfg MainConfig
	err := envconfig.Process("", &cfg)
	if err != nil {
		helper.PanicIfError(err, "Error when reading env config.")
	}
	return &cfg
}
