package transaction

import (
	"github.com/egorlichenkoam/bgo3/pkg/card"
	"github.com/egorlichenkoam/bgo3/pkg/money"
	"github.com/egorlichenkoam/bgo3/pkg/person"
	"math/rand"
)

func GenerateTestData() (*person.Service, *card.Service, *Service, map[*person.Person]map[Mcc]money.Money, *person.Person) {
	personSvc := person.TestData()

	cardSvc := card.TestData(personSvc.Ids())

	txSvc := TestData(cardSvc.Ids())

	standard := map[*person.Person]map[Mcc]money.Money{}

	for _, p := range personSvc.Persons {
		pCards := cardSvc.ByPersonId(p.Id)
		standard[p] = txSvc.SumByMCCs(txSvc.Transactions, pCards)
	}

	standardKeys := make([]*person.Person, 0)
	for key := range standard {
		standardKeys = append(standardKeys, key)
	}
	keyIdx := rand.Intn(len(standardKeys))
	p := standardKeys[keyIdx]

	return personSvc, cardSvc, txSvc, standard, p
}
