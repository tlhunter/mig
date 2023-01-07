package commands

import (
	"fmt"
	"time"

	"github.com/fatih/color"
	"github.com/tlhunter/mig/config"
	"github.com/tlhunter/mig/database"
	"github.com/tlhunter/mig/migrations"
)

func CommandList(cfg config.MigConfig) error {
	dbox, err := database.Connect(cfg.Connection)

	if err != nil {
		return err
	}

	defer dbox.Db.Close()

	status, err := migrations.GetStatus(cfg, dbox)

	if err != nil {
		return err
	}

	color.White("%5s %-48s %5s %-20s %-20s", "ID", "Migration", "Batch", "Time of Run", "Note")

	for _, entry := range status.History {
		switch entry.Status {
		case "applied":
			color.Green("%5d %-48s %5d %20s %-20s", entry.Migration.Id, entry.Migration.Name, entry.Migration.Batch, entry.Migration.Time.Format(time.RFC3339), "Applied")
		case "skipped":
			color.Red("%5s %-48s %5s %20s %-20s", "", entry.Migration.Name, "", "", "Migration Skipped!")
		case "missing":
			color.Yellow("%5d %-48s %5d %20s %-20s", entry.Migration.Id, entry.Migration.Name, entry.Migration.Batch, entry.Migration.Time.Format(time.RFC3339), "Missing File!")
		case "unapplied":
			color.Cyan("%5s %-48s %5s %20s %-20s", "", entry.Migration.Name, "", "", "Ready to Run")
		}
	}

	if status.Missing > 0 || status.Skipped > 0 {
		fmt.Println()
	}

	if status.Skipped > 0 {
		color.Red("* A skipped migration was encountered. If editing locally you may need to rename the file to the current time.")
	}

	if status.Missing > 0 {
		color.Yellow("* A missing migration was encountered. You might need to pull changes from repo.")
	}

	color.HiWhite("Applied: %d, Unapplied: %d, Skipped: %d, Missing: %d", status.Applied, status.Unapplied, status.Skipped, status.Missing)

	return nil
}
