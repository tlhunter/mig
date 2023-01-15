package commands

import (
	"time"

	"github.com/fatih/color"
	"github.com/tlhunter/mig/config"
	"github.com/tlhunter/mig/database"
	"github.com/tlhunter/mig/migrations"
	"github.com/tlhunter/mig/result"
)

func CommandList(cfg config.MigConfig) result.Response {
	dbox, err := database.Connect(cfg.Connection)

	if err != nil {
		return *result.NewErrorWithDetails("database connection error", "db_conn", err)
	}

	defer dbox.Db.Close()

	status, err := migrations.GetStatus(cfg, dbox)

	if err != nil {
		return *result.NewErrorWithDetails("unable to get migration status", "unable_get_status", err)
	}

	res := result.NewSerializable(color.WhiteString("%5s %-48s %5s %-20s %-20s", "ID", "Migration", "Batch", "Time of Run", "Note"), status.History)

	for _, entry := range status.History {
		switch entry.Status {
		case "applied":
			res.AddSuccessLn(color.GreenString("%5d %-48s %5d %20s %-20s", entry.Migration.Id, entry.Migration.Name, entry.Migration.Batch, entry.Migration.Time.Format(time.RFC3339), "Applied"))
		case "skipped":
			res.AddSuccessLn(color.RedString("%5s %-48s %5s %20s %-20s", "", entry.Migration.Name, "", "", "Migration Skipped!"))
		case "missing":
			res.AddSuccessLn(color.YellowString("%5d %-48s %5d %20s %-20s", entry.Migration.Id, entry.Migration.Name, entry.Migration.Batch, entry.Migration.Time.Format(time.RFC3339), "Missing File!"))
		case "unapplied":
			res.AddSuccessLn(color.CyanString("%5s %-48s %5s %20s %-20s", "", entry.Migration.Name, "", "", "Unapplied"))
		}
	}

	if status.Missing > 0 || status.Skipped > 0 {
		res.AddSuccessLn("")
	}

	if status.Skipped > 0 {
		res.AddSuccessLn(color.RedString("* A skipped migration was encountered. If editing locally you may need to rename the file to the current time."))
	}

	if status.Missing > 0 {
		res.AddSuccessLn(color.YellowString("* A missing migration was encountered. You might need to pull changes from repo."))
	}

	res.AddSuccessLn(color.HiWhiteString("Applied: %d, Unapplied: %d, Skipped: %d, Missing: %d", status.Applied, status.Unapplied, status.Skipped, status.Missing))

	return *res
}
