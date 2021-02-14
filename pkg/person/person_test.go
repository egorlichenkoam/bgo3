package person

import (
	"os"
	"reflect"
	"testing"
)

var gPersonSvc *Service

func TestExportJson(t *testing.T) {
	gPersonSvc = TestData()
	type args struct {
		persons []*Person
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{{
		name: name + "-TestExportJson",
		args: args{
			persons: gPersonSvc.Persons,
		},
		wantErr: false,
	}}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := ExportJson(tt.args.persons); (err != nil) != tt.wantErr {
				t.Errorf("ExportJson() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestImportJson(t *testing.T) {
	fPath, _ := os.Getwd()
	fPath = fPath + "/personsExport.json"
	type args struct {
		filePath string
	}
	tests := []struct {
		name        string
		args        args
		wantPersons []*Person
		wantErr     bool
	}{{
		name: name + "-TestImportJson",
		args: args{
			filePath: fPath,
		},
		wantPersons: gPersonSvc.Persons,
		wantErr:     false,
	}}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotCards, err := ImportJson(tt.args.filePath)
			if (err != nil) != tt.wantErr {
				t.Errorf("ImportJson() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotCards, tt.wantPersons) {
				t.Errorf("ImportJson() gotCards = %v, want %v", gotCards, tt.wantPersons)
			}
		})
	}
}
