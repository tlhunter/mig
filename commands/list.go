package commands

import (
	"fmt"

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

	status, err := migrations.GetStatus(cfg, db, true)

	if err != nil {
		return err
	}

	if status.Missing > 0 || status.Skipped > 0 {
		fmt.Println()
	}

	if status.Skipped > 0 {
		fmt.Print(colorRed)
		fmt.Println("* A skipped migration was encountered. If editing locally you may need to rename the file to the current time.")
		fmt.Print(colorReset)
	}

	if status.Missing > 0 {
		fmt.Print(colorYellow)
		fmt.Println("* A missing migration was encountered. You might need to pull changes from repo.")
		fmt.Print(colorReset)
	}

	fmt.Printf("Applied: %d, Unapplied: %d, Skipped: %d, Missing: %d", status.Applied, status.Unapplied, status.Skipped, status.Missing)

	return nil
}
