package transaction

import (
	"01/pkg/card"
	"01/pkg/money"
	"01/pkg/person"
	"bytes"
	"encoding/csv"
	"io"
	"io/ioutil"
	"math/rand"
	"os"
	"sort"
	"strconv"
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
	CardId   int64
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
		CardId:   card.Id,
		Type:     fromTo,
	}
	s.Transactions = append(s.Transactions, &tx)
	return &tx
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
		if transaction.CardId == card.Id {
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
			if tx.CardId == c.Id {
				result[tx.Mcc] += tx.Amount
				break
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

func (s *Service) SumByPersonAndMccsWithChannels(transactions []*Transaction, person *person.Person) map[Mcc]money.Money {
	partCount := 10
	result := make(map[Mcc]money.Money)
	chMap := make(chan map[Mcc]money.Money)
	partSize := len(transactions) / partCount
	for i := 0; i < partCount; i++ {
		part := transactions[i*partSize : (i+1)*partSize]
		if i == partCount-1 {
			for _, value := range transactions[(i+1)*partSize:] {
				part = append(part, value)
			}
		}
		go func(chMap chan<- map[Mcc]money.Money) {
			chMap <- s.SumByPersonAndMccs(part, person)
		}(chMap)
	}
	finished := 0
	for value := range chMap {
		for key, value := range value {
			result[key] += value
		}
		finished++
		if finished == partCount {
			break
		}
	}
	return result
}

func (s *Service) SumByPersonAndMccsWithMutexStraightToMap(transactions []*Transaction, person *person.Person) map[Mcc]money.Money {
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
			for _, tx := range part {
				for _, c := range person.Cards {
					if tx.CardId == c.Id {
						mu.Lock()
						result[tx.Mcc] += tx.Amount
						mu.Unlock()
						break
					}
				}
			}
			wg.Done()
		}()
	}
	wg.Wait()
	return result
}

func (t Transaction) strings() (result []string) {
	result = make([]string, 0)
	result = append(result, strconv.Itoa(int(t.Id)))
	result = append(result, strconv.Itoa(int(t.Amount)))
	result = append(result, strconv.Itoa(int(t.Datetime)))
	result = append(result, string(t.Mcc))
	result = append(result, string(t.Status))
	result = append(result, strconv.Itoa(int(t.CardId)))
	result = append(result, strconv.Itoa(int(t.Type)))
	return result
}

func (t *Transaction) mapRowToTransaction(content []string) (err error) {
	for key, value := range content {
		switch key {
		case 0:
			t.Id, err = strconv.ParseInt(value, 10, 64)
			if err != nil {
				return err
			}
			break
		case 1:
			var amount int64 = 0
			amount, err = strconv.ParseInt(value, 10, 64)
			if err != nil {
				return err
			}
			t.Amount = money.Money(amount)
			break
		case 2:
			t.Datetime, err = strconv.ParseInt(value, 10, 64)
			if err != nil {
				return err
			}
			break
		case 3:
			t.Mcc = Mcc(value)
			break
		case 4:
			t.Status = Status(value)
			break
		case 5:
			t.CardId, err = strconv.ParseInt(value, 10, 64)
			if err != nil {
				return err
			}
			break
		case 6:
			var txType int
			txType, err = strconv.Atoi(value)
			if err != nil {
				return err
			}
			t.Type = Type(txType)
			break
		}
	}
	return nil
}

func ExportCsv(transactions []*Transaction) (err error) {
	file, err := os.Create("exports.csv")
	if err != nil {
		return err
	}
	defer func(c io.Closer) {
		if cerr := c.Close(); cerr != nil {
			err = cerr
		}
	}(file)
	writer := csv.NewWriter(file)
	defer writer.Flush()
	for _, tx := range transactions {
		err = writer.Write(tx.strings())
		if err != nil {
			return err
		}
	}
	return nil
}

func ImportCsv(filePath string) ([]*Transaction, error) {
	transactions := make([]*Transaction, 0)
	data, err := ioutil.ReadFile(filePath)
	if err != nil {
		return nil, err
	}
	reader := csv.NewReader(bytes.NewReader(data))
	records, err := reader.ReadAll()
	if err != nil {
		return nil, err
	}
	for _, content := range records {
		tx := Transaction{}
		err = tx.mapRowToTransaction(content)
		if err != nil {
			return nil, err
		}
		transactions = append(transactions, &tx)
	}
	return transactions, nil
}
