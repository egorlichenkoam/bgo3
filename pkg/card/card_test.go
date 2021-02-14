package card

import (
	"os"
	"reflect"
	"testing"
)

var gCardSvc *Service

func TestExportJson(t *testing.T) {
	gCardSvc = TestData(nil)
	type args struct {
		cards []*Card
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{{
		name: name + "-TestExportJson",
		args: args{
			cards: gCardSvc.Cards,
		},
		wantErr: false,
	}}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := ExportJson(tt.args.cards); (err != nil) != tt.wantErr {
				t.Errorf("ExportJson() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestImportJson(t *testing.T) {
	fPath, _ := os.Getwd()
	fPath = fPath + "/cardsExport.json"
	type args struct {
		filePath string
	}
	tests := []struct {
		name      string
		args      args
		wantCards []*Card
		wantErr   bool
	}{{
		name: name + "-TestImportJson",
		args: args{
			filePath: fPath,
		},
		wantCards: gCardSvc.Cards,
		wantErr:   false,
	}}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotCards, err := ImportJson(tt.args.filePath)
			if (err != nil) != tt.wantErr {
				t.Errorf("ImportJson() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotCards, tt.wantCards) {
				t.Errorf("ImportJson() gotCards = %v, want %v", gotCards, tt.wantCards)
			}
		})
	}
}
