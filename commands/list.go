package commands

import (
	"fmt"
	"time"

	"github.com/tlhunter/mig/config"
	"github.com/tlhunter/mig/database"
	"github.com/tlhunter/mig/migrations"
)

const (
	colorReset  = "\033[0m"
	colorRed    = "\033[31m"
	colorGreen  = "\033[32m"
	colorYellow = "\033[33m"
	colorBlue   = "\033[34m"
	colorPurple = "\033[35m"
	colorCyan   = "\033[36m"
	colorWhite  = "\033[37m"
)

func CommandList(cfg config.MigConfig) error {
	db := database.Connect(cfg.Connection)
	defer db.Close()

	migFiles, err := migrations.ListFiles(cfg.Migrations)

	if err != nil {
		return err
	}

	migRows, err := migrations.ListRows(db)

	if err != nil {
		return err
	}

	fmt.Printf("%5s %-48s %5s %-20s %-20s\n", "ID", "Migration", "Batch", "Time of Run", "Note")

	encounterSkipped := false
	encounterMissing := false
	mfi := 0
	mri := 0

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
			fmt.Print(colorGreen)
			fmt.Printf("%5d %-48s %5d %20s %-20s\n", migRow.Id, migRow.Name, migRow.Batch, migRow.Time.Format(time.RFC3339), "Applied")
			fmt.Print(colorReset)
			mfi++
			mri++
		} else if migFile < migRow.Name {
			// This migration is present on disk but not in database and is ready to run
			fmt.Print(colorRed)
			fmt.Printf("%5s %-48s %5s %20s %-20s\n", "", migFile, "", "", "Migration Skipped!")
			fmt.Print(colorReset)
			mfi++
			encounterSkipped = true
		} else if migFile > migRow.Name {
			// This migration is missing on disk which is a pretty weird scenario
			fmt.Print(colorYellow)
			fmt.Printf("%5d %-48s %5d %20s %-20s\n", migRow.Id, migRow.Name, migRow.Batch, migRow.Time.Format(time.RFC3339), "Missing File!")
			fmt.Print(colorReset)
			mri++
			encounterMissing = true
		}
	}

	if mfi >= len(migFiles) {
		// There are still rows in the database to print
		for i := mri; i < len(migRows); i++ {
			migRow := migRows[i]
			fmt.Print(colorGreen)
			fmt.Printf("%5d %-48s %5d %20s %-20s\n", migRow.Id, migRow.Name, migRow.Batch, migRow.Time.Format(time.RFC3339), "Applied")
			fmt.Print(colorReset)
		}
	}

	if mri >= len(migRows) {
		// There are still files on disk to print
		for i := mfi; i < len(migFiles); i++ {
			migFile := migFiles[i]
			fmt.Print(colorCyan)
			fmt.Printf("%5s %-48s %5s %20s %-20s\n", "", migFile, "", "", "Ready to Run")
			fmt.Print(colorReset)
		}
	}

	if encounterMissing || encounterSkipped {
		fmt.Println()
	}

	if encounterSkipped {
		fmt.Print(colorRed)
		fmt.Println("* A skipped migration was encountered. If editing locally you may need to rename the file to the current time.")
		fmt.Print(colorReset)
	}

	if encounterMissing {
		fmt.Print(colorYellow)
		fmt.Println("* A missing migration was encountered. You might need to pull changes from repo.")
		fmt.Print(colorReset)
	}

	return nil
}
