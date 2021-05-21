package swarmtool

import (
	"context"
	"fmt"
	"github.com/mholt/archiver/v3"
	"github.com/minio/minio-go/v7"
	"log"
	"strings"
	"time"
)

type (
	backup struct {
		backupOutputPath string
		filename         string
		backupDir        string
		hot              bool
		s3Client         S3Client
		s3Bucket         string
	}
	BackupOpts struct {
		BackupOutputPath string
		Filename         string
		BackupDir        string
		Hot              bool
		S3Client         S3Client
		S3Bucket         string
	}
	S3Client interface {
		Upload(ctx context.Context, bucketName, filePath string) error
	}
	MinIOClient struct {
		*minio.Client
	}
)

func (c *MinIOClient) Upload(ctx context.Context, bucketName, filePath string) error {
	path := strings.Split(filePath, "/")
	objName := path[len(path)-1]
	contentType := "application/gzip"
	log.Print("uploading backup to s3")
	_, err := c.FPutObject(ctx, bucketName, objName, filePath, minio.PutObjectOptions{ContentType: contentType})
	if err == nil {
		log.Print("backup uploaded")
	}
	return err
}

func NewBackup(bo *BackupOpts) *backup {
	return &backup{
		backupOutputPath: bo.BackupOutputPath,
		filename:         bo.Filename,
		backupDir:        bo.BackupDir,
		hot:              bo.Hot,
		s3Client:         bo.S3Client,
		s3Bucket:         bo.S3Bucket,
	}
}

func (b *backup) compress() (string, error) {
	timeFormat := "2006-01-02T15:04:05"
	now := time.Now().Format(timeFormat)
	tmpFile := fmt.Sprintf("%s/%s-%s.tar.gz", b.backupOutputPath, b.filename, now)
	log.Printf("creating backup %s from %s ...", tmpFile, b.backupDir)
	err := archiver.Archive([]string{b.backupDir}, tmpFile)
	if err != nil {
		log.Printf("%s compression failed with error %s", b.backupDir, err)
		return "", err
	}
	log.Print("backup created")
	return tmpFile, nil
}

func (b *backup) Run() error {
	if b.hot {
		return b.hotBackup()
	}
	return nil
}

func (b *backup) hotBackup() error {
	filePath, err := b.compress()
	if err != nil {
		return err
	}
	return b.upload(filePath)
}

func (b *backup) upload(filePath string) error {
	if b.s3Client != nil {
		err := b.s3Client.Upload(context.Background(), b.s3Bucket, filePath)
		if err != nil {
			return err
		}
	}
	return nil
}
