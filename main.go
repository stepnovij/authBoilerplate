package main

import (
	"flag"
	"log"
	"os"

	"github.com/stepnovij/authBoilerplate/daemon"
)

func processFlags() *daemon.Config {
	cfg := &daemon.Config{}
	connectionString := os.Getenv("CONNECTION_STRING")
	listenAddr := os.Getenv("LISTEN_ADDR")
	flag.StringVar(&cfg.ListenSpec, "listen", listenAddr, "HTTP listen spec")
	flag.StringVar(&cfg.Db.ConnectString, "db-connect", connectionString, "DB connect String")
	flag.Parse()
	return cfg
}

func main() {

	cfg := processFlags()
	if err := daemon.Run(cfg); err != nil {
		log.Printf("Error in main(): %v", err)
	}

}
