package transfer

import (
	"01/pkg/card"
	"01/pkg/money"
	"01/pkg/transaction"
	"errors"
)

var (
	errNotEnoughMoney    = errors.New("not enough money")
	errCardFromNotFound  = errors.New("card 'from' not found")
	errCardToNotFound    = errors.New("card 'to' not found")
	errCardNumberInvalid = errors.New("card number invalid")
)

type Commission struct {
	Percent float64
	Minimum money.Money
}

type Service struct {
	CardSvc              *card.Service
	TransactionSvc       *transaction.Service
	InBank               Commission
	ToDifferentBank      Commission
	BetweenDifferentBank Commission
}

func NewService(cardSvc *card.Service, transactionSvc *transaction.Service, inBank Commission, toDifferentBank Commission, betweenDifferentBank Commission) *Service {
	return &Service{
		CardSvc:              cardSvc,
		TransactionSvc:       transactionSvc,
		InBank:               inBank,
		ToDifferentBank:      toDifferentBank,
		BetweenDifferentBank: betweenDifferentBank,
	}
}

func (s *Service) Card2Card(from, to string, amount money.Money) (total money.Money, e error) {
	e = nil
	total = 0
	if !s.CardSvc.CheckByLuna(from) || !s.CardSvc.CheckByLuna(to) {
		e = errCardNumberInvalid
		return total, e
	}
	cardFrom := s.CardSvc.ByNumber(from)
	cardTo := s.CardSvc.ByNumber(to)
	total = s.total(amount, s.commission(cardFrom, cardTo))
	if cardFrom == nil {
		e = errCardFromNotFound
		return
	}
	if cardTo == nil {
		e = errCardToNotFound
		return
	}
	e = s.transfer(cardFrom, total, transaction.From)
	if e == nil {
		e = s.transfer(cardTo, amount, transaction.To)
	}
	return
}

func (s *Service) commission(cardFrom, cardTo *card.Card) *Commission {
	if cardFrom == nil && cardTo == nil {
		return &s.BetweenDifferentBank
	}
	if cardFrom != nil && cardTo == nil {
		return &s.ToDifferentBank
	}
	return &s.InBank
}

func (s *Service) total(amount money.Money, commission *Commission) money.Money {
	internalCommission := money.Money(float64(amount) / 100 * commission.Percent)
	if internalCommission < commission.Minimum {
		internalCommission = commission.Minimum
	}
	return amount + internalCommission
}

func (s *Service) transfer(card *card.Card, amount money.Money, fromTo transaction.Type) (e error) {
	e = nil
	tx := s.TransactionSvc.CreateTransaction(amount, "", card, fromTo)
	if fromTo == transaction.From {
		if card.Balance >= amount {
			card.Balance -= amount
			tx.Status = transaction.Ok
		} else {
			tx.Status = transaction.Fail
			e = errNotEnoughMoney
		}
	} else if fromTo == transaction.To {
		card.Balance += amount
		tx.Status = transaction.Ok
	}
	return
}
