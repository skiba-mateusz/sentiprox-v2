package config

import (
	"fmt"
	"os"
	"strings"

	"gopkg.in/yaml.v3"
)

const (
	StaticRoute RouteType = iota
	PassRoute
)

type RouteType int

type Route struct {
	Path    string    `yaml:"path"`
	Target  string    `yaml:"target"`
	Type    RouteType `yaml:"type"`
	MLCheck bool      `yaml:"ml_check"`
}

type Proxy struct {
	Addr string `yaml:"addr"`
}

type Database struct {
	Addr string `yaml:"addr"`
}

type ML struct {
	Addr             string  `yaml:"addr"`
	Score            float64 `yaml:"score"`
	TimeoutMS        int     `yaml:"timeout_ms"`
	FailOpen         bool    `yaml:"fail_open"`
	InspectBodyBytes int     `yaml:"inspect_body_bytes"`
}

type Config struct {
	Proxy    *Proxy
	Database *Database
	ML       *ML
	Routes   []Route `yaml:"routes"`
}

func Load(path string) (*Config, error) {
	f, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	var c Config
	if err = yaml.Unmarshal(f, &c); err != nil {
		return nil, err
	}

	return &c, nil
}

func (r *RouteType) UnmarshalYAML(value *yaml.Node) error {
	var s string
	if err := value.Decode(&s); err != nil {
		return err
	}

	switch strings.ToLower(s) {
	case "static":
		*r = StaticRoute
	case "pass":
		*r = PassRoute
	default:
		return fmt.Errorf("invalid route type: %s", s)
	}

	return nil
}
