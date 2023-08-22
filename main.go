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

type Template struct {
	Name   string
	Fields map[string]string
}

type Group struct {
	Name      string
	Templates []Template
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

	templates := make([]Template, 0)

	for cnfTmplName, cnfTmpl := range cnf.Templates {
		tmpl := Template{
			Name:   cnfTmplName,
			Fields: cnfTmpl.Fields,
		}

		templates = append(templates, tmpl)
	}

	groups := make([]Group, 0)

	for cnfGroupName, cnfGroup := range cnf.Groups {
		group := Group{
			Name:      cnfGroupName,
			Templates: make([]Template, 0),
		}

		// Add templates to group
		for _, groupTmpl := range cnfGroup.Templates {
			for _, tmpl := range templates {
				if groupTmpl == tmpl.Name {
					group.Templates = append(group.Templates, tmpl)
				}
			}
		}

		groups = append(groups, group)
	}

	log.Printf("%#v", groups)

	return nil
}

func startServer(log *log.Logger, cnf config.Config) error {
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
