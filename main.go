package main

import (
	"errors"
	"fmt"
	"net/http"
	"os"
	"time"

	"git.bode.fun/adsig/config"
	"git.bode.fun/adsig/server"
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

	addr := fmt.Sprintf("%s:%d", cnf.Server.Host, cnf.Server.Port)
	handler := server.New()

	// TODO: find a good value for the timeouts.
	// Fixes gosec issue G114
	srv := &http.Server{
		Addr:         addr,
		Handler:      handler,
		ReadTimeout:  time.Duration(cnf.Server.ReadTimeout) * time.Second,
		WriteTimeout: time.Duration(cnf.Server.WriteTimeout) * time.Second,
	}

	log.Infof("Starting server on %s", addr)
	if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		return err
	}

	return nil
}
