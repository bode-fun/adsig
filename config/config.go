package config

import (
	"errors"
	"io"
	"strings"
	"time"

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

type configMember interface {
	setDefaults()
	normalize()
}

// Server
// ------------------------------------------------------------------------

var _ configMember = (*Server)(nil)

type Server struct {
	Host string `yaml:"host"`
	Port int    `yaml:"port"`
	// Optional
	ReadTimeout time.Duration `yaml:"readTimeout,omitempty"`
	// Optional
	WriteTimeout time.Duration `yaml:"writeTimeout,omitempty"`
}

func (s *Server) setDefaults() {
	var defaultTimeout time.Duration = 30

	if s.ReadTimeout == 0 {
		s.ReadTimeout = defaultTimeout
	}

	if s.WriteTimeout == 0 {
		s.WriteTimeout = defaultTimeout
	}
}

func (s *Server) normalize() {}

// Connection
// ------------------------------------------------------------------------

var _ configMember = (*Connection)(nil)

type Connection struct {
	Address  string `yaml:"address"`
	UserDN   string `yaml:"userDN"`   //nolint
	Password string `yaml:"password"` // TODO: Make this optional and introduce enum.
}

func (c *Connection) setDefaults() {}
func (c *Connection) normalize()   {}

// Group
// ------------------------------------------------------------------------

var _ configMember = (*Group)(nil)

type Group struct {
	// Optional
	AdFilter string `yaml:"adFilter,omitempty"`
	BaseDN   string `yaml:"baseDN"` //nolint
	// Optional
	ExcludeEmails []string `yaml:"excludeEmails,omitempty"`
	Templates     []string
}

func (g *Group) setDefaults() {
	if g.AdFilter == "" {
		g.AdFilter = "(&(objectclass=person)(mail=*))"
	}
}

func (g *Group) normalize() {
	for i, email := range g.ExcludeEmails {
		g.ExcludeEmails[i] = strings.TrimSpace(strings.ToLower(email))
	}
}

// Template
// ------------------------------------------------------------------------

var _ configMember = (*Template)(nil)

type Template struct {
	Fields map[string]string `yaml:"fields"`
}

func (t *Template) setDefaults() {}
func (t *Template) normalize()   {}

// Config
// ------------------------------------------------------------------------

var _ configMember = (*Config)(nil)

type Config struct {
	Server     Server              `yaml:"server"`
	Connection Connection          `yaml:"connection"`
	Groups     map[string]Group    `yaml:"groups"`
	Templates  map[string]Template `yaml:"templates"`
}

func (c *Config) setDefaults() {
	c.Connection.setDefaults()
	c.Server.setDefaults()

	for gName, g := range c.Groups {
		g.setDefaults()
		c.Groups[gName] = g
	}

	for tName, t := range c.Templates {
		t.setDefaults()
		c.Templates[tName] = t
	}
}

func (c *Config) normalize() {
	c.Connection.normalize()
	c.Server.normalize()

	for gName, g := range c.Groups {
		g.normalize()
		c.Groups[gName] = g
	}

	for tName, t := range c.Templates {
		t.normalize()
		c.Templates[tName] = t
	}
}

// Public Functions
// ------------------------------------------------------------------------

// TODO: Validate.
func FromYAML(r io.Reader) (Config, error) {
	var cnf Config

	decoder := yaml.NewDecoder(r)
	decoder.KnownFields(true)

	// err can be io.EOF or yaml.TypeError
	if err := decoder.Decode(&cnf); err != nil {
		var yamlTypeErr yaml.TypeError
		if errors.Is(err, &yamlTypeErr) {
			// TODO: Wrap the error.
			return Config{}, ErrDecode
		} else if errors.Is(err, io.EOF) {
			return Config{}, ErrNoContent
		} else {
			return Config{}, ErrDecode
		}
	}

	cnf.setDefaults()
	cnf.normalize()

	return cnf, nil
}
