package migrations

import (
	"github.com/tlhunter/mig/config"
	"github.com/tlhunter/mig/database"
)

type MigrationStatus struct {
	Applied   int                  `json:"applied"`
	Unapplied int                  `json:"unapplied"`
	Skipped   int                  `json:"skipped"`        // number of skipped migrations
	Missing   int                  `json:"missing"`        // number of locally missing file migrations
	Last      *MigrationRow        `json:"last,omitempty"` // last successfully executed migration
	Next      string               `json:"next"`           // the next migration to execute
	History   []MigrationRowStatus `json:"history,omitempty"`
}

type MigrationRowStatus struct {
	Migration MigrationRow `json:"migration"`
	Status    string       `json:"status"`
}

func GetStatus(cfg config.MigConfig, dbox database.DbBox) (MigrationStatus, error) {
	var status MigrationStatus

	migFiles, err := ListFiles(cfg.Migrations)

	if err != nil {
		return status, err
	}

	migRows, err := ListRows(dbox)

	if err != nil {
		return status, err
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
			status.Last = &migRow
			mfi++
			mri++
			status.Applied++
			status.History = append(status.History, MigrationRowStatus{
				Migration: migRow,
				Status:    "applied",
			})
		} else if migFile < migRow.Name {
			// This migration is present on disk but not in database and is ready to run
			mfi++
			status.Skipped++
			status.Unapplied++
			status.History = append(status.History, MigrationRowStatus{
				Migration: MigrationRow{
					Name: migFile,
				},
				Status: "skipped",
			})
			if !didFindNext {
				status.Next = migFile
				didFindNext = true
			}
		} else if migFile > migRow.Name {
			// This migration is missing on disk which is a pretty weird scenario
			mri++
			status.Missing++
			status.Applied++
			status.History = append(status.History, MigrationRowStatus{
				Migration: migRow,
				Status:    "missing",
			})
		}
	}

	if mfi >= len(migFiles) {
		// There are still rows in the database to print
		for i := mri; i < len(migRows); i++ {
			migRow := migRows[i]
			status.Last = &migRow
			status.Applied++
			status.History = append(status.History, MigrationRowStatus{
				Migration: migRow,
				Status:    "applied",
			})
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
			status.Unapplied++
			status.History = append(status.History, MigrationRowStatus{
				Migration: MigrationRow{
					Name: migFile,
				},
				Status: "unapplied",
			})
		}
	}

	return status, nil
}
