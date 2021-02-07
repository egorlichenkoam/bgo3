package card

import (
	"github.com/egorlichenkoam/bgo3/pkg/money"
	"math/rand"
	"strconv"
	"strings"
)

type Currency string

const (
	Rub Currency = "RUB"
)

type Card struct {
	Id       int64
	Issuer   string
	Balance  money.Money
	Currency Currency
	Number   string
	Icon     string
}

func (s *Service) Create(balance money.Money, currency Currency, number string) *Card {
	return s.Add(Card{Id: rand.Int63(), Issuer: s.Issuer, Balance: balance, Currency: currency, Number: number, Icon: ""})
}

type Service struct {
	IssuerId string
	Issuer   string
	Cards    []Card
}

func NewService(issuerId, issuer string) *Service {
	return &Service{
		IssuerId: issuerId,
		Issuer:   issuer,
	}
}

func (s *Service) Add(card Card) *Card {
	s.Cards = append(s.Cards, card)
	return &card
}

func (s *Service) ByNumber(number string) (card *Card) {
	card = nil
	if s.isOurCard(number) {
		for i, c := range s.Cards {
			if c.Number == number {
				card = &s.Cards[i]
			}
		}
		if card == nil {
			card = s.Create(0, Rub, number)
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
