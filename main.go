package main

import (
	"flag"
	"log"
	// "net/http"
	"simpledex/daemon"
)

var assetsPath string

func processFlags() *daemon.Config {
	cfg := &daemon.Config{}

	flag.StringVar(&cfg.ListenSpec, "listen", "localhost:3001", "HTTP listen spec")
	flag.StringVar(&cfg.Db.ConnectString, "db-connect", "user=i.stepnov dbname=simpledex sslmode=disable", "DB connect String")
	flag.Parse()
	return cfg
}

func main() {

	cfg := processFlags()
	if err := daemon.Run(cfg); err != nil {
		log.Printf("Error in main(): %v", err)
	}

}