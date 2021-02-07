package main

import (
	"bufio"
	"bytes"
	"fmt"
	"github.com/egorlichenkoam/bgo3/pkg/card"
	"github.com/egorlichenkoam/bgo3/pkg/dailycurrencies"
	"github.com/egorlichenkoam/bgo3/pkg/money"
	"github.com/egorlichenkoam/bgo3/pkg/person"
	"github.com/egorlichenkoam/bgo3/pkg/transaction"
	"io"
	"io/ioutil"
	"log"
	"net"
	"os"
	"strconv"
	"strings"
	"time"
)

const CRLF = "\r\n"

var GTransactions []*transaction.Transaction = nil
var GStandard map[*person.Person]map[transaction.Mcc]money.Money = nil
var GPers *person.Person = nil

func main() {
	printVersion()
	dailyCurrenciesSvc := dailycurrencies.NewService()
	err := dailyCurrenciesSvc.Extract()
	if err != nil {
		log.Fatal(err)
	}
	GTransactions, GStandard, GPers = transaction.GenerateTestData()
	if err := execute(); err != nil {
		log.Fatal(err)
	}
}

func execute() (err error) {
	listener, err := net.Listen("tcp", "0.0.0.0:9999")
	if err != nil {
		log.Fatal(err)
	}
	defer func(c io.Closer) {
		if cerr := c.Close(); cerr != nil {
			log.Fatal(cerr)
		}
	}(listener)
	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Println(err)
			continue
		}
		go handle(conn)
	}
}

func handle(conn net.Conn) {
	defer func() {
		if cerr := conn.Close(); cerr != nil {
			doFatal(cerr)
		}
	}()
	reader := bufio.NewReader(conn)
	const delimiter = '\n'
	line, err := reader.ReadString(delimiter)
	if err != nil {
		if err != io.EOF {
			doFatal(err)
		}
		return
	}
	log.Printf("received: %s\n", line)

	parts := strings.Split(line, " ")
	if len(parts) != 3 {
		log.Printf("invalide request line: %s", line)
	}

	path := parts[1]
	switch path {
	case "/":
		err = writeIndex(conn)
	case "/operations.csv":
		err = writeOperation(conn, "csv")
	case "/operations.xml":
		err = writeOperation(conn, "xml")
	case "/operations.json":
		err = writeOperation(conn, "json")
	default:
		err = write404(conn)
	}
	if err != nil {
		log.Print(err)
	}
}

func writeIndex(conn net.Conn) error {
	username := GPers.Name
	balance := money.Money(0)
	for _, card := range GPers.Cards {
		balance = balance + card.Balance
	}
	page, err := ioutil.ReadFile("web/template/index.html")
	page = bytes.ReplaceAll(page, []byte("{username}"), []byte(username))
	page = bytes.ReplaceAll(page, []byte("{balance}"), []byte(strconv.Itoa(int(balance))))
	writer := bufio.NewWriter(conn)
	writer.WriteString("HTTP/1.1 200" + CRLF)
	writer.WriteString("Content-Type: text/html;charset=utf-8" + CRLF)
	writer.WriteString(fmt.Sprintf("Content-Length: %d", len(page)) + CRLF)
	writer.WriteString("Connection: close" + CRLF + CRLF)
	writer.Write(page)
	err = writer.Flush()
	if err != nil {
		doFatal(err)
	}
	return err
}

func writeOperation(conn net.Conn, format string) error {
	var page []byte
	contentType := "Content-Type: %s;charset=utf-8" + CRLF
	contentLength := "Content-Length: %d" + CRLF
	switch format {
	case "csv":
		page = transaction.ExportCsvToBytes(GTransactions)
		contentType = fmt.Sprintf(contentType, "text/csv")
	case "xml":
		page = transaction.ExportXmlToBytes(GTransactions)
		contentType = fmt.Sprintf(contentType, "application/xml")
	case "json":
		page = transaction.ExportJsonToBytes(GTransactions)
		contentType = fmt.Sprintf(contentType, "application/json")
	default:
		page = []byte("Что-то пошло не так... Ёпрст\n")
		contentType = fmt.Sprintf(contentType, "plain/text")
	}
	writer := bufio.NewWriter(conn)
	writer.WriteString("HTTP/1.1 200" + CRLF)
	writer.WriteString(contentType)
	writer.WriteString(fmt.Sprintf(contentLength, len(page)))
	writer.WriteString("Connection: close" + CRLF + CRLF)
	nn, err := writer.Write(page)
	log.Print(nn)
	if err != nil {
		log.Print(err)
	}
	err = writer.Flush()
	return err
}

func write404(conn net.Conn) error {
	page, err := ioutil.ReadFile("web/template/404.html")
	writer := bufio.NewWriter(conn)
	writer.WriteString("HTTP/1.1 404" + CRLF)
	writer.WriteString("Content-Type: text/html;charset=utf-8" + CRLF)
	writer.WriteString(fmt.Sprintf("Content-Length: %d", len(page)) + CRLF)
	writer.WriteString("Connection: close" + CRLF + CRLF)
	writer.Write(page)
	err = writer.Flush()
	return err
}

func doFatal(err error) {
	log.Fatal(err)
}

func printCards(cards []card.Card) {
	for _, c := range cards {
		fmt.Println(c)
	}
	fmt.Println("")
}

func printTransactions(txs []*transaction.Transaction) {
	for _, tx := range txs {
		fmt.Println(tx, tx.CardId)
	}
	fmt.Println("")
}

func sumConcurrently() {
	cardSvc := card.NewService("510621", "BABANK")
	transactionSvc := transaction.NewService()
	card00 := cardSvc.Create(1000_000_00, card.Rub, "5106212879499054")

	tx := transactionSvc.CreateTransaction(1_000_00, "", card00, transaction.From)
	tx.Datetime = time.Date(2020, 12, 1, 0, 0, 0, 0, time.Local).Unix()
	tx = transactionSvc.CreateTransaction(12_000_00, "", card00, transaction.From)
	tx.Datetime = time.Date(2020, 12, 5, 0, 0, 0, 0, time.Local).Unix()

	tx = transactionSvc.CreateTransaction(10_000_00, "", card00, transaction.From)
	tx.Datetime = time.Date(2020, 11, 1, 0, 0, 0, 0, time.Local).Unix()
	tx = transactionSvc.CreateTransaction(22_000_00, "", card00, transaction.From)
	tx.Datetime = time.Date(2020, 11, 5, 0, 0, 0, 0, time.Local).Unix()

	tx = transactionSvc.CreateTransaction(100_00, "", card00, transaction.From)
	tx.Datetime = time.Date(2020, 9, 1, 0, 0, 0, 0, time.Local).Unix()
	tx = transactionSvc.CreateTransaction(200_00, "", card00, transaction.From)
	tx.Datetime = time.Date(2020, 9, 5, 0, 0, 0, 0, time.Local).Unix()

	tx = transactionSvc.CreateTransaction(800_000_00, "", card00, transaction.From)
	tx.Datetime = time.Date(2020, 6, 1, 0, 0, 0, 0, time.Local).Unix()
	tx = transactionSvc.CreateTransaction(2_000_000_00, "", card00, transaction.From)
	tx.Datetime = time.Date(2020, 6, 5, 0, 0, 0, 0, time.Local).Unix()

	tx = transactionSvc.CreateTransaction(8_700_000_00, "", card00, transaction.From)
	tx.Datetime = time.Date(2020, 3, 1, 0, 0, 0, 0, time.Local).Unix()

	tx = transactionSvc.CreateTransaction(3_000_000_00, "", card00, transaction.From)
	tx.Datetime = time.Date(2020, 4, 5, 0, 0, 0, 0, time.Local).Unix()

	tx = transactionSvc.CreateTransaction(1_000_000_00, "", card00, transaction.From)
	tx.Datetime = time.Date(2020, 1, 5, 0, 0, 0, 0, time.Local).Unix()

	result := transactionSvc.SumConcurrentlyByCardAndYearMonth(card00, time.Date(2020, 1, 1, 0, 0, 0, 0, time.Local).Unix(), time.Date(2020, 12, 1, 0, 0, 0, 0, time.Local).Unix(), transaction.From)

	fmt.Println("------------------------------------------------------------------ ")
	keys := make([]time.Time, 0)
	result.Range(func(key, value interface{}) bool {
		k, _ := key.(time.Time)
		keys = append(keys, k)
		return true
	})
	for _, key := range keys {
		value, _ := result.Load(key)
		fmt.Println(key, " - ", value)
	}
	fmt.Println("------------------------------------------------------------------")
}

func exportImport() {
	cardSvc := card.NewService("510621", "BABANK")
	transactionSvc := transaction.NewService()
	card00 := cardSvc.Create(1000_000_00, card.Rub, "5106212879499054")

	tx := transactionSvc.CreateTransaction(1_000_00, "5812", card00, transaction.From)
	tx.Datetime = time.Date(2020, 12, 1, 0, 0, 0, 0, time.Local).Unix()
	tx = transactionSvc.CreateTransaction(12_000_00, "5812", card00, transaction.From)
	tx.Datetime = time.Date(2020, 12, 5, 0, 0, 0, 0, time.Local).Unix()

	tx = transactionSvc.CreateTransaction(10_000_00, "5812", card00, transaction.From)
	tx.Datetime = time.Date(2020, 11, 1, 0, 0, 0, 0, time.Local).Unix()
	tx = transactionSvc.CreateTransaction(22_000_00, "5812", card00, transaction.From)
	tx.Datetime = time.Date(2020, 11, 5, 0, 0, 0, 0, time.Local).Unix()

	tx = transactionSvc.CreateTransaction(100_00, "5812", card00, transaction.From)
	tx.Datetime = time.Date(2020, 9, 1, 0, 0, 0, 0, time.Local).Unix()
	tx = transactionSvc.CreateTransaction(200_00, "5812", card00, transaction.From)
	tx.Datetime = time.Date(2020, 9, 5, 0, 0, 0, 0, time.Local).Unix()

	tx = transactionSvc.CreateTransaction(800_000_00, "5812", card00, transaction.From)
	tx.Datetime = time.Date(2020, 6, 1, 0, 0, 0, 0, time.Local).Unix()
	tx = transactionSvc.CreateTransaction(2_000_000_00, "5812", card00, transaction.From)
	tx.Datetime = time.Date(2020, 6, 5, 0, 0, 0, 0, time.Local).Unix()

	tx = transactionSvc.CreateTransaction(8_700_000_00, "5812", card00, transaction.From)
	tx.Datetime = time.Date(2020, 3, 1, 0, 0, 0, 0, time.Local).Unix()

	tx = transactionSvc.CreateTransaction(3_000_000_00, "5812", card00, transaction.From)
	tx.Datetime = time.Date(2020, 4, 5, 0, 0, 0, 0, time.Local).Unix()

	tx = transactionSvc.CreateTransaction(1_000_000_00, "5812", card00, transaction.From)
	tx.Datetime = time.Date(2020, 1, 5, 0, 0, 0, 0, time.Local).Unix()

	log.Println("CSV")
	err := transaction.ExportCsv(transactionSvc.Transactions)
	if err != nil {
		log.Println(err)
	} else {
		path, _ := os.Getwd()
		path = path + "/exports.csv"
		txs, err := transaction.ImportCsv(path)
		if err != nil {
			log.Println(err)
		} else {
			for _, tx := range txs {
				log.Println(tx)
			}
		}
	}

	log.Println("JSON")
	err = transaction.ExportJson(transactionSvc.Transactions)
	if err != nil {
		log.Println(err)
	} else {
		path, _ := os.Getwd()
		path = path + "/exports.json"
		txs, err := transaction.ImportJson(path)
		if err != nil {
			log.Println(err)
		} else {
			for _, tx := range txs {
				log.Println(tx)
			}
		}
	}

	log.Println("XML")
	err = transaction.ExportXml(transactionSvc.Transactions)
	if err != nil {
		log.Println(err)
	} else {
		path, _ := os.Getwd()
		path = path + "/exports.xml"
		txs, err := transaction.ImportXml(path)
		if err != nil {
			log.Println(err)
		} else {
			for _, tx := range txs {
				log.Println(tx)
			}
		}
	}
}

func printVersion() {
	fmt.Println("03.02.01")
}
