package main

import (
	"os"

	"git.bode.fun/adsig/config"
	"github.com/charmbracelet/log"
)

var Name = "adsig" //nolint

func main() {
	log := log.NewWithOptions(os.Stderr, log.Options{ //nolint
		Prefix:          Name,
		ReportTimestamp: true,
		ReportCaller:    true,
	})

	if err := mainE(log); err != nil {
		log.Fatal(err)
	}
}

func mainE(log *log.Logger) error {
	cnfFile, err := os.Open("adsig.yml")
	if err != nil {
		return err
	}
	defer cnfFile.Close()

	cnf, err := config.FromYAML(cnfFile)
	if err != nil {
		return err
	}

	log.Infof("%#v", cnf)

	return nil
}
