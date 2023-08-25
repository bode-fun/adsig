package main

import (
	"errors"
	"flag"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"git.bode.fun/adsig"
	"git.bode.fun/adsig/config"
	"git.bode.fun/adsig/server"
	"github.com/charmbracelet/log"
	"github.com/dgraph-io/badger/v3"
	"github.com/go-ldap/ldap/v3"
	"github.com/google/uuid"
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

	cnfFile, err := os.Open("adsig.yml")
	if err != nil {
		return err
	}
	defer cnfFile.Close()

	cnf, err := config.FromYAML(cnfFile)
	if err != nil {
		return err
	}

	db, err := badger.Open(badger.DefaultOptions("").WithInMemory(true))
	if err != nil {
		return err
	}
	defer db.Close()

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

	renderedTmplsPerMember := make(map[string][]struct {
		ID   uuid.UUID
		File string
		Data string
	}, 0)

	for _, group := range groups {
		if ok, member := group.MemberBySamAccountName(accountSearch); ok {
			log.Infof("%s is a member of the group \"%s\" with %d members", accountSearch, group.Name, len(group.Members))

			renderedTmplsPerMember[accountSearch] = make([]struct {
				ID   uuid.UUID
				File string
				Data string
			}, 0)

			for _, sig := range group.Signatures {
				tmpls, err := sig.ParseFiles()
				if err != nil {
					return err
				}

				var tmplData struct {
					Fields map[string]string
				}

				tmplData.Fields = make(map[string]string)

				for key, attribute := range sig.Fields {
					tmplData.Fields[key] = member.GetAttributeValue(attribute)
				}

				for _, tmpl := range tmpls {
					var b strings.Builder

					if err := tmpl.Execute(&b, tmplData); err != nil {
						break
					}

					u := uuid.New()

					ext := filepath.Ext(tmpl.Name())

					renderedTmplsPerMember[accountSearch] = append(renderedTmplsPerMember[accountSearch], struct {
						ID   uuid.UUID
						File string
						Data string
					}{
						ID:   u,
						Data: b.String(),
						File: sig.Name + ext,
					})
				}
			}
		}
	}

	for memberAccount, renderedTmpls := range renderedTmplsPerMember {
		log.Infof("Printing tmpls for %s", memberAccount)

		for _, tmpl := range renderedTmpls {
			db.Update(func(txn *badger.Txn) error {
				e := badger.NewEntry(tmpl.ID[:], []byte(tmpl.Data)).WithTTL(time.Hour)
				err := txn.SetEntry(e)
				return err
			})

			valCopy := make([]byte, 0)

			db.View(func(txn *badger.Txn) error {
				itm, err := txn.Get(tmpl.ID[:])
				if err != nil {
					return err
				}

				cpy, err := itm.ValueCopy(nil)
				if err != nil {
					return err
				}

				valCopy = append(valCopy, cpy...)

				return nil
			})

			db.Update(func(txn *badger.Txn) error {
				err := txn.Delete(tmpl.ID[:])
				return err
			})

			log.Infof("name=%s val=%s", tmpl.File, string(valCopy))
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
