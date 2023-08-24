package main

import (
	"errors"
	"flag"
	"fmt"
	"net/http"
	"os"
	"time"

	"git.bode.fun/adsig"
	"git.bode.fun/adsig/config"
	"git.bode.fun/adsig/internal/util"
	"git.bode.fun/adsig/server"
	"github.com/charmbracelet/log"
	"github.com/go-ldap/ldap/v3"
)

var Name = "adsig" //nolint

func main() {
	var logOpt log.Options

	logOpt.Prefix = Name
	logOpt.ReportTimestamp = true
	logOpt.ReportCaller = true

	log := log.NewWithOptions(os.Stderr, logOpt)

	if err := mainE(log); err != nil {
		log.Fatal(err)
	}
}

func mainE(log *log.Logger) error {
	var accountSearch string

	flag.StringVar(&accountSearch, "account", "", "The account to search for")

	flag.Parse()

	if accountSearch == "" {
		return errors.New("email required but not provided")
	}

	accountSearch = util.NormalizeEmail(accountSearch)

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

	groups, err := adsig.GroupsFromConfig(cnf, conn)
	if err != nil {
		return err
	}

	for _, group := range groups {
		if ok, member := group.MemberBySamAccountName(accountSearch); ok {
			log.Infof("%s is a member of the group \"%s\" with %d members", accountSearch, group.Name, len(group.Members))

			for _, sig := range group.Signatures {
				tmpls, err := sig.ParseFiles()
				if err != nil {
					return err
				}

				var v struct {
					Fields map[string]string
				}

				v.Fields = make(map[string]string)

				for key, attribute := range sig.Fields {
					v.Fields[key] = member.GetAttributeValue(attribute)
				}

				for _, tmpl := range tmpls {
					fmt.Println()
					fmt.Println(tmpl.Name())
					fmt.Println("------------------------------------------------")
					tmpl.Execute(os.Stdout, v)
				}
			}
		}
	}
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
		ReadTimeout:  cnf.Server.ReadTimeout * time.Second,
		WriteTimeout: cnf.Server.WriteTimeout * time.Second,
	}

	log.Infof("Starting server on %s", addr)

	if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		return err
	}

	return nil
}
