// Copyright 2020 Clivern. All rights reserved.
// Use of this source code is governed by the MIT
// license that can be found in the LICENSE file.

package backup

import (
	"fmt"
	"time"

	"github.com/clivern/walrus/core/model"
	"github.com/clivern/walrus/core/module"
	"github.com/clivern/walrus/core/storage"
	"github.com/clivern/walrus/core/util"

	"github.com/spf13/viper"
)

// Manager type
type Manager struct {
	S3 *storage.S3
}

// NewManager create a new backup manager
func NewManager(S3 *storage.S3) *Manager {
	return &Manager{
		S3: S3,
	}
}

// ProcessBackup process a backup request
func (m *Manager) ProcessBackup(message module.BackupMessage) error {
	var err error

	fileName := fmt.Sprintf("%s.tar.gz", time.Now().Format("2006-01-02_15-04-05"))

	localPath := fmt.Sprintf("%s/%s", viper.GetString("agent.backup.tmpDir"), fileName)

	remotePath := fmt.Sprintf(
		"%s/%s/%s",
		message.Cron.Hostname,
		message.Cron.ID,
		fileName,
	)

	// BackupDirectory
	if message.Cron.Request.Type == model.BackupDirectory {
		err := m.BackupDirectory(
			message.Cron.Request.Directory,
			localPath,
		)

		if err != nil {
			return err
		}

		defer util.DeleteFile(localPath)

	} else if message.Cron.Request.Type == model.BackupSQLite {
		err := m.BackupDirectory(
			message.Cron.Request.SQLitePath,
			localPath,
		)

		if err != nil {
			return err
		}

		defer util.DeleteFile(localPath)

	} else if message.Cron.Request.Type == model.BackupMySQL {
		mysql := &MySQL{
			Host:         message.Cron.Request.MySQLHost,
			Port:         message.Cron.Request.MySQLPort,
			Username:     message.Cron.Request.MySQLUsername,
			Password:     message.Cron.Request.MySQLPassword,
			AllDatabases: message.Cron.Request.MySQLAllDatabases,
			Database:     message.Cron.Request.MySQLDatabase,
			Table:        message.Cron.Request.MySQLTable,
			Options:      message.Cron.Request.MySQLOptions,
		}

		err := m.BackupMySQL(
			mysql,
			localPath,
		)

		if err != nil {
			return err
		}

		defer util.DeleteFile(localPath)
	}

	// Create bucket if not exist (ignore error)
	// TODO consider the case where bucket is missing and it fails to create
	m.S3.CreateBucket(message.Settings["backup_s3_bucket"])

	// Upload to S3
	err = m.S3.UploadFile(
		message.Settings["backup_s3_bucket"],
		localPath,
		remotePath,
		true,
	)

	if err != nil {
		return err
	}

	// Delete Old Backups (Retention Policy)
	_, err = m.S3.CleanupOld(
		message.Settings["backup_s3_bucket"],
		fmt.Sprintf("%s/%s", message.Cron.Hostname, message.Cron.ID),
		message.Cron.Request.RetentionDays,
	)

	return err
}
