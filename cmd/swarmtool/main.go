package main

import (
	"github.com/saeed617/swarmtool"
	"log"
)

func main() {
	c := swarmtool.Config{
		BackupFilePath: "./harmony",
		BackupDir:      "/home/saeed/dev/harmony",
		HotBackup:      true,
	}
	b := swarmtool.NewBackup(&swarmtool.BackupOpts{
		FilePath:  c.BackupFilePath,
		BackupDir: c.BackupDir,
		Hot:       c.HotBackup,
	})
	err := b.Backup()
	if err != nil {
		log.Fatal(err)
	}
}
