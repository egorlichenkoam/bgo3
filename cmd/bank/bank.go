package bank

import (
	"github.com/egorlichenkoam/bgo3/cmd/bank/app"
	"github.com/egorlichenkoam/bgo3/pkg/card"
	"github.com/egorlichenkoam/bgo3/pkg/person"
	"github.com/egorlichenkoam/bgo3/pkg/transaction"
	"log"
	"net/http"
	"time"
)

//starts bank server at given address
func Execute(addr string, personSvc *person.Service, cardSvc *card.Service, txSvc *transaction.Service) error {
	log.Println("Execute starting bank server")
	mux := http.NewServeMux()
	application := app.NewServer(personSvc, cardSvc, txSvc, mux)
	application.Init()
	server := &http.Server{
		Addr:              addr,
		Handler:           application,
		ReadTimeout:       1 * time.Second,
		ReadHeaderTimeout: 1 * time.Second,
		WriteTimeout:      1 * time.Second,
	}
	return server.ListenAndServe()
}
