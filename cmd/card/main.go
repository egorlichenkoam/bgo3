package main

import (
	"01/pkg/card"
	"01/pkg/transaction"
	"01/pkg/transfer"
	"fmt"
	"log"
	"os"
	"time"
)

func main() {
	cardSvc := card.NewService("510621", "BANK")
	transactionSvc := transaction.NewService()
	commissions := transfer.Commission{
		PercentInBank:       0,
		MinimumInBank:       0,
		PercentToDiffBank:   0.5,
		MinimumToDiffBank:   10_00,
		PercentBetweenBanks: 1.5,
		MinimumBetweenBanks: 30_00,
	}
	transferSvc := transfer.NewService(cardSvc, transactionSvc, commissions)

	cardSvc.Create(10_000_00, card.Rub, "5106212879499054")
	cardSvc.Create(20_000_00, card.Rub, "5106212548197220")
	cardSvc.Create(30_000_00, card.Rub, "5106211562724463")

	printCards(cardSvc.Cards)
	printTransactions(transactionSvc.Transactions)

	fmt.Println(transferSvc.Card2Card(transferSvc.CardSvc.Cards[0].Number, transferSvc.CardSvc.Cards[1].Number, 1_000_00))
	fmt.Println(transferSvc.Card2Card(transferSvc.CardSvc.Cards[1].Number, transferSvc.CardSvc.Cards[2].Number, 1_000_00))
	fmt.Println(transferSvc.Card2Card(transferSvc.CardSvc.Cards[2].Number, transferSvc.CardSvc.Cards[0].Number, 1_000_00))

	fmt.Println("")

	printCards(cardSvc.Cards)
	printTransactions(transactionSvc.Transactions)

	sumConcurrently()

	exportImport()

	printVersion()
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
	fmt.Println("03.01.02.02")
}
