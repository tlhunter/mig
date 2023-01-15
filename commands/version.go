package commands

import (
	"github.com/fatih/color"
	"github.com/tlhunter/mig/config"
	"github.com/tlhunter/mig/result"
)

var Version string   // set at compile time
var BuildTime string // set at compile time

type CommandVersionResult struct {
	Version   string `json:"version"`
	BuildTime string `json:"build_time"`
}

func CommandVersion(cfg config.MigConfig) result.Response {
	data := CommandVersionResult{
		Version:   Version,
		BuildTime: BuildTime,
	}

	res := result.NewSerializable(color.GreenString("mig version: "+Version), data)

	res.AddSuccessLn(color.WhiteString("build time:  " + BuildTime))

	return *res
}
