package daemon

import (
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"

	"github.com/stepnovij/authBoilerplate/db"
	"github.com/stepnovij/authBoilerplate/model"
	"github.com/stepnovij/authBoilerplate/view"
)

type Config struct {
	ListenSpec string

	Db db.Config
}

func Run(cfg *Config) error {
	log.Printf("Starting, HTTP on: %s\n", cfg.ListenSpec)

	db, err := db.InitDb(cfg.Db)
	if err != nil {
		log.Printf("Error initializing database: %v\n", err)
		return err
	}
	log.Printf("Initializing database is OK")

	m := model.New(db)
	log.Printf("Connecting models")

	l, err := net.Listen("tcp", cfg.ListenSpec)
	if err != nil {
		log.Printf("Error creating listener: %v\n", err)
		return err
	}
	view.Start(m, l)
	waitForSignal()
	log.Printf("%v  %v", m, l)
	return nil
}

func waitForSignal() {
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, syscall.SIGINT, syscall.SIGTERM)
	s := <-ch
	log.Printf("Got signal: %v, exiting.", s)
}
