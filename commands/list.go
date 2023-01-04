package commands

import (
	"fmt"

	"github.com/fatih/color"
	"github.com/tlhunter/mig/config"
	"github.com/tlhunter/mig/database"
	"github.com/tlhunter/mig/migrations"
)

func CommandList(cfg config.MigConfig) error {
	dbox := database.Connect(cfg.Connection)
	defer dbox.Db.Close()

	status, err := migrations.GetStatus(cfg, dbox, true)

	if err != nil {
		return err
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
