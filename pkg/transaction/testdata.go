package transaction

import (
	"github.com/egorlichenkoam/bgo3/pkg/card"
	"github.com/egorlichenkoam/bgo3/pkg/money"
	"github.com/egorlichenkoam/bgo3/pkg/person"
	"math/rand"
)

func GenerateTestData() ([]*Transaction, map[*person.Person]map[Mcc]money.Money, *person.Person) {
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
		c := cardSvc.Create(10_000_000_00, card.Rub, number)
		personSvc.AddCard(personSvc.Persons[rand.Intn(personsCount)], c)
	}

	transactionSvc := NewService()
	transactions := make([]*Transaction, 1000)
	standard := map[*person.Person]map[Mcc]money.Money{}

	mccs := make([]Mcc, 0)
	for key := range Mccs() {
		mccs = append(mccs, key)
	}

	for i := range transactions {
		pers := personSvc.Persons[rand.Intn(personsCount)]
		cardIdx := rand.Intn(len(pers.Cards))
		standardMap := standard[pers]
		if standardMap == nil {
			standardMap = map[Mcc]money.Money{}
		}
		mccIdx := rand.Intn(len(mccs))
		tx := transactionSvc.CreateTransaction(100_00, mccs[mccIdx], pers.Cards[cardIdx], From)
		transactions[i] = tx
		standardMap[tx.Mcc] += tx.Amount
		standard[pers] = standardMap
	}

	standardKeys := make([]*person.Person, 0)
	for key := range standard {
		standardKeys = append(standardKeys, key)
	}
	keyIdx := rand.Intn(len(standardKeys))
	pers := standardKeys[keyIdx]

	return transactions, standard, pers
}
