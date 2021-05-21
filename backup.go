package main

import (
	"fmt"
	"github.com/mholt/archiver/v3"
	"log"
	"time"
)

type Backup struct {
	config *Config
}

func (b *Backup) compress() (string, error) {
	now := time.Now()
	filename := fmt.Sprintf("/tmp/%s-%s.tar.gz", b.config.BackupFilename, now.Format(time.RFC3339))
	log.Printf("creating backup %s from %s ...", filename, b.config.BackupDir)
	err := archiver.Archive([]string{b.config.BackupDir}, filename)
	if err != nil {
		log.Printf("%s compression failed with error %s", b.config.BackupDir, err)
		return "", err
	}
	log.Print("backup created")
	return filename, nil
}
