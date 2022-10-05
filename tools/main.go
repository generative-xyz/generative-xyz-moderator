package main

import (
	"log"
	"os"
	"rederinghub.io/tools/protoc"

	"github.com/urfave/cli/v2"
)

func main() {
	app := &cli.App{
		Commands: []*cli.Command{{
			Name:   "protoc",
			Usage:  "Gen from proto files",
			Action: protoc.Action,
		}},
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
