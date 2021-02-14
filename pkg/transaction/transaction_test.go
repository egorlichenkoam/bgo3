package transaction

import (
	"github.com/egorlichenkoam/bgo3/pkg/card"
	"github.com/egorlichenkoam/bgo3/pkg/money"
	"github.com/egorlichenkoam/bgo3/pkg/person"
	"os"
	"reflect"
	"testing"
)

var GPersonSvc *person.Service = nil
var GCardSvc *card.Service = nil
var GTransactionSvc *Service = nil
var GStandard map[*person.Person]map[Mcc]money.Money = nil
var GPerson *person.Person = nil

func CreateTestData() {
	if (GPersonSvc == nil) || (GCardSvc == nil) || (GTransactionSvc == nil) || (GStandard == nil) || (GPerson == nil) {
		GPersonSvc, GCardSvc, GTransactionSvc, GStandard, GPerson = GenerateTestData()
	}
}

func TestService_SortedByType(t *testing.T) {
	cardSvc := card.NewService("510621", "VISA")
	transactionSvc := NewService()
	personSvc := person.NewService()

	p := personSvc.Create("Иванов Иван Иванович")
	card00 := cardSvc.Create(p.Id, 1000_000_00, card.Rub, "5106212879499054")

	transactionSvc.CreateTransaction(1_000_00, "", card00.Id, From)
	transactionSvc.CreateTransaction(5_000_00, "", card00.Id, From)
	transactionSvc.CreateTransaction(6_000_00, "", card00.Id, From)
	transactionSvc.CreateTransaction(500_00, "", card00.Id, From)
	transactionSvc.CreateTransaction(50_000_00, "", card00.Id, From)
	transactionSvc.CreateTransaction(49_000_00, "", card00.Id, From)

	transactions := []Transaction{
		{
			Amount: 50_000_00,
			CardId: card00.Id,
			Type:   From,
		},
		{
			Amount: 49_000_00,
			CardId: card00.Id,
			Type:   From,
		},
		{
			Amount: 6_000_00,
			CardId: card00.Id,
			Type:   From,
		},
		{
			Amount: 5_000_00,
			CardId: card00.Id,
			Type:   From,
		},
		{
			Amount: 1_000_00,
			CardId: card00.Id,
			Type:   From,
		},
		{
			Amount: 500_00,
			CardId: card00.Id,
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
		if (gotTx.CardId != wantTx.CardId) && (gotTx.Amount != wantTx.Amount) {
			return false
		}
	}
	return true
}

func TestService_SumByMCCs(t *testing.T) {
	CreateTestData()

	type fields struct {
		Transactions []*Transaction
	}
	type args struct {
		transactions []*Transaction
		cards        []*card.Card
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   map[Mcc]money.Money
	}{
		{
			name: "TestService_SumByMCCs",
			fields: fields{
				Transactions: GTransactionSvc.Transactions,
			},
			args: args{
				transactions: GTransactionSvc.Transactions,
				cards:        GCardSvc.ByPersonId(GPerson.Id),
			},
			want: GStandard[GPerson],
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &Service{
				Transactions: tt.fields.Transactions,
			}
			if gotResult := s.SumByMCCs(tt.args.transactions, tt.args.cards); !reflect.DeepEqual(gotResult, tt.want) {
				t.Errorf("SumByMCCs() = %v, want %v", gotResult, tt.want)
			}
		})
	}
}

func BenchmarkSumByMCCs(b *testing.B) {
	CreateTestData()

	s := NewService()
	want := GStandard[GPerson]
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		result := s.SumByMCCs(GTransactionSvc.Transactions, GCardSvc.ByPersonId(GPerson.Id))
		b.StopTimer()
		if !reflect.DeepEqual(result, want) {
			b.Fatalf("invalid result, got %v, want %v", result, want)
		}
		b.StartTimer()
	}
}

func TestService_SumByMCCsWithMutex(t *testing.T) {
	CreateTestData()

	type fields struct {
		Transactions []*Transaction
	}
	type args struct {
		transactions []*Transaction
		cards        []*card.Card
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   map[Mcc]money.Money
	}{
		{
			name: "TestService_SumByMCCsWithMutex",
			fields: fields{
				Transactions: GTransactionSvc.Transactions,
			},
			args: args{
				transactions: GTransactionSvc.Transactions,
				cards:        GCardSvc.ByPersonId(GPerson.Id),
			},
			want: GStandard[GPerson],
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &Service{
				Transactions: tt.fields.Transactions,
			}
			if got := s.SumByMCCsWithMutex(tt.args.transactions, tt.args.cards); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("SumByMccWithMutex() = %v, want %v", got, tt.want)
			}
		})
	}
}

func BenchmarkSumByMCCsWithMutex(b *testing.B) {
	CreateTestData()

	s := NewService()
	want := GStandard[GPerson]
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		result := s.SumByMCCsWithMutex(GTransactionSvc.Transactions, GCardSvc.ByPersonId(GPerson.Id))
		b.StopTimer()
		if !reflect.DeepEqual(result, want) {
			b.Fatalf("invalid result, got %v, want %v", result, want)
		}
		b.StartTimer()
	}
}

func TestService_SumByMCCsWithChannels(t *testing.T) {
	CreateTestData()

	type fields struct {
		Transactions []*Transaction
	}
	type args struct {
		transactions []*Transaction
		cards        []*card.Card
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   map[Mcc]money.Money
	}{
		{
			name: "TestService_SumByMCCsWithChannels",
			fields: fields{
				Transactions: GTransactionSvc.Transactions,
			},
			args: args{
				transactions: GTransactionSvc.Transactions,
				cards:        GCardSvc.ByPersonId(GPerson.Id),
			},
			want: GStandard[GPerson],
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &Service{
				Transactions: tt.fields.Transactions,
			}
			if got := s.SumByMCCsWithChannels(tt.args.transactions, tt.args.cards); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("SumByMCCsWithChannels() = %v, want %v", got, tt.want)
			}
		})
	}
}

func BenchmarkSumByMCCsWithChannels(b *testing.B) {
	CreateTestData()

	s := NewService()
	want := GStandard[GPerson]
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		result := s.SumByMCCsWithChannels(GTransactionSvc.Transactions, GCardSvc.ByPersonId(GPerson.Id))
		b.StopTimer()
		if !reflect.DeepEqual(result, want) {
			b.Fatalf("invalid result, got %v, want %v", result, want)
		}
		b.StartTimer()
	}
}

func TestService_SumByMCCsWithMutexStraightToMap(t *testing.T) {
	CreateTestData()

	type fields struct {
		Transactions []*Transaction
	}
	type args struct {
		transactions []*Transaction
		cards        []*card.Card
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   map[Mcc]money.Money
	}{
		{
			name: "TestService_SumByMCCsWithMutexStraightToMap",
			fields: fields{
				Transactions: GTransactionSvc.Transactions,
			},
			args: args{
				transactions: GTransactionSvc.Transactions,
				cards:        GCardSvc.ByPersonId(GPerson.Id),
			},
			want: GStandard[GPerson],
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &Service{
				Transactions: tt.fields.Transactions,
			}
			if got := s.SumByMCCsWithMutexStraightToMap(tt.args.transactions, tt.args.cards); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("SumByMCCsWithMutexStraightToMap() = %v, want %v", got, tt.want)
			}
		})
	}
}

func BenchmarkSumByMCCsWithMutexStraightToMap(b *testing.B) {
	CreateTestData()

	s := NewService()
	want := GStandard[GPerson]
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		result := s.SumByMCCsWithMutexStraightToMap(GTransactionSvc.Transactions, GCardSvc.ByPersonId(GPerson.Id))
		b.StopTimer()
		if !reflect.DeepEqual(result, want) {
			b.Fatalf("invalid result, got %v, want %v", result, want)
		}
		b.StartTimer()
	}
}

func TestExportCsv(t *testing.T) {
	CreateTestData()
	type args struct {
		transactions []*Transaction
	}
	tests := []struct {
		name string
		args args
		want error
	}{
		{
			name: "TestExportCsv",
			args: args{
				transactions: GTransactionSvc.Transactions,
			},
			want: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := ExportCsv(tt.args.transactions); err != tt.want {
				t.Errorf("ExportCsv() error = %v, wantErr %v", err, tt.want)
			}
		})
	}
}

func TestImportCsv(t *testing.T) {
	CreateTestData()
	fPath, _ := os.Getwd()
	fPath = fPath + "/exports.csv"
	type args struct {
		filePath string
	}
	tests := []struct {
		name    string
		args    args
		want    []*Transaction
		wantErr error
	}{
		{
			name: "TestImportCsv",
			args: args{
				filePath: fPath,
			},
			want:    GTransactionSvc.Transactions,
			wantErr: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ImportCsv(tt.args.filePath)
			if err != tt.wantErr {
				t.Errorf("ImportCsv() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ImportCsv() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestExportJson(t *testing.T) {
	CreateTestData()
	type args struct {
		transactions []*Transaction
	}
	tests := []struct {
		name    string
		args    args
		wantErr error
	}{
		{
			name: "TestExportJson",
			args: args{
				transactions: GTransactionSvc.Transactions,
			},
			wantErr: nil,
		}}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := ExportJson(tt.args.transactions); err != tt.wantErr {
				t.Errorf("ExportJson() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestImportJson(t *testing.T) {
	CreateTestData()
	fPath, _ := os.Getwd()
	fPath = fPath + "/exports.json"
	type args struct {
		filePath string
	}
	tests := []struct {
		name    string
		args    args
		want    []*Transaction
		wantErr error
	}{
		{
			name: "TestImportJson",
			args: args{
				filePath: fPath,
			},
			want:    GTransactionSvc.Transactions,
			wantErr: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ImportJson(tt.args.filePath)
			if err != tt.wantErr {
				t.Errorf("ImportJson() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ImportJson() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestExportXml(t *testing.T) {
	CreateTestData()
	type args struct {
		transactions []*Transaction
	}
	tests := []struct {
		name    string
		args    args
		wantErr error
	}{
		{
			name: "TestExportXml",
			args: args{
				transactions: GTransactionSvc.Transactions,
			},
			wantErr: nil,
		}}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := ExportXml(tt.args.transactions); err != tt.wantErr {
				t.Errorf("ExportXml() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestImportXml(t *testing.T) {
	CreateTestData()
	fPath, _ := os.Getwd()
	fPath = fPath + "/exports.xml"
	type args struct {
		filePath string
	}
	tests := []struct {
		name    string
		args    args
		want    []*Transaction
		wantErr error
	}{
		{
			name: "TestImportXml",
			args: args{
				filePath: fPath,
			},
			want:    GTransactionSvc.Transactions,
			wantErr: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ImportXml(tt.args.filePath)
			if err != tt.wantErr {
				t.Errorf("ImportXml() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ImportXml() got = %v, want %v", got, tt.want)
			}
		})
	}
}
