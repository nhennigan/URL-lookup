package main

import (
	"testing"

	_ "github.com/go-sql-driver/mysql"
)

func TestInit(t *testing.T) {
	t.Run("checks db is created", func(t *testing.T) {
		got := createDb()
		want := int64(1)

		if got != want {
			t.Errorf("return is incorrect - got %q , want %q", got, want)
		}

	})
	t.Run("checks db url opens", func(t *testing.T) {
		got := loadData()
		// want := nil

		if got != nil {
			t.Errorf("return is incorrect - got %q ", got)
		}

	})
}
