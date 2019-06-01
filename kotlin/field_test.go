package kotlin

import (
	"testing"

	"github.com/AntonRutkevich/astparser"
)

func Test_jsonTagFieldNameConverter(t *testing.T) {
	tests := []struct {
		name string
		def astparser.FieldDef
		want string
	}{
		{
			"single word",
			astparser.FieldDef{JsonName:"name"},
			"name",
		},
		{
			"single word first upper case",
			astparser.FieldDef{JsonName:"Name"},
			"name",
		},
		{
			"multi word with underscore first upper case",
			astparser.FieldDef{JsonName:"Name_name"},
			"nameName",
		},
		{
			"multi word with underscore first lower case",
			astparser.FieldDef{JsonName:"name_name"},
			"nameName",
		},
		{
			"multi word with underscore first & second upper case",
			astparser.FieldDef{JsonName:"Name_Name"},
			"nameName",
		},
		{
			"multi word first & second upper case",
			astparser.FieldDef{JsonName:"NameName"},
			"nameName",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := jsonTagFieldNameConverter(tt.def); got != tt.want {
				t.Errorf("jsonTagFieldNameConverter() = %v, want %v", got, tt.want)
			}
		})
	}
}
