package main

import (
	"fmt"
	"os"

	"github.com/charmbracelet/log"
	"gopkg.in/yaml.v3"
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
	cnf, err := os.Open("adsig.yml")
	if err != nil {
		return err
	}
	defer cnf.Close()

	var v struct {
		Connection struct {
			User     string
			Password string `yaml:",omitempty"`
		}
		Groups map[string]struct {
			AdFilter      string   `yaml:"adFilter,omitempty"`
			BaseDN        string   `yaml:"baseDN"`
			ExcludeEmails []string `yaml:"excludeEmails,omitempty"`
			Templates     []string
		}
		Templates map[string]struct {
			FieldMapping map[string]string `yaml:"fieldMapping"` // TODO: Maybe I will need to case the index? https://stackoverflow.com/questions/75535015/go-to-unmarshal-into-uppercase-keys
		}
	}

	ymlDecoder := yaml.NewDecoder(cnf)
	ymlDecoder.KnownFields(false)
	if err := ymlDecoder.Decode(&v); err != nil {
		return err
	}

	fmt.Printf("%+v", v)

	return nil
}
