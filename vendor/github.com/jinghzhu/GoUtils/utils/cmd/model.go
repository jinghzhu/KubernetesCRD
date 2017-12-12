package cmd

import "time"

type Mount struct {
	Name   string `string:"name,omitempty"`
	Type   string `string:"type,omitempty"`
	Share  string `string:"share"`
	Server string `string:"server"`
}

const (
	ShowmountTimeout = 15 * time.Second
	CmdShowmount     = "showmount"
	CmdShowmountOptE = "-e"
)
