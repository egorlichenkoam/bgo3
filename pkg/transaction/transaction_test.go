package transaction

import (
	"01/pkg/card"
	"testing"
)

func TestService_SortedByType(t *testing.T) {
	cardSvc := card.NewService("510621")
	transactionSvc := NewService()
	card00 := cardSvc.NewCard("BABANK", 1000_000_00, card.Rub, "5106212879499054")

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
			if got := transactionSvc.SortedByType(tt.args.card, From); !areTransactionsEquals(got, tt.want) {
				t.Errorf("\n got  = %v,\n want = %v", got, tt.want)
			}
		})
	}
}

//Сторонняя функция проверки используя потому, что транзакциям при переводе
//выставляется время и идентификатор автоматически и повторить их в тестовых данных
//для сравнения будет проблематично
func areTransactionsEquals(got []Transaction, want []Transaction) bool {
	if len(got) != len(want) {
		return false
	}
	for n := range want {
		gotTx := got[n]
		wantTx := want[n]
		if gotTx.Card.Number != wantTx.Card.Number {
			return false
		}
		if gotTx.Amount != wantTx.Amount {
			return false
		}
	}
	return true
}
