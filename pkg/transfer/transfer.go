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
	PercentInBank       float64
	MinimumInBank       money.Money
	PercentToDiffBank   float64
	MinimumToDiffBank   money.Money
	PercentBetweenBanks float64
	MinimumBetweenBanks money.Money
}

type Service struct {
	CardSvc        *card.Service
	TransactionSvc *transaction.Service
	commissions    Commission
}

func NewService(cardSvc *card.Service, transactionSvc *transaction.Service, commissions Commission) *Service {
	return &Service{
		CardSvc:        cardSvc,
		TransactionSvc: transactionSvc,
		commissions:    commissions,
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
	percent, minimum := s.commission(cardFrom, cardTo)
	total = s.total(amount, percent, minimum)
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

func (s *Service) commission(cardFrom, cardTo *card.Card) (percent float64, minimum money.Money) {
	if cardFrom == nil && cardTo == nil {
		return s.commissions.PercentBetweenBanks, s.commissions.MinimumBetweenBanks
	}
	if cardFrom != nil && cardTo == nil {
		return s.commissions.PercentToDiffBank, s.commissions.MinimumToDiffBank
	}
	return s.commissions.PercentInBank, s.commissions.MinimumInBank
}

func (s *Service) total(amount money.Money, percent float64, minimum money.Money) money.Money {
	internalCommission := money.Money(float64(amount) / 100 * percent)
	if internalCommission < minimum {
		internalCommission = minimum
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
