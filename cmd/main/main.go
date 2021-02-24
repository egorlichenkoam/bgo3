package main

import (
	"fmt"
	"github.com/egorlichenkoam/bgo3/cmd/bank"
	"github.com/egorlichenkoam/bgo3/pkg/card"
	"github.com/egorlichenkoam/bgo3/pkg/person"
	"github.com/egorlichenkoam/bgo3/pkg/transaction"
	"log"
	"net"
	"os"
)

//const CRLF = "\r\n"

const defaultHost = "0.0.0.0"
const defaultPort = "9999"
const name = "main"

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
	if err := bank.Execute(net.JoinHostPort(host, port), personSvc, cardSvc, txSvc); err != nil {
		log.Println(err)
		os.Exit(1)
	}
	log.Printf("%s - %s", name, "main - end")
}

//func execute() (err error) {
//	listener, err := net.Listen("tcp", "0.0.0.0:9999")
//	if err != nil {
//		log.Fatal(err)
//	}
//	defer func(c io.Closer) {
//		if cErr := c.Close(); cErr != nil {
//			log.Fatal(cErr)
//		}
//	}(listener)
//	for {
//		conn, err := listener.Accept()
//		if err != nil {
//			log.Println(err)
//			continue
//		}
//		go handle(conn)
//	}
//}
//
//func handle(conn net.Conn) {
//	defer func() {
//		if cErr := conn.Close(); cErr != nil {
//			doFatal(cErr)
//		}
//	}()
//	reader := bufio.NewReader(conn)
//	const delimiter = '\n'
//	line, err := reader.ReadString(delimiter)
//	if err != nil {
//		if err != io.EOF {
//			doFatal(err)
//		}
//		return
//	}
//	log.Printf("received: %s\n", line)
//
//	parts := strings.Split(line, " ")
//	if len(parts) != 3 {
//		log.Printf("invalide request line: %s", line)
//	}
//
//	path := parts[1]
//	switch path {
//	case "/":
//		err = writeIndex(conn)
//	case "/operations.csv":
//		err = writeOperation(conn, "csv")
//	case "/operations.xml":
//		err = writeOperation(conn, "xml")
//	case "/operations.json":
//		err = writeOperation(conn, "json")
//	default:
//		err = write404(conn)
//	}
//	if err != nil {
//		log.Print(err)
//	}
//}
//
//func writeIndex(conn net.Conn) error {
//	username := GPerson.Name
//	balance := money.Money(0)
//	for _, c := range GCardSvc.ByPersonId(GPerson.Id) {
//		balance = balance + c.Balance
//	}
//	page, err := ioutil.ReadFile("web/template/index.html")
//	page = bytes.ReplaceAll(page, []byte("{username}"), []byte(username))
//	page = bytes.ReplaceAll(page, []byte("{balance}"), []byte(strconv.Itoa(int(balance))))
//	writer := bufio.NewWriter(conn)
//	writer.WriteString("HTTP/1.1 200" + CRLF)
//	writer.WriteString("Content-Type: text/html;charset=utf-8" + CRLF)
//	writer.WriteString(fmt.Sprintf("Content-Length: %d", len(page)) + CRLF)
//	writer.WriteString("Connection: close" + CRLF + CRLF)
//	writer.Write(page)
//	err = writer.Flush()
//	if err != nil {
//		doFatal(err)
//	}
//	return err
//}
//
//func writeOperation(conn net.Conn, format string) error {
//	var page []byte
//	contentType := "Content-Type: %s;charset=utf-8" + CRLF
//	contentLength := "Content-Length: %d" + CRLF
//	switch format {
//	case "csv":
//		page = transaction.ExportCsvToBytes(GTransactionSvc.Transactions)
//		contentType = fmt.Sprintf(contentType, "text/csv")
//	case "xml":
//		page = transaction.ExportXmlToBytes(GTransactionSvc.Transactions)
//		contentType = fmt.Sprintf(contentType, "application/xml")
//	case "json":
//		page = transaction.ExportJsonToBytes(GTransactionSvc.Transactions)
//		contentType = fmt.Sprintf(contentType, "application/json")
//	default:
//		page = []byte("Что-то пошло не так... Ёпрст\n")
//		contentType = fmt.Sprintf(contentType, "plain/text")
//	}
//	writer := bufio.NewWriter(conn)
//	writer.WriteString("HTTP/1.1 200" + CRLF)
//	writer.WriteString(contentType)
//	writer.WriteString(fmt.Sprintf(contentLength, len(page)))
//	writer.WriteString("Connection: close" + CRLF + CRLF)
//	_, err := writer.Write(page)
//	if err != nil {
//		log.Print(err)
//	}
//	err = writer.Flush()
//	return err
//}
//
//func write404(conn net.Conn) error {
//	page, err := ioutil.ReadFile("web/template/404.html")
//	writer := bufio.NewWriter(conn)
//	writer.WriteString("HTTP/1.1 404" + CRLF)
//	writer.WriteString("Content-Type: text/html;charset=utf-8" + CRLF)
//	writer.WriteString(fmt.Sprintf("Content-Length: %d", len(page)) + CRLF)
//	writer.WriteString("Connection: close" + CRLF + CRLF)
//	writer.Write(page)
//	err = writer.Flush()
//	return err
//}
//
//func doFatal(err error) {
//	log.Fatal(err)
//}

func printVersion() {
	fmt.Println("03.04.02")
}
