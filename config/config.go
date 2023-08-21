package config

import (
	"io"

	"gopkg.in/yaml.v3"
)

// Type Definition
// ------------------------------------------------------------------------

type Connection struct {
	UserDN   string `yaml:"userDN"`   //nolint
	Password string `yaml:"password"` // TODO: Make this optional
}

type Group struct {
	AdFilter      string   `yaml:"adFilter"`      // TODO: Make this optional
	BaseDN        string   `yaml:"baseDN"`        //nolint
	ExcludeEmails []string `yaml:"excludeEmails"` // TODO: Make this optional
	Templates     []string
}

type Template struct {
	// TODO: Maybe I will need to case the keys? https://stackoverflow.com/questions/75535015/go-to-unmarshal-into-uppercase-keys
	Fields map[string]string `yaml:"fields"`
}

type Config struct {
	Connection Connection
	Groups     map[string]Group
	Templates  map[string]Template
}

// Public Functions
// ------------------------------------------------------------------------

// TODO: Validate and set defaults
func FromYAML(r io.Reader) (Config, error) {
	var cnf Config

	decoder := yaml.NewDecoder(r)
	decoder.KnownFields(true)

	err := decoder.Decode(&cnf)

	return cnf, err
}
