package swarmtool

type Config struct {
	BackupOutputDir   string `mapstructure:"backup_output_dir"`
	Filename          string `mapstructure:"filename"`
	BackupDir         string `mapstructure:"backup_dir"`
	HotBackup         bool   `mapstructure:"hot_backup"`
	S3AccessKeyID     string `mapstructure:"s3_access_key_id"`
	S3SecretAccessKey string `mapstructure:"s3_secret_access_key"`
	S3BucketName      string `mapstructure:"s3_bucket"`
	S3EndpointUrl     string `mapstructure:"s3_endpoint_url"`
}
