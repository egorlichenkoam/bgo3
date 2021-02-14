package main

import (
	"bufio"
	"bytes"
	"fmt"
	"github.com/egorlichenkoam/bgo3/pkg/card"
	"github.com/egorlichenkoam/bgo3/pkg/dailycurrencies"
	"github.com/egorlichenkoam/bgo3/pkg/money"
	"github.com/egorlichenkoam/bgo3/pkg/person"
	"github.com/egorlichenkoam/bgo3/pkg/qrcodegenerator"
	"github.com/egorlichenkoam/bgo3/pkg/transaction"
	"io"
	"io/ioutil"
	"log"
	"net"
	"strconv"
	"strings"
)

const CRLF = "\r\n"

var GPersonSvc *person.Service = nil
var GCardSvc *card.Service = nil
var GTransactionSvc *transaction.Service = nil
var GStandard map[*person.Person]map[transaction.Mcc]money.Money = nil
var GPerson *person.Person = nil

func createTestData() {
	if (GPersonSvc == nil) || (GCardSvc == nil) || (GTransactionSvc == nil) || (GStandard == nil) || (GPerson == nil) {
		GPersonSvc, GCardSvc, GTransactionSvc, GStandard, GPerson = transaction.GenerateTestData()
	}
}

func main() {
	printVersion()
	createTestData()
	dailyCurrenciesSvc := dailycurrencies.NewService()
	err := dailyCurrenciesSvc.Extract()
	if err != nil {
		log.Fatal(err)
	}
	qrcodegeneratorSvc := qrcodegenerator.NewServive(3000)
	fp, err := qrcodegeneratorSvc.Encode("Привет с большого бодуна!", "qrcode.png")
	if err != nil {
		log.Fatal(err)
	}
	log.Println(fp)
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
	username := GPerson.Name
	balance := money.Money(0)
	for _, card := range GCardSvc.ByPersonId(GPerson.Id) {
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
		page = transaction.ExportCsvToBytes(GTransactionSvc.Transactions)
		contentType = fmt.Sprintf(contentType, "text/csv")
	case "xml":
		page = transaction.ExportXmlToBytes(GTransactionSvc.Transactions)
		contentType = fmt.Sprintf(contentType, "application/xml")
	case "json":
		page = transaction.ExportJsonToBytes(GTransactionSvc.Transactions)
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

func printVersion() {
	fmt.Println("03.03")
}
