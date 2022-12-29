package commands

import (
	"github.com/fatih/color"
	"github.com/tlhunter/mig/config"
)

var Version string   // set at compile time
var BuildTime string // set at compile time

func CommandVersion(cfg config.MigConfig) error {
	color.Green("mig version: " + Version)
	color.White("build time:  " + BuildTime)

	return nil
}
