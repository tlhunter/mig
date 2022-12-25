package commands

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/tlhunter/mig/config"
	"github.com/tlhunter/mig/database"
)

const (
	LIST        = `SELECT id, name, batch, migration_time FROM migrations ORDER BY id ASC;`
	colorReset  = "\033[0m"
	colorRed    = "\033[31m"
	colorGreen  = "\033[32m"
	colorYellow = "\033[33m"
	colorBlue   = "\033[34m"
	colorPurple = "\033[35m"
	colorCyan   = "\033[36m"
	colorWhite  = "\033[37m"
)

type MigrationRow struct {
	Id    int
	Name  string
	Batch int
	Time  time.Time
}

func CommandList(cfg config.MigConfig) error {
	db := database.Connect(cfg.Connection)
	defer db.Close()

	directory := cfg.Migrations

	files, err := os.ReadDir(directory)

	if err != nil {
		fmt.Println("unable to read migrations direcotry")
		return err
	}

	var migFiles []string

	for _, entry := range files {
		name := entry.Name()

		if entry.IsDir() || strings.HasPrefix(name, ".") || !strings.HasSuffix(name, ".sql") {
			continue
		}

		migFiles = append(migFiles, entry.Name())
	}

	rows, err := db.Query(LIST)

	if err != nil {
		return err
	}

	defer rows.Close()

	var migRows []MigrationRow

	for rows.Next() {
		var id int
		var name string
		var batch int
		var time time.Time
		err = rows.Scan(&id, &name, &batch, &time)

		if err != nil {
			return err
		}

		migRows = append(migRows, MigrationRow{
			Id:    id,
			Name:  name,
			Batch: batch,
			Time:  time,
		})
	}

	fmt.Printf("%5s %-48s %5s %-20s\n", "ID", "Migration Name", "Batch", "Time")

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
			fmt.Print(colorGreen)
			fmt.Printf("%5d %-48s %5d %20s\n", migRow.Id, migRow.Name, migRow.Batch, migRow.Time.Format(time.RFC3339))
			fmt.Print(colorReset)
			mfi++
			mri++
		} else if migFile < migRow.Name {
			fmt.Print(colorYellow)
			fmt.Printf("%5s %-48s\n", "", migFile)
			fmt.Print(colorReset)
			mfi++
		} else if migFile > migRow.Name {
			fmt.Print(colorRed)
			fmt.Printf("%5d %-48s %5d %20s\n", migRow.Id, migRow.Name, migRow.Batch, migRow.Time.Format(time.RFC3339))
			fmt.Print(colorReset)
			mri++
		}
	}

	if mfi >= len(migFiles) {
		// TODO
	}
	if mri >= len(migRows) {
		for i := mfi; i < len(migFiles); i++ {
			migFile := migFiles[i]
			fmt.Print(colorYellow)
			fmt.Printf("%5s %-48s\n", "N/A", migFile)
			fmt.Print(colorReset)
		}
	}

	return nil
}
