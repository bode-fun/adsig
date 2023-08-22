package main

import (
	"errors"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"git.bode.fun/adsig/config"
	"git.bode.fun/adsig/server"
	"github.com/charmbracelet/log"
	"github.com/go-ldap/ldap/v3"
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
	Files  []string
}

type Group struct {
	Name      string
	Templates []Template
	Members   []*ldap.Entry
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

	conn, err := ldap.DialURL(cnf.Connection.Address)
	if err != nil {
		return err
	}

	err = conn.Bind(cnf.Connection.UserDN, cnf.Connection.Password)
	if err != nil {
		return err
	}

	groups, err := collectTemplateGroups(cnf, conn)
	if err != nil {
		return err
	}

	for _, group := range groups {
		log.Infof("Group %s: %d Members", group.Name, len(group.Members))
	}
	return nil
}

func collectTemplateGroups(cnf config.Config, conn *ldap.Conn) ([]Group, error) {
	templates, err := collectTemplates(cnf)
	if err != nil {
		return nil, err
	}

	groups := make([]Group, 0)

	for cnfGroupName, cnfGroup := range cnf.Groups {
		group := Group{
			Name:      cnfGroupName,
			Templates: make([]Template, 0),
		}

		searchRequest := &ldap.SearchRequest{
			BaseDN: cnfGroup.BaseDN,
			Scope:  ldap.ScopeWholeSubtree,
			Filter: cnfGroup.AdFilter,
		}

		searchRes, err := conn.Search(searchRequest)
		if err != nil {
			return nil, err
		}

		for _, entry := range searchRes.Entries {
			email := strings.ToLower(strings.TrimSpace(entry.GetAttributeValue("mail")))
			if email != "" {
				inExcludeList := false

				for _, excludedEmail := range cnfGroup.ExcludeEmails {
					if excludedEmail == email {
						inExcludeList = true

						break
					}
				}

				if !inExcludeList {
					group.Members = append(group.Members, entry)
				}
			}
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

	return groups, nil
}

func collectTemplates(cnf config.Config) ([]Template, error) {
	// Get templates folder
	wd, err := os.Getwd()
	if err != nil {
		return nil, err
	}

	templateDir := filepath.Join(wd, "templates")

	fileInfo, err := os.Stat(templateDir)
	if err != nil || !fileInfo.IsDir() {
		return nil, errors.New("main: templates folder is not present or can not be opened")
	}

	// Get templates
	templates := make([]Template, 0)

	for cnfTmplName, cnfTmpl := range cnf.Templates {
		tmpl := Template{
			Name:   cnfTmplName,
			Fields: cnfTmpl.Fields,
			Files:  make([]string, 0),
		}

		// Add template files to template
		tmplDir := filepath.Join(templateDir, tmpl.Name)

		fileInfo, err := os.Stat(tmplDir)
		if err != nil || !fileInfo.IsDir() {
			return nil, errors.New("main: template folder is not present or can not be opened")
		}

		files := []string{
			"signature.html",
			"signature.rtf",
			"signature.txt",
		}

		for _, file := range files {
			signaturePath := filepath.Join(tmplDir, file)
			if _, err := os.Stat(signaturePath); err != nil {
				return nil, err
			}

			tmpl.Files = append(tmpl.Files, signaturePath)
		}

		templates = append(templates, tmpl)
	}

	return templates, nil
}

func startServer(log *log.Logger, cnf config.Config) error {
	addr := fmt.Sprintf("%s:%d", cnf.Server.Host, cnf.Server.Port)
	handler := server.New()

	// TODO: find a good value for the timeouts.
	// Fixes gosec issue G114
	srv := &http.Server{
		Addr:         addr,
		Handler:      handler,
		ReadTimeout:  cnf.Server.ReadTimeout * time.Second,
		WriteTimeout: cnf.Server.WriteTimeout * time.Second,
	}

	log.Infof("Starting server on %s", addr)
	if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		return err
	}

	return nil
}
