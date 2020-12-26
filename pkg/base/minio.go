package base

import (
	"context"
	"log"

	minio "github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"github.com/tivizi/forarun/approot/config"
)

var minioClient *minio.Client

func init() {
	config := config.GetContext().MinioConfig
	if !config.Enabled {
		return
	}
	log.Println("Minio: Enabled")
	cli, err := minio.New(config.Endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(config.AccessKey, config.SecretKey, ""),
		Secure: config.HTTPS,
	})
	if err != nil {
		panic(err)
	}
	minioClient = cli
	bucketName := "forarun-files"
	err = minioClient.MakeBucket(context.Background(), bucketName, minio.MakeBucketOptions{})
	if err != nil {
		// Check to see if we already own this bucket (which happens if you run this twice)
		exists, errBucketExists := minioClient.BucketExists(context.Background(), bucketName)
		if errBucketExists == nil && exists {
			log.Printf("We already own %s\n", bucketName)
		} else {
			panic(err)
		}
	}
}

// MinioCli 客户端
func MinioCli() *minio.Client {
	return minioClient
}
