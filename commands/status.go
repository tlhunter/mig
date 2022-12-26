package commands

import "github.com/tlhunter/mig/config"

func CommandStatus(cfg config.MigConfig) error {

	// Check if migration and lock tables exist
	//   If missing, recommend `mig init`, HALT
	//   If one is present and one missing, warn, recommend deleting, HALT
	// Check if tables adhere to correct format (useful in case mig changes?)
	//  If missing, warn about data corruption, recommend deleting tables and running `mig init`, HALT
	// Check if locked
	//   If so, mention someone might be running a long migration right now
	//   If so, recommend checking if someone ran a migration that failed if remains locked
	//   Mention that `mig unlock` can fix this if everything is OK
	// Check migrations on disk and migrations that have executed
	// Display the name of the last run migration
	// Display count of executed and unexecuted migrations
	//   If there is a skipped migration, display error, HALT
	// Display the name of the next-to-run migration, and mention `mig up` will run it

	return nil
}
