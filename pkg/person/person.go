package person

import (
	"encoding/json"
	"io"
	"log"
	"math/rand"
	"os"
	"sync"
)

const name = "person"

type Person struct {
	Id   int64  `json:"id"`
	Name string `json:"name"`
}

type Service struct {
	name    string
	mu      sync.RWMutex
	Persons []*Person `json:"person"`
}

func NewService() *Service {
	return &Service{name: "Person service",
		mu: sync.RWMutex{}}
}

func (s *Service) Create(name string) *Person {
	return s.add(&Person{
		Id:   rand.Int63(),
		Name: name,
	})
}

func (s *Service) add(person *Person) *Person {
	log.Printf("%s - %s", s.name, "Start add")
	s.mu.Lock()
	defer s.mu.Unlock()
	s.Persons = append(s.Persons, person)
	log.Printf("%s - %s", s.name, "End add")
	return person
}

func (s *Service) Ids() []int64 {
	log.Printf("%s -%s", s.name, "Start Ids")
	ids := make([]int64, len(s.Persons))
	for _, p := range s.Persons {
		ids = append(ids, p.Id)
	}
	log.Printf("%s -%s", s.name, "End Ids")
	return ids
}

//check if person exist in service
func (s *Service) Exist(personId int64) bool {
	log.Printf("%s - %s : %d", s.name, "Start Exist", personId)
	s.mu.RLock()
	defer s.mu.RUnlock()
	result := false
	for _, id := range s.Ids() {
		if id == personId {
			result = true
			break
		}
	}
	log.Printf("%s - %s : %d", s.name, "End Exist", personId)
	return result
}

func ExportJson(persons []*Person) error {
	log.Printf("%s - %s", name, "Start export")
	file, err := os.Create("personsExport.json")
	if err != nil {
		return err
	}
	defer func(c io.Closer) {
		if cErr := c.Close(); cErr != nil {
			err = cErr
		}
	}(file)
	encoder := json.NewEncoder(file)
	err = encoder.Encode(persons)
	if err != nil {
		return err
	}
	log.Printf("%s - %s", name, "End export")
	return nil
}

func ImportJson(filePath string) ([]*Person, error) {
	log.Printf("%s - %s", name, "Start import")
	persons := make([]*Person, 0)
	reader, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer func(c io.Closer) {
		if cErr := c.Close(); cErr != nil {
			err = cErr
		}
	}(reader)
	err = json.NewDecoder(reader).Decode(&persons)
	if err != nil {
		return nil, err
	}
	log.Printf("%s - %s", name, "End import")
	return persons, nil
}

func TestData() *Service {
	log.Printf("%s - %s", name, "Start test data")
	personSvc := NewService()
	personSvc.Create("Иванов Иван Иванович")
	personSvc.Create("Петров Перт Петрович")
	personSvc.Create("Александров Александр Александрович")
	log.Printf("%s - %s", name, "End test data")
	return personSvc
}
