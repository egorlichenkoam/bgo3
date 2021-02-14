package card

import (
	"encoding/json"
	"github.com/egorlichenkoam/bgo3/pkg/money"
	"io"
	"log"
	"math/rand"
	"os"
	"strconv"
	"strings"
	"sync"
)

type Currency string

const (
	Rub  Currency = "RUB"
	name string   = "Card"
)

type Card struct {
	Id       int64       `json:"id"`
	Issuer   string      `json:"issuer"`
	Balance  money.Money `json:"balance"`
	Currency Currency    `json:"currency"`
	Number   string      `json:"number"`
	Icon     string      `json:"icon"`
	Type     string      `json:"type"`
	PersonId int64       `json:"personId"`
}

type CardDTO struct {
	Issuer string `json:"issuer"`
	Number string `json:"number"`
	Type   string `json:"type"`
}
type Service struct {
	name     string
	mu       sync.RWMutex
	IssuerId string
	Issuer   string
	Cards    []*Card
}

func (c Card) DTO() CardDTO {
	log.Printf("%s - %s", name, "Start DTO")
	dto := CardDTO{
		Issuer: c.Issuer,
		Number: c.Number,
		Type:   c.Type,
	}
	log.Printf("%s - %s", name, "End DTO")
	return dto
}

func NewService(issuerId, issuer string) *Service {
	return &Service{
		name:     "Card service",
		mu:       sync.RWMutex{},
		IssuerId: issuerId,
		Issuer:   issuer,
	}
}

func (s *Service) Create(issuer string, personId int64, balance money.Money, currency Currency, number string, cardType string) *Card {
	log.Printf("%s - %s", s.name, "Start Create")
	s.mu.Lock()
	defer s.mu.Unlock()
	c := Card{Id: rand.Int63(), Issuer: issuer, Balance: balance, Currency: currency, Number: number, Icon: "", PersonId: personId, Type: cardType}
	log.Printf("%s - card added : %s", s.name, c.Number)
	log.Printf("%s - %s", s.name, "End Create")
	return s.Add(&c)
}

func (s *Service) Add(card *Card) *Card {
	s.Cards = append(s.Cards, card)
	return card
}

//return slice of cards by person id
func (s *Service) ByPersonId(personId int64) []*Card {
	personCards := make([]*Card, 0)
	for _, card := range s.Cards {
		if card.PersonId == personId {
			personCards = append(personCards, card)
		}
	}
	return personCards
}

func (s *Service) ByNumber(number string) (card *Card) {
	card = nil
	if s.isOurCard(number) {
		for i, c := range s.Cards {
			if c.Number == number {
				card = s.Cards[i]
			}
		}
		if card == nil {
			card = s.Create("VISA", 0, 0, Rub, number, "VIRTUAL")
		}
	}
	return
}

func (s *Service) isOurCard(number string) bool {
	if strings.HasPrefix(number, s.IssuerId) {
		return true
	}
	return false
}

func (s *Service) CheckByLuna(number string) bool {
	numberInString := strings.Split(strings.ReplaceAll(number, " ", ""), "")
	sum := 0
	for idx := range numberInString {
		if sn, e := strconv.Atoi(numberInString[idx]); e == nil {
			if (idx+1)%2 > 0 {
				sn = sn * 2
				if sn > 9 {
					sn = sn - 9
				}
			}
			sum += sn
		} else {
			return false
		}
	}
	return sum%10 == 0
}

func (s *Service) Ids() []int64 {
	log.Printf("%s -%s", s.name, "Start Ids")
	ids := make([]int64, len(s.Cards))
	for _, c := range s.Cards {
		ids = append(ids, c.Id)
	}
	log.Printf("%s -%s", s.name, "End Ids")
	return ids
}

func ExportJson(cards []*Card) error {
	log.Printf("%s - %s", name, "Start export")
	file, err := os.Create("cardsExport.json")
	if err != nil {
		return err
	}
	defer func(c io.Closer) {
		if cErr := c.Close(); cErr != nil {
			err = cErr
		}
	}(file)
	encoder := json.NewEncoder(file)
	err = encoder.Encode(cards)
	if err != nil {
		return err
	}
	log.Printf("%s - %s", name, "End export")
	return nil
}

func ImportJson(filePath string) (cards []*Card, err error) {
	log.Printf("%s - %s", name, "Start import")
	cards = make([]*Card, 0)
	reader, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer func(c io.Closer) {
		if cErr := c.Close(); cErr != nil {
			err = cErr
		}
	}(reader)
	err = json.NewDecoder(reader).Decode(&cards)
	if err != nil {
		return nil, err
	}
	log.Printf("%s - %s", name, "End import")
	return cards, nil
}

func TestData(personsIds []int64) *Service {
	log.Printf("%s - %s", name, "Start test data")
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
	cardSvc := NewService("510621", "VISA")
	for _, number := range cardsNumbers {
		personId := int64(0)
		if personsIds != nil {
			if len(personsIds) > 0 {
				personId = personsIds[rand.Intn(len(personsIds))]
			}
		}
		cardSvc.Create("VISA", personId, 10_000_000_00, Rub, number, "PLASTIC")
	}
	log.Printf("%s - %s", name, "End test data")
	return cardSvc
}
