package main

import (
	"log"
	"net/http"
	"os"

	"github.com/gotestbootcamp/go-todo-app/config"
	"github.com/gotestbootcamp/go-todo-app/controller"
	"github.com/gotestbootcamp/go-todo-app/ledger"
	"github.com/gotestbootcamp/go-todo-app/store"
	"github.com/gotestbootcamp/go-todo-app/store/fake"
)

func main() {
	cfg, err := config.FromFlags(os.Args[1:]...)
	if err != nil {
		log.Printf("error parsing flags: %v", err)
		os.Exit(0)
	}
	log.Printf("ready: configuration:\n%s", cfg.String())

	var st store.Storage
	if cfg.Redis.URL != "" {
		log.Printf("store: using backend \"redis\"")
		st, err = store.NewRedis(cfg.Redis.URL, cfg.Redis.Password, cfg.Redis.Database)
	} else {
		log.Printf("store: using backend \"fake\"")
		st, err = fake.NewMem()
	}
	if err != nil {
		log.Printf("error creating store backend: %v", err)
	}
	log.Printf("ready: store backend")

	ldg, err := ledger.New(st)
	if err != nil {
		log.Printf("error parsing flags: %v", err)
	}
	log.Printf("ready: data ledger")

	ctrl := controller.New(ldg)
	log.Printf("ready: controller")

	log.Printf("start serving on address %q", cfg.Address)
	log.Fatal(http.ListenAndServe(cfg.Address, ctrl))
}
