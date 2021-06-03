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
	Backup struct {
		BackupOutputDir string
		Filename        string
		BackupDir       string
		Hot             bool
		S3Client        S3Client
		S3Bucket        string
		Cluster         *Cluster
		Dockerd         *Dockerd
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

func (b *Backup) compress() (string, error) {
	timeFormat := "2006-01-02T15:04:05"
	now := time.Now().Format(timeFormat)
	tmpFile := fmt.Sprintf("%s/%s-%s.tar.gz", b.BackupOutputDir, b.Filename, now)
	log.Printf("creating backup %s from %s ...", tmpFile, b.BackupDir)
	err := archiver.Archive([]string{b.BackupDir}, tmpFile)
	if err != nil {
		log.Printf("%s compression failed with error %s", b.BackupDir, err)
		return "", err
	}
	log.Print("backup created")
	return tmpFile, nil
}

func (b *Backup) Run() error {
	if b.Hot {
		return b.hotBackup()
	}
	if b.Cluster.IsSafeToShutdown() {
		return b.coldBackup()
	}

	return nil
}

func (b *Backup) hotBackup() error {
	filePath, err := b.compress()
	if err != nil {
		return err
	}
	return b.upload(filePath)
}

func (b *Backup) coldBackup() error {
	err := b.Dockerd.Stop()
	defer func() {
		err := b.Dockerd.Start()
		if err != nil {
			log.Printf("starting docker failed with err %s", err)
		}
	}()
	if err != nil {
		return err
	}
	return b.hotBackup()
}

func (b *Backup) upload(filePath string) error {
	if b.S3Client != nil {
		err := b.S3Client.Upload(context.Background(), b.S3Bucket, filePath)
		if err != nil {
			return err
		}
	}
	return nil
}
