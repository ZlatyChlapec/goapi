package internal

import (
	"testing"
)

func Test_checkDuplicits(t *testing.T) {
	type args struct {
		tokens *[]authToken
	}
	tests := []struct {
		name string
		panic bool
		args args
	}{
		{"noDuplicities", false, args{tokens: &[]authToken{{"123", "Spidy"}, {"124", "Baty"}} }},
		{"duplicities", true, args{tokens: &[]authToken{{"123", "Spidy"}, {"124", "Spidy"}} }},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			defer func() {
				r := recover()
				if (r != nil) != tt.panic {
					t.Error("panic")
				}
			}()
			checkDuplicits(tt.args.tokens)
		})
	}
}
