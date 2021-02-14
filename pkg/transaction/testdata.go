package transaction

import (
	"github.com/egorlichenkoam/bgo3/pkg/card"
	"github.com/egorlichenkoam/bgo3/pkg/money"
	"github.com/egorlichenkoam/bgo3/pkg/person"
	"math/rand"
)

func GenerateTestData() (*person.Service, *card.Service, *Service, map[*person.Person]map[Mcc]money.Money, *person.Person) {
	personSvc := person.NewService()
	personSvc.Create("Иванов Иван Иванович")
	personSvc.Create("Петров Перт Петрович")
	personSvc.Create("Александров Александр Александрович")
	personsCount := len(personSvc.Persons)

	cardsNumbers := []string{
		"5106218416444735",
		"5106213218822113",
		"5106212866596714",
		"5106217691072252",
		"5106212352395522",
		"5106213096028379",
		"5106212135434895",
		"5106216399162894",
		"5106215378054189",
		"5106212023035804",
		"5106212615962522",
		"5106215392336513",
		"5106216651506119",
		"5106219357347762",
		"5106211376685587",
		"5106217418637700",
		"5106213096531406"}
	cardSvc := card.NewService("510621", "BABABANK")
	for _, number := range cardsNumbers {
		cardSvc.Create(personSvc.Persons[rand.Intn(personsCount)].Id, 10_000_000_00, card.Rub, number)
	}

	transactionSvc := NewService()
	transactions := make([]*Transaction, 100000)
	standard := map[*person.Person]map[Mcc]money.Money{}

	mccs := make([]Mcc, 0)
	for key := range MCCs() {
		mccs = append(mccs, key)
	}

	for range transactions {
		cardIdx := rand.Intn(len(cardSvc.Cards))
		mccIdx := rand.Intn(len(mccs))
		transactionSvc.CreateTransaction(100_00, mccs[mccIdx], cardSvc.Cards[cardIdx].Id, From)
	}

	for _, p := range personSvc.Persons {
		pCards := cardSvc.ByPersonId(p.Id)
		standard[p] = transactionSvc.SumByMCCs(transactionSvc.Transactions, pCards)
	}

	standardKeys := make([]*person.Person, 0)
	for key := range standard {
		standardKeys = append(standardKeys, key)
	}
	keyIdx := rand.Intn(len(standardKeys))
	p := standardKeys[keyIdx]

	return personSvc, cardSvc, transactionSvc, standard, p
}
