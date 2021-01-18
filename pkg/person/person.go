package person

import (
	"01/pkg/card"
	"math/rand"
)

type Person struct {
	Id    int64
	Name  string
	Cards []*card.Card
}

type Service struct {
	Persons []*Person
}

func NewService() *Service {
	return &Service{}
}

func (s *Service) Create(name string) *Person {
	return s.add(&Person{
		Id:    rand.Int63(),
		Name:  name,
		Cards: make([]*card.Card, 0),
	})
}

func (s *Service) add(person *Person) *Person {
	s.Persons = append(s.Persons, person)
	return person
}

func (s *Service) AddCard(person *Person, card *card.Card) {
	person.Cards = append(person.Cards, card)
}
