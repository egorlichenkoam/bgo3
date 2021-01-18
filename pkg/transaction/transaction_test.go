package transaction

import (
	"01/pkg/card"
	"01/pkg/money"
	"01/pkg/person"
	"math/rand"
	"reflect"
	"testing"
)

func testData() ([]*Transaction, map[*person.Person]map[Mcc]money.Money, *person.Person) {
	personSvc := person.NewService()
	personSvc.Create("Иванов Иван Иванович")
	personSvc.Create("Петров Перт Петрович")
	personSvc.Create("Александров Александр Александрович")
	personsCount := len(personSvc.Persons)

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
	cardSvc := card.NewService("510621", "BABABANK")
	for _, number := range cardsNumbers {
		c := cardSvc.Create(10_000_000_00, card.Rub, number)
		personSvc.AddCard(personSvc.Persons[rand.Intn(personsCount)], c)
	}

	transactionSvc := NewService()
	transactions := make([]*Transaction, 100000)
	standard := map[*person.Person]map[Mcc]money.Money{}

	mccs := make([]Mcc, 0)
	for key := range Mccs() {
		mccs = append(mccs, key)
	}

	for i := range transactions {
		pers := personSvc.Persons[rand.Intn(personsCount)]
		cardIdx := rand.Intn(len(pers.Cards))
		standardMap := standard[pers]
		if standardMap == nil {
			standardMap = map[Mcc]money.Money{}
		}
		mccIdx := rand.Intn(len(mccs))
		tx := transactionSvc.CreateTransaction(100_00, mccs[mccIdx], pers.Cards[cardIdx], From)
		transactions[i] = tx
		standardMap[tx.Mcc] += tx.Amount
		standard[pers] = standardMap
	}

	standardKeys := make([]*person.Person, 0)
	for key := range standard {
		standardKeys = append(standardKeys, key)
	}
	keyIdx := rand.Intn(len(standardKeys))
	pers := standardKeys[keyIdx]

	return transactions, standard, pers
}

func TestService_SortedByType(t *testing.T) {
	cardSvc := card.NewService("510621", "BABANK")
	transactionSvc := NewService()
	personSvc := person.NewService()

	pers := personSvc.Create("Иванов Иван Иванович")
	card00 := cardSvc.Create(1000_000_00, card.Rub, "5106212879499054")
	personSvc.AddCard(pers, card00)

	transactionSvc.CreateTransaction(1_000_00, "", card00, From)
	transactionSvc.CreateTransaction(5_000_00, "", card00, From)
	transactionSvc.CreateTransaction(6_000_00, "", card00, From)
	transactionSvc.CreateTransaction(500_00, "", card00, From)
	transactionSvc.CreateTransaction(50_000_00, "", card00, From)
	transactionSvc.CreateTransaction(49_000_00, "", card00, From)

	transactions := []Transaction{
		{
			Amount: 50_000_00,
			Card:   card00,
			Type:   From,
		},
		{
			Amount: 49_000_00,
			Card:   card00,
			Type:   From,
		},
		{
			Amount: 6_000_00,
			Card:   card00,
			Type:   From,
		},
		{
			Amount: 5_000_00,
			Card:   card00,
			Type:   From,
		},
		{
			Amount: 1_000_00,
			Card:   card00,
			Type:   From,
		},
		{
			Amount: 500_00,
			Card:   card00,
			Type:   From,
		},
	}

	type fields struct {
		TransactionSvc *Service
	}
	type args struct {
		card *card.Card
	}

	tests := []struct {
		name   string
		fields fields
		args   args
		want   []Transaction
	}{
		{
			name: "Сортировка транзакций",
			fields: fields{
				TransactionSvc: transactionSvc,
			},
			args: args{
				card: card00,
			},
			want: transactions,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.fields.TransactionSvc.SortByCardAndType(tt.args.card, From); !areTransactionsEquals(got, tt.want) {
				t.Errorf("\n got  = %v,\n want = %v", got, tt.want)
			}
		})
	}
}

//Сторонняя функция проверки используя потому, что транзакциям при переводе
//выставляется время и идентификатор автоматически и повторить их в тестовых данных
//для сравнения будет проблематично
func areTransactionsEquals(got []*Transaction, want []Transaction) bool {
	if len(got) != len(want) {
		return false
	}
	for n := range want {
		gotTx := got[n]
		wantTx := want[n]
		if (gotTx.Card.Number != wantTx.Card.Number) && (gotTx.Amount != wantTx.Amount) {
			return false
		}
	}
	return true
}

func TestService_SumByPersonAndMccs(t *testing.T) {
	transactions, standard, pers := testData()

	type fields struct {
		Transactions []*Transaction
	}
	type args struct {
		transactions []*Transaction
		person       *person.Person
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   map[Mcc]money.Money
	}{
		{
			name: "Вывод группированных по MCC затрат",
			fields: fields{
				Transactions: transactions,
			},
			args: args{
				transactions: transactions,
				person:       pers,
			},
			want: standard[pers],
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &Service{
				Transactions: tt.fields.Transactions,
			}
			if gotResult := s.SumByPersonAndMccs(tt.args.transactions, tt.args.person); !reflect.DeepEqual(gotResult, tt.want) {
				t.Errorf("SumByPersonAndMccs() = %v, want %v", gotResult, tt.want)
			}
		})
	}
}

func TestService_SumByPersonAndMccsWithMutex(t *testing.T) {
	transactions, standard, pers := testData()

	type fields struct {
		Transactions []*Transaction
	}
	type args struct {
		transactions []*Transaction
		person       *person.Person
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   map[Mcc]money.Money
	}{
		{
			name: "Вывод группированных по MCC затрат с mutex",
			fields: fields{
				Transactions: transactions,
			},
			args: args{
				transactions: transactions,
				person:       pers,
			},
			want: standard[pers],
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &Service{
				Transactions: tt.fields.Transactions,
			}
			if got := s.SumByPersonAndMccsWithMutex(tt.args.transactions, tt.args.person); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("SumByPersonAndMccWithMutex() = %v, want %v", got, tt.want)
			}
		})
	}
}
