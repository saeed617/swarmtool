package main

import (
	"fmt"
	"github.com/mholt/archiver/v3"
	"log"
	"time"
)

type (
	backup struct {
		filePath  string
		backupDir string
		hot       bool
	}
	BackupOpts struct {
		FilePath  string
		BackupDir string
		Hot       bool
	}
)

func NewBackup(bo *BackupOpts) *backup {
	return &backup{
		filePath:  bo.FilePath,
		backupDir: bo.BackupDir,
		hot:       bo.Hot,
	}
}

func (b *backup) compress() (string, error) {
	now := time.Now()
	tmpFile := fmt.Sprintf("%s-%s.tar.gz", b.filePath, now.Format(time.RFC3339))
	log.Printf("creating backup %s from %s ...", tmpFile, b.backupDir)
	err := archiver.Archive([]string{b.backupDir}, tmpFile)
	if err != nil {
		log.Printf("%s compression failed with error %s", b.backupDir, err)
		return "", err
	}
	log.Print("backup created")
	return tmpFile, nil
}

func (b *backup) Backup() error {
	if b.hot {
		return b.hotBackup()
	}
	return nil
}

func (b *backup) hotBackup() error {
	_, err := b.compress()
	if err != nil {
		return err
	}
	return nil
}
