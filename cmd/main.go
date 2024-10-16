package main

import (
	"log"
	"net/http"
	"os"

	"github.com/ffromani/go-todo-app/config"
	"github.com/ffromani/go-todo-app/controller"
	"github.com/ffromani/go-todo-app/ledger"
	"github.com/ffromani/go-todo-app/store/fake"
)

func main() {
	cfg, err := config.FromFlags(os.Args...)
	if err != nil {
		log.Printf("error parsing flags: %v", err)
	}
	log.Printf("ready: configuration:\n%s", cfg.String())

	st, err := fake.NewMem()
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
