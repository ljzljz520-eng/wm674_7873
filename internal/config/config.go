package config

import "strings"

type Config struct {
	DatabasePath string
	Actor        string
	Command      string
}

func Default() Config { return Config{DatabasePath: "meters.db", Actor: "operator", Command: "demo"} }
func (c Config) Normalize() Config {
	c.DatabasePath = strings.TrimSpace(c.DatabasePath)
	if c.DatabasePath == "" {
		c.DatabasePath = "meters.db"
	}
	c.Actor = strings.TrimSpace(c.Actor)
	if c.Actor == "" {
		c.Actor = "operator"
	}
	c.Command = strings.TrimSpace(c.Command)
	if c.Command == "" {
		c.Command = "demo"
	}
	return c
}
