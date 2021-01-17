package main

import (
	"01/pkg/card"
	"01/pkg/transaction"
	"01/pkg/transfer"
	"fmt"
)

func main() {
	cardSvc := card.NewService("510621")
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

	cardSvc.NewCard("BANK", 10_000_00, card.Rub, "5106212879499054")
	cardSvc.NewCard("BANK", 20_000_00, card.Rub, "5106212548197220")
	cardSvc.NewCard("BANK", 30_000_00, card.Rub, "5106211562724463")

	printCards(cardSvc.Cards)
	printTransactions(transactionSvc.Transactions)

	fmt.Println(transferSvc.Card2Card(transferSvc.CardSvc.Cards[0].Number, transferSvc.CardSvc.Cards[1].Number, 1_000_00))
	fmt.Println(transferSvc.Card2Card(transferSvc.CardSvc.Cards[1].Number, transferSvc.CardSvc.Cards[2].Number, 1_000_00))
	fmt.Println(transferSvc.Card2Card(transferSvc.CardSvc.Cards[2].Number, transferSvc.CardSvc.Cards[0].Number, 1_000_00))

	fmt.Println("")

	printCards(cardSvc.Cards)
	printTransactions(transactionSvc.Transactions)
}

func printCards(cards []card.Card) {
	for _, c := range cards {
		fmt.Println(c)
	}
	fmt.Println("")
}

func printTransactions(txs []transaction.Transaction) {
	for _, tx := range txs {
		fmt.Println(tx, tx.Card.Number)
	}
	fmt.Println("")
}

func printVersion() {
	fmt.Println("02.03.01")
}
