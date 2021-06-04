package swarmtool

import (
	"context"
	"github.com/minio/minio-go/v7"
	"log"
	"strings"
)

type (
	// S3Client amazon s3 compatible methods.
	S3Client interface {
		Upload(ctx context.Context, bucketName, filePath string) error
	}
	// MinIOClient implements S3Client.
	MinIOClient struct {
		Client *minio.Client
	}
)

// Upload creates an object in a bucket, with contents from file at filePath.
func (c *MinIOClient) Upload(ctx context.Context, bucketName, filePath string) error {
	path := strings.Split(filePath, "/")
	objName := path[len(path)-1]
	contentType := "application/gzip"
	log.Print("uploading backup to s3")
	_, err := c.Client.FPutObject(ctx, bucketName, objName, filePath, minio.PutObjectOptions{ContentType: contentType})
	if err == nil {
		log.Print("backup uploaded")
	}
	return err
}
