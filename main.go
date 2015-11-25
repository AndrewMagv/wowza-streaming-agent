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
		addr = c.String("addr")

		stop = make(chan struct{})
	)

	nodeInfo, err := getNodeInfo()
	if err != nil {
		log.Warning(err)
		os.Exit(1)
	}

	service.Advertise = c.String("advertise")
	if service.Advertise == "" {
		if adver, ok := nodeInfo["public_ipv4"]; !ok {
			cli.ShowAppHelp(c)
			log.Warning(err)
			log.Error("Required flag --advertise missing")
			os.Exit(1)
		} else {
			service.Advertise = adver
		}
	}
	if host, ok := nodeInfo["host"]; ok {
		service.Host = host
	} else {
		log.Error("Required info public host endpoint required")
		os.Exit(1)
	}
	if node, ok := nodeInfo["node"]; ok {
		service.Node = node
	} else {
		log.Error("Required info node identity required")
		os.Exit(1)
	}

	if addr != "" {
		log.WithFields(log.Fields{"addr": addr, "advertise": service.Advertise}).Info("API endpoint begin")
		go runAPIEndpoint(addr, stop)
		<-stop // we should never reach pass this point
	} else {
		log.Warning("API endpoint disabled")
	}

	log.Warning("nothing to do; quit now")
}
