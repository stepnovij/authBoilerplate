package daemon

import (
	"log"
	"net"
    "os"
    "os/signal"
    "syscall"
    "simpledex/db"
    "simpledex/model"
    "simpledex/view"
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
    ch := make(chan os.Signal)
    signal.Notify(ch, syscall.SIGINT, syscall.SIGTERM)
    s := <-ch
    log.Printf("Got signal: %v, exiting.", s)
}
