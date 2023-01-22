package main

import (
	"log"

	"github.com/mahditakrim/template/config"
	"github.com/mahditakrim/template/internal/logger"
	"github.com/mahditakrim/template/luncher"
	"github.com/mahditakrim/template/setup"
)

func main() {

	conf, err := config.Init()
	if err != nil {
		log.Panic(err)
	}

	f, err := logger.Init(conf.LogPath)
	if err != nil {
		log.Panic(err)
	}
	defer f.Close()

	runnables, err := setup.Init(conf)
	if err != nil {
		log.Panic(err)
	}

	luncher.Start(runnables...)
}
