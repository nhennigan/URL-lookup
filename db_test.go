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
}

func TestSampleData(t *testing.T) {
	t.Run("checks db is populated", func(t *testing.T) {
		initializeDb()
		got, err := malwareCheck("abc.com")
		want := "yes"

		if got != want || err != nil {
			t.Errorf("return is incorrect - got %q , want %q", got, want)
		}
	})
}

var malwareCheckTests = []struct {
	name      string
	input     string
	output    string
	wantError bool
}{
	{"check db quey valid input", "abc.com", "yes", false},
	{"check db query empty input", "", "", true},
	{"check db query invalid input", "678.com", "", false},
}

func TestMalwareCheck(t *testing.T) {
	for _, tt := range malwareCheckTests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := malwareCheck(tt.input)
			if got != tt.output {
				t.Errorf("db returned val is incorrect - got %q , want %q", got, tt.output)
			}
			if err == nil && tt.wantError == true {
				t.Errorf("did not return error as expected - got %q ", got)
			}
		})
	}
}

func TestSetMalwareState(t *testing.T) {
	t.Run("checks db update func", func(t *testing.T) {
		setMalwareState("def.com", "yes")
		got, _ := malwareCheck("def.com")
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

var testAddEntry = []struct {
	name  string
	input []inputData
}{
	{"checks new entries added to db correctly",
		[]inputData{
			{"qrs.com", "no"},
			{"tuv.com", "yes"},
			{"wxy.com", "yes"},
			{"123.com", "no"}},
	},
	{"checks updated entries",
		[]inputData{
			{"qrs.com", "no"},
			{"tuv.com", "no"},
			{"wxy.com", "yes"},
			{"123.com", "no"},
			{"456.com", "yes"}},
	},
}

func TestAddNewEntry(t *testing.T) {
	for _, tt := range testAddEntry {
		t.Run(tt.name, func(t *testing.T) {
			for i, _ := range tt.input {
				addNewEntry(tt.input)
				got, _ := malwareCheck(tt.input[i].URL)
				want := tt.input[i].Malware
				if got != tt.input[i].Malware {
					t.Errorf("return is incorrect - got %q , want %q for entry %s", got, want, tt.input[i].URL)
				}
			}
		})
	}
}
