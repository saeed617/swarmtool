package cmd

import (
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"github.com/saeed617/swarmtool"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(backupCmd)
}

var backupCmd = &cobra.Command{
	Use:   "backup",
	Short: "Backup swarm cluster state to s3",
	RunE: func(cmd *cobra.Command, args []string) error {
		return backup()
	},
}

func backup() error {
	var s3Client swarmtool.S3Client
	var err error
	if config.S3AccessKeyID != "" {
		s3Client, err = createS3Client()
		if err != nil {
			return err
		}
	}

	b := &swarmtool.Backup{
		BackupOutputDir: config.BackupOutputDir,
		Filename:        config.Filename,
		BackupDir:       config.BackupDir,
		Hot:             config.HotBackup,
		S3Client:        s3Client,
		S3Bucket:        config.S3BucketName,
	}

	err = b.Run()
	if err != nil {
		return err
	}
	return nil
}

func createS3Client() (swarmtool.S3Client, error) {
	minioClient, err := minio.New(config.S3EndpointUrl, &minio.Options{
		Creds:  credentials.NewStaticV4(config.S3AccessKeyID, config.S3SecretAccessKey, ""),
		Secure: true,
	})
	if err != nil {
		return nil, err
	}
	return &swarmtool.MinIOClient{minioClient}, nil
}
