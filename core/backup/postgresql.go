// Copyright 2020 Clivern. All rights reserved.
// Use of this source code is governed by the MIT
// license that can be found in the LICENSE file.

package backup

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
)

// PostgreSQL type
type PostgreSQL struct {
	Host       string
	Port       string
	Username   string
	Password   string
	Database   string
	Table      string
	Options    string
	OutputFile string
}

// DumpOptions dump backup options
func (m *PostgreSQL) DumpOptions() string {
	if m.Table == "" {
		return fmt.Sprintf(
			"-h %s -p %s -U %s -W %s -f %s -d %s %s",
			m.Host,
			m.Port,
			m.Username,
			m.Password,
			m.OutputFile,
			m.Database,
			strings.Replace(m.Options, ",", " ", -1),
		)
	}

	return fmt.Sprintf(
		"-h %s -p %s -U %s -W %s -f %s -d %s -t %s %s",
		m.Host,
		m.Port,
		m.Username,
		m.Password,
		m.OutputFile,
		m.Database,
		m.Table,
		strings.Replace(m.Options, ",", " ", -1),
	)
}

// BackupPostgreSQL backup and compress the postgresql dump file
func (m *Manager) BackupPostgreSQL(mysql *PostgreSQL, archive string) error {
	// Get the full executable path for the editor.
	executable, err := exec.LookPath("pg_dump")

	sqlDumpFile := strings.Replace(archive, ".tar.gz", ".sql", -1)

	mysql.OutputFile = sqlDumpFile

	if err != nil {
		return err
	}

	command := strings.Split(
		fmt.Sprintf(`%s %s`, executable, mysql.DumpOptions()),
		" ",
	)

	cmd := exec.Command(command[0], command[1:]...)

	cmd.Stdin = os.Stdin
	// cmd.Stdout = os.Stdout
	// cmd.Stderr = os.Stderr

	err = cmd.Run()

	if err != nil {
		return err
	}

	return m.BackupDirectory(sqlDumpFile, archive)
}
