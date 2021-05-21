package swarmtool

type Config struct {
	BackupFilePath    string
	BackupDir         string
	HotBackup         bool
	S3AccessKeyID     string
	S3SecretAccessKey string
	S3BucketName      string
	S3EndpointUrl     string
}
