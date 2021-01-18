package transaction

import (
	"01/pkg/card"
	"01/pkg/money"
	"01/pkg/person"
	"math/rand"
	"sort"
	"sync"
	"time"
)

type Type int

const (
	From Type = iota
	To
)

type Status string

const (
	Ok   Status = "Ok"
	Fail        = "Fail"
	Wait        = "Wait"
)

type Transaction struct {
	Id       int64
	Amount   money.Money
	Datetime int64
	Mcc      Mcc
	Status   Status
	Card     *card.Card
	Type     Type
}

type Service struct {
	Transactions []*Transaction
}

func NewService() *Service {
	return &Service{}
}

func (s *Service) CreateTransaction(amount money.Money, mcc Mcc, card *card.Card, fromTo Type) *Transaction {
	tx := Transaction{
		Id:       rand.Int63(),
		Amount:   amount,
		Datetime: time.Now().Unix(),
		Mcc:      mcc,
		Status:   Wait,
		Card:     card,
		Type:     fromTo,
	}
	s.Transactions = append(s.Transactions, &tx)
	return s.ById(tx.Id)
}

func (s *Service) ById(id int64) *Transaction {
	for i, tx := range s.Transactions {
		if tx.Id == id {
			return s.Transactions[i]
		}
	}
	return nil
}

func (s *Service) ByCard(card *card.Card) []*Transaction {
	result := make([]*Transaction, 0)
	for _, transaction := range s.Transactions {
		if transaction.Card.Id == card.Id {
			result = append(result, transaction)
		}
	}
	return result
}

func (s *Service) LastNTransactions(card *card.Card, n int) []*Transaction {
	transactions := s.ByCard(card)
	if len(transactions) < n {
		n = len(transactions)
	}
	n = len(transactions) - n
	transactions = transactions[n:]
	for i := len(transactions)/2 - 1; i >= 0; i-- {
		flipIdx := len(transactions) - 1 - i
		transactions[i], transactions[flipIdx] = transactions[flipIdx], transactions[i]
	}
	return transactions
}

func (s *Service) SumByMcc(card *card.Card, mccs []Mcc) money.Money {
	var result money.Money = 0
	transactions := filterTransactionsByMcc(s.ByCard(card), mccs)
	for _, transaction := range transactions {
		result = result + transaction.Amount
	}
	return result
}

func filterTransactionsByMcc(transactions []*Transaction, mccs []Mcc) []*Transaction {
	result := make([]*Transaction, 0)
	for _, transaction := range transactions {
		for _, mcc := range mccs {
			if transaction.Mcc == mcc {
				result = append(result, transaction)
			}
		}
	}
	return result
}

func (s *Service) TranslateMcc(code Mcc) string {
	result := "Категория не указана"
	value, ok := Mccs()[code]
	if ok {
		result = value
	}
	return result
}

func (s *Service) ByCardAndType(card *card.Card, fromTo Type) []*Transaction {
	result := make([]*Transaction, 0)
	cardTransactions := s.ByCard(card)
	//транзакции по типу (списание/зачисление)
	for n := range cardTransactions {
		tx := cardTransactions[n]
		if tx.Type == fromTo {
			result = append(result, tx)
		}
	}
	return result
}

func (s *Service) SortByCardAndType(card *card.Card, fromTo Type) []*Transaction {
	result := s.ByCardAndType(card, fromTo)
	sort.SliceStable(result, func(i, j int) bool {
		return result[i].Amount > result[j].Amount
	})
	return result
}

func makeYearMonthKey(unixTime int64) time.Time {
	t := time.Unix(unixTime, 0)
	return time.Date(t.Year(), t.Month(), 1, 0, 0, 0, 0, time.Local)
}

func (s *Service) GroupByCardAndYearMonth(card *card.Card, startTime, endTime int64, fromTo Type) map[time.Time][]*Transaction {
	if startTime < endTime {
		groupedTransactions := make(map[time.Time][]*Transaction, 0)
		next := time.Unix(startTime, 0)
		for next.Before(time.Unix(endTime, 0)) {
			groupedTransactions[makeYearMonthKey(next.Unix())] = make([]*Transaction, 0)
			next = next.AddDate(0, 1, 0)
		}
		groupedTransactions[makeYearMonthKey(endTime)] = make([]*Transaction, 0)
		transactions := s.ByCardAndType(card, fromTo)
		for n := range transactions {
			tx := transactions[n]
			mapKey := makeYearMonthKey(tx.Datetime)
			transactions, found := groupedTransactions[mapKey]
			if found {
				groupedTransactions[mapKey] = append(transactions, tx)
			}
		}
		return groupedTransactions
	}
	return nil
}

func (s *Service) SumConcurrentlyByCardAndYearMonth(card *card.Card, startTime, endTime int64, fromTo Type) (result sync.Map) {
	groupedTransactions := s.GroupByCardAndYearMonth(card, startTime, endTime, fromTo)
	count := len(groupedTransactions)
	wg := sync.WaitGroup{}
	wg.Add(count)
	for key, value := range groupedTransactions {
		go func(mark time.Time, transactions []*Transaction) {
			sum := money.Money(0)
			for i := range transactions {
				sum += transactions[i].Amount
			}
			result.Store(mark, sum)
			wg.Done()
		}(key, value)
	}
	wg.Wait()
	return result
}

func (s *Service) SumByPersonAndMccs(transactions []*Transaction, person *person.Person) (result map[Mcc]money.Money) {
	result = make(map[Mcc]money.Money)
	for _, tx := range transactions {
		for _, c := range person.Cards {
			if tx.Card == c {
				result[tx.Mcc] += tx.Amount
			}
		}
	}
	return result
}

func (s *Service) SumByPersonAndMccsWithMutex(transactions []*Transaction, person *person.Person) map[Mcc]money.Money {
	partCount := 10
	wg := sync.WaitGroup{}
	wg.Add(partCount)
	mu := sync.Mutex{}
	result := make(map[Mcc]money.Money)
	partSize := len(transactions) / partCount
	for i := 0; i < partCount; i++ {
		part := transactions[i*partSize : (i+1)*partSize]
		if i == partCount-1 {
			for _, value := range transactions[(i+1)*partSize:] {
				part = append(part, value)
			}
		}
		go func() {
			m := s.SumByPersonAndMccs(part, person)
			mu.Lock()
			for key, value := range m {
				result[key] += value
			}
			mu.Unlock()
			wg.Done()
		}()
	}
	wg.Wait()
	return result
}
