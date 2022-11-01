package main

import (
	"log"
	"os"

	"github.com/urfave/cli/v2"

	"rederinghub.io/tools/protoc"
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
