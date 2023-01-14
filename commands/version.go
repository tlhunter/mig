package commands

import (
	"github.com/fatih/color"
)

var Version string   // set at compile time
var BuildTime string // set at compile time

func CommandVersion() error {
	color.Green("mig version: " + Version)
	color.White("build time:  " + BuildTime)

	return nil
}
