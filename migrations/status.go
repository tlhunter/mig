package migrations

import (
	"database/sql"
	"time"

	"github.com/fatih/color"
	"github.com/tlhunter/mig/config"
)

type MigrationStatus struct {
	Skipped   int          // number of skipped migrations
	Missing   int          // number of locally missing file migrations
	Last      MigrationRow // last successfully executed migration
	Next      string       // the next migration to execute
	Applied   int
	Unapplied int
}

// TODO: This shouldn't print anything at all but should instead return an array of migration data
// Printing and color constants should be in the CommandList function

func GetStatus(cfg config.MigConfig, db *sql.DB, print bool) (MigrationStatus, error) {
	var status MigrationStatus

	migFiles, err := ListFiles(cfg.Migrations)

	if err != nil {
		return status, err
	}

	migRows, err := ListRows(db)

	if err != nil {
		return status, err
	}

	if print {
		color.White("%5s %-48s %5s %-20s %-20s\n", "ID", "Migration", "Batch", "Time of Run", "Note")
	}

	mfi := 0
	mri := 0
	didFindNext := false

	for {
		if mfi >= len(migFiles) {
			break
		}

		if mri >= len(migRows) {
			break
		}

		migFile := migFiles[mfi]
		migRow := migRows[mri]

		if migFile == migRow.Name {
			// This migration is present both on disk and in the database
			status.Last = migRow
			if print {
				color.Green("%5d %-48s %5d %20s %-20s\n", migRow.Id, migRow.Name, migRow.Batch, migRow.Time.Format(time.RFC3339), "Applied")
			}
			mfi++
			mri++
			status.Applied++
		} else if migFile < migRow.Name {
			// This migration is present on disk but not in database and is ready to run
			if print {
				color.Red("%5s %-48s %5s %20s %-20s\n", "", migFile, "", "", "Migration Skipped!")
			}
			mfi++
			status.Skipped++
			status.Unapplied++
			if !didFindNext {
				status.Next = migFile
				didFindNext = true
			}
		} else if migFile > migRow.Name {
			// This migration is missing on disk which is a pretty weird scenario
			if print {
				color.Yellow("%5d %-48s %5d %20s %-20s\n", migRow.Id, migRow.Name, migRow.Batch, migRow.Time.Format(time.RFC3339), "Missing File!")
			}
			mri++
			status.Missing++
			status.Applied++
		}
	}

	if mfi >= len(migFiles) {
		// There are still rows in the database to print
		for i := mri; i < len(migRows); i++ {
			migRow := migRows[i]
			status.Last = migRow
			if print {
				color.Green("%5d %-48s %5d %20s %-20s\n", migRow.Id, migRow.Name, migRow.Batch, migRow.Time.Format(time.RFC3339), "Applied")
			}
			status.Applied++
		}
	}

	if mri >= len(migRows) {
		// There are still files on disk to print
		for i := mfi; i < len(migFiles); i++ {
			migFile := migFiles[i]
			if !didFindNext {
				status.Next = migFile
				didFindNext = true
			}
			if print {
				color.Cyan("%5s %-48s %5s %20s %-20s\n", "", migFile, "", "", "Ready to Run")
			}
			status.Unapplied++
		}
	}

	return status, nil
}
