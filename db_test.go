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
	// t.Run("checks db url created and populated", func(t *testing.T) {
	// 	got := loadData()
	// 	want := int64(4)

	// 	if got != want {
	// 		t.Errorf("return is incorrect - got %q ", got)
	// 	}

	// })
}

func TestMalwareCheck(t *testing.T) {
	t.Run("checks db query fucn", func(t *testing.T) {
		got := malwareCheck("abc.com")
		want := "yes"

		if got != want {
			t.Errorf("return is incorrect - got %q , want %q", got, want)
		}

	})
}

// func TestSampleData(t *testing.T) {
// 	t.Run("checks db is populated", func(t *testing.T) {
// 		got := loadData()
// 		want := int64(1)

// 		if got != want {
// 			t.Errorf("return is incorrect - got %q , want %q", got, want)
// 		}

// 	})
// 	// t.Run("checks db url opens", func(t *testing.T) {
// 	// 	got := loadData()
// 	// 	// want := nil

// 	// 	if got != nil {
// 	// 		t.Errorf("return is incorrect - got %q ", got)
// 	// 	}

// 	// })
// }
