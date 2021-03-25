package main

import (
	"fmt"
	"github.com/egorlichenkoam/bgo3/cmd/bank/app"
	"github.com/egorlichenkoam/bgo3/pkg/card"
	"github.com/egorlichenkoam/bgo3/pkg/person"
	"github.com/egorlichenkoam/bgo3/pkg/transaction"
	"log"
	"net"
	"net/http"
	"os"
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
		ReadTimeout:       5 * time.Second,
		ReadHeaderTimeout: 5 * time.Second,
		WriteTimeout:      5 * time.Second,
	}
	return server.ListenAndServe()
}

const defaultHost = "0.0.0.0"
const defaultPort = "9999"
const name = "bank"

func createTestData() (*person.Service, *card.Service, *transaction.Service) {
	log.Printf("%s - %s", name, "Start createTestData")
	filePath, _ := os.Getwd()
	personSvc := person.TestData()
	persons, err := person.ImportJson(filePath + "/personsExport.json")
	if err != nil {
		log.Println(err)
		person.ExportJson(personSvc.Persons)
	} else {
		personSvc.Persons = persons
	}
	for _, p := range personSvc.Persons {
		log.Printf("%s - %s : %d", name, "available person", p.Id)
	}

	cardSvc := card.TestData(personSvc.Ids())
	cards, err := card.ImportJson(filePath + "/cardsExport.json")
	if err != nil {
		log.Println(err)
		card.ExportJson(cardSvc.Cards)
	} else {
		cardSvc.Cards = cards
	}

	txSvc := transaction.TestData(cardSvc.Ids())
	txs, err := transaction.ImportJson(filePath + "/txsExport.json")
	if err != nil {
		log.Println(err)
		transaction.ExportJson(txSvc.Transactions)
	} else {
		txSvc.Transactions = txs
	}
	log.Printf("%s - %s", name, "End createTestData")

	personSvc.Persons = append(personSvc.Persons, &person.Person{
		Id:   111,
		Name: "Пустой пользователь",
	})

	return personSvc, cardSvc, txSvc
}

func main() {
	printVersion()
	log.Printf("%s - %s", name, "main - start")
	personSvc, cardSvc, txSvc := createTestData()
	host, ok := os.LookupEnv("HOST")
	if !ok {
		host = defaultHost
	}
	port, ok := os.LookupEnv("PORT")
	if !ok {
		port = defaultPort
	}
	if err := Execute(net.JoinHostPort(host, port), personSvc, cardSvc, txSvc); err != nil {
		log.Println(err)
		os.Exit(1)
	}
	log.Printf("%s - %s", name, "main - end")
}

func printVersion() {
	fmt.Println("04.05")
}
