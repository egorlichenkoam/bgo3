package person

import (
	"math/rand"
)

type Person struct {
	Id   int64  `json:"id"`
	Name string `json:"name"`
}

type Service struct {
	Persons []*Person `json:"person"`
}

func NewService() *Service {
	return &Service{}
}

func (s *Service) Create(name string) *Person {
	return s.add(&Person{
		Id:   rand.Int63(),
		Name: name,
	})
}

func (s *Service) add(person *Person) *Person {
	s.Persons = append(s.Persons, person)
	return person
}
