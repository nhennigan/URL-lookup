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
	// 	got := initializeData()
	// 	want := int64(4)

	// 	if got != want {
	// 		t.Errorf("return is incorrect - got %q ", got)
	// 	}

	// })
}

func TestMalwareCheck(t *testing.T) {
	t.Run("checks db query func", func(t *testing.T) {
		got := malwareCheck("abc.com")
		want := "yes"

		if got != want {
			t.Errorf("return is incorrect - got %q , want %q", got, want)
		}

	})
}
func TestSetMalwareSafe(t *testing.T) {
	t.Run("checks db update func", func(t *testing.T) {
		setMalwareSafe("def.com", "yes")
		got := malwareCheck("def.com")
		want := "yes"

		if got != want {
			t.Errorf("return is incorrect - got %q , want %q", got, want)
		}

	})
}

func TestReadNewData(t *testing.T) {
	t.Run("checks entries read in correctly", func(t *testing.T) {
		got := readNewData()
		want := []inputData{
			{"qrs.com", "no"},
			{"tuv.com", "no"},
			{"wxy.com", "yes"},
			{"123.com", "no"}}

		for i, _ := range got {
			if got[i] != want[i] {
				t.Errorf("return is incorrect - got %q , want %q", got, want)
			}
		}

	})
}

func TestAddNewEntry(t *testing.T) {
	t.Run("checks new entries added to db correctly", func(t *testing.T) {
		input := []inputData{
			{"qrs.com", "no"},
			{"tuv.com", "no"},
			{"wxy.com", "yes"},
			{"123.com", "no"}}
		addNewEntry(input)
		// got := readNewData()

		for i, _ := range input {
			got := malwareCheck(input[i].URL)
			want := input[i].Malware
			if got != input[i].Malware {
				t.Errorf("return is incorrect - got %q , want %q", got, want)
			}
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
