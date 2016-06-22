package main

import (
	log "github.com/Sirupsen/logrus"
	"github.com/jawher/mow.cli"
	"os"
)

func init() {
	log.SetFormatter(new(log.JSONFormatter))
}

func main() {
	app := cli.App("curated-authors-transformer", "A RESTful API for transforming Bertha Curated Authors to UP People JSON")

	baseURL := app.String(cli.StringOpt{
		Name:   "base-url",
		Value:  "http://localhost:8080/transformers/curated-authors/",
		Desc:   "Base url",
		EnvVar: "BASE_URL",
	})
	port := app.Int(cli.IntOpt{
		Name:   "port",
		Value:  8080,
		Desc:   "Port to listen on",
		EnvVar: "PORT",
	})
	berthaSrcUrl := app.String(cli.StringOpt{
		Name:   "bertha-source-url",
		Value:  "http://bertha.ig.ft.com/view/publish/gss/1wEdVRLtayZ6-XBfYM3vKAGaOV64cNJD3L8MlLM8-uFY/TestAuthors",
		Desc:   "The URL of the Bertha Authors JSON source",
		EnvVar: "BERTHA_SOURCE_URL",
	})

	app.Action = func() {
		log.Info("App started!!!")
		log.Info(*baseURL)
		log.Info(*port)
		log.Info(*berthaSrcUrl)
	}

	app.Run(os.Args)
}
