package cmd

import (
	"context"
	"github.com/coreos/go-systemd/v22/dbus"
	"github.com/docker/docker/client"
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
	var dbusConn swarmtool.Connection
	var dockerd *swarmtool.Dockerd
	var cluster *swarmtool.Cluster
	var err error
	if config.S3AccessKeyID != "" {
		s3Client, err = createS3Client()
		if err != nil {
			return err
		}
	}
	if !config.HotBackup {
		cluster, err = createCluster()
		if err != nil {
			return err
		}
		dbusConn, err = createDbusConn()
		if err != nil {
			return err
		}
		defer dbusConn.Close()
		dockerd = createDockerd(dbusConn)
	}

	b := &swarmtool.Backup{
		BackupOutputDir: config.BackupOutputDir,
		Filename:        config.Filename,
		BackupDir:       config.BackupDir,
		Hot:             config.HotBackup,
		S3Client:        s3Client,
		S3Bucket:        config.S3BucketName,
		Cluster:         cluster,
		Dockerd:         dockerd,
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
	return &swarmtool.MinIOClient{Client: minioClient}, nil
}

func createCluster() (*swarmtool.Cluster, error) {
	dockerClient, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		return nil, err
	}
	return &swarmtool.Cluster{
		Client: &swarmtool.ClusterClient{
			Client: dockerClient,
		},
	}, nil
}

func createDbusConn() (swarmtool.Connection, error) {
	ctx := context.Background()
	conn, err := dbus.NewSystemdConnectionContext(ctx)
	if err != nil {
		return nil, err
	}
	dbusConn := &swarmtool.DbusConnection{Conn: conn}
	return dbusConn, nil
}

func createDockerd(conn swarmtool.Connection) *swarmtool.Dockerd {
	return &swarmtool.Dockerd{DbusConn: conn}
}
