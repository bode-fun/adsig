package main

import (
	"errors"
	"flag"
	"os"
	"path/filepath"
	"strings"

	"git.bode.fun/adsig"
	"git.bode.fun/adsig/config"
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
	var account string

	flag.StringVar(&account, "account", "", "The account to create the signature for")

	flag.Parse()

	if account == "" {
		return errors.New("account name is required but not provided")
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

	conn, err := ldap.DialURL(cnf.Connection.Address)
	if err != nil {
		return err
	}
	defer conn.Close()

	err = conn.Bind(cnf.Connection.UserDN, cnf.Connection.Password)
	if err != nil {
		return err
	}
	defer conn.Unbind()

	groups, err := adsig.GroupsFromConfig(cnf, conn)
	if err != nil {
		return err
	}

	renderedTmpls := make(map[string]string, 0)

	for _, group := range groups {
		if ok, member := group.MemberBySamAccountName(account); ok {
			log.Infof("%s is a member of the group \"%s\" with %d members", account, group.Name, len(group.Members))

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

					ext := filepath.Ext(tmpl.Name())

					renderedTmpls[sig.Name+ext] = b.String()
				}
			}
		}
	}

	for name, tmplData := range renderedTmpls {
		log.Infof("name: %s, data: %s", name, tmplData)
	}

	return nil
}
