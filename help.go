package main

import (
	cli "github.com/codegangsta/cli"
)

func init() {
	cli.AppHelpTemplate = `Usage: {{.Name}} {{if .Flags}}[OPTIONS]{{end}}

{{.Usage}}

Version: {{.Version}}

{{if .Flags}}Options:
	{{range .Flags}}{{.}}
	{{end}}{{end}}
`
}
