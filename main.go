package main

import (
	"github.com/andrewmagv/wowza-streaming-agent/api/service"

	log "github.com/Sirupsen/logrus"
	cli "github.com/codegangsta/cli"

	"os"
)

var (
	Flags = []cli.Flag{
		cli.StringFlag{
			Name:  "addr",
			Usage: "API endpoint for stream controller",
		},
		cli.StringFlag{
			Name:  "advertise",
			Usage: "The netloc of this node seen by other nodes",
		},
	}
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

func main() {
	app := cli.NewApp()
	app.Name = "wowza-agent"
	app.Usage = "Process stream adminitration to Wowza"
	app.Authors = []cli.Author{
		cli.Author{"Yi-Hung Jen", "yihungjen@gmail.com"},
	}
	app.Flags = Flags
	app.Action = Agent
	app.Run(os.Args)
}

func Agent(c *cli.Context) {
	var (
		addr      = c.String("addr")
		advertise = c.String("advertise")

		stop = make(chan struct{})
	)

	if advertise == "" {
		cli.ShowAppHelp(c)
		log.Error("Required flag --advertise missing")
		os.Exit(1)
	} else {
		service.Advertise = advertise
	}

	if addr != "" {
		log.WithFields(log.Fields{"addr": addr}).Info("API endpoint begin")
		go runAPIEndpoint(addr, stop)
		<-stop // we should never reach pass this point
	} else {
		log.Warning("API endpoint disabled")
	}

	log.Warning("nothing to do; quit now")
}
