package main

import (
	"os/exec"
	"time"
)

// Config - provider configuration
type Config struct {
	Host        string
	User        string
	KeyFilePath string
	Timeout     time.Duration

	cmds []*exec.Cmd
}

// AddCmd adds a *execCmd to the list of commands the provider is responsible for
func (c *Config) AddCmd(cmd *exec.Cmd) {
	c.cmds = append(c.cmds, cmd)
}

// KillAll is a helper function that kills all ssh processes that were started by this provider
func (c *Config) KillAll() {
	for _, cmd := range c.cmds {
		cmd.Process.Kill()
	}
}
