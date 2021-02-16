// Copyright 2020 Clivern. All rights reserved.
// Use of this source code is governed by the MIT
// license that can be found in the LICENSE file.

// +build unit

package backup

import (
	"testing"

	"github.com/franela/goblin"
)

// TestPostgreSQLType
func TestPostgreSQLType(t *testing.T) {
	g := goblin.Goblin(t)

	g.Describe("#TestDumpOptions", func() {

		g.It("It should return expected options", func() {
			postgresql := &PostgreSQL{
				Host:       "127.0.0.1",
				Port:       "3306",
				Username:   "root",
				Password:   "root",
				Database:   "walrus",
				Table:      "options",
				Options:    "--no-tablespaces,--strict-names",
				OutputFile: "/tmp/result.sql",
			}

			g.Assert(postgresql.DumpOptions()).Equal("-h 127.0.0.1 -p 3306 -U root -W root -f /tmp/result.sql -d walrus -t options --no-tablespaces --strict-names")
		})

		g.It("It should return expected options", func() {
			postgresql := &PostgreSQL{
				Host:       "127.0.0.1",
				Port:       "3306",
				Username:   "root",
				Password:   "root",
				Database:   "walrus",
				Table:      "options",
				Options:    "--no-tablespaces",
				OutputFile: "/tmp/result.sql",
			}

			g.Assert(postgresql.DumpOptions()).Equal("-h 127.0.0.1 -p 3306 -U root -W root -f /tmp/result.sql -d walrus -t options --no-tablespaces")
		})

		g.It("It should return expected options", func() {
			postgresql := &PostgreSQL{
				Host:       "127.0.0.1",
				Port:       "3306",
				Username:   "root",
				Password:   "root",
				Database:   "walrus",
				Options:    "--no-tablespaces,--strict-names",
				OutputFile: "/tmp/result.sql",
			}

			g.Assert(postgresql.DumpOptions()).Equal("-h 127.0.0.1 -p 3306 -U root -W root -f /tmp/result.sql -d walrus --no-tablespaces --strict-names")
		})
	})
}
