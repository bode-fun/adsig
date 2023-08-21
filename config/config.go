package config

import (
	"errors"
	"io"

	"gopkg.in/yaml.v3"
)

// Errors
// ------------------------------------------------------------------------

// TODO: Make a custom error type to wrap in.

var (
	ErrDecode    = errors.New("config: Can not content of decode config")
	ErrNoContent = errors.New("config: The content file holds no content")
)

// Type Definition
// ------------------------------------------------------------------------

type Server struct {
	Host         string `yaml:"host"`
	Port         int    `yaml:"port"`
	ReadTimeout  int    `yaml:"readTimeout"`
	WriteTimeout int    `yaml:"writeTimeout"`
}

type Connection struct {
	UserDN   string `yaml:"userDN"`   //nolint
	Password string `yaml:"password"` // TODO: Make this optional.
}

type Group struct {
	AdFilter      string   `yaml:"adFilter"`      // TODO: Make this optional.
	BaseDN        string   `yaml:"baseDN"`        //nolint
	ExcludeEmails []string `yaml:"excludeEmails"` // TODO: Make this optional.
	Templates     []string
}

type Template struct {
	Fields map[string]string `yaml:"fields"`
}

type Config struct {
	Server     Server              `yaml:"server"`
	Connection Connection          `yaml:"connection"`
	Groups     map[string]Group    `yaml:"groups"`
	Templates  map[string]Template `yaml:"templates"`
}

// Public Functions
// ------------------------------------------------------------------------

// TODO: Validate and set defaults.
func FromYAML(r io.Reader) (Config, error) {
	var cnf Config

	decoder := yaml.NewDecoder(r)
	decoder.KnownFields(true)

	// err can be io.EOF or yaml.TypeError
	if err := decoder.Decode(&cnf); err != nil {
		if errors.Is(err, &yaml.TypeError{}) { //nolint
			// TODO: Wrap the error.
			return Config{}, ErrDecode
		} else if errors.Is(err, io.EOF) {
			return Config{}, ErrNoContent
		} else {
			return Config{}, ErrDecode
		}
	}

	return cnf, nil
}
