package charset

import (
	"strings"
	"testing"
)

// Round 1: 定数の文字数検証
func TestConstantLengths(t *testing.T) {
	tests := []struct {
		name     string
		value    string
		expected int
	}{
		{"Lower", Lower, 25},
		{"Upper", Upper, 24},
		{"Digits", Digits, 8},
		{"Symbols", Symbols, 4},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := len(tt.value); got != tt.expected {
				t.Errorf("%s: len = %d, want %d", tt.name, got, tt.expected)
			}
		})
	}
}

// Round 2: 曖昧文字の不在検証
func TestNoAmbiguousCharacters(t *testing.T) {
	ambiguous := "lIO01"
	sets := []struct {
		name  string
		value string
	}{
		{"Lower", Lower},
		{"Upper", Upper},
		{"Digits", Digits},
		{"Symbols", Symbols},
	}
	for _, c := range ambiguous {
		for _, s := range sets {
			if strings.ContainsRune(s.value, c) {
				t.Errorf("ambiguous character %q found in %s", c, s.name)
			}
		}
	}
}

// Round 3: All() 関数
func TestAll(t *testing.T) {
	all := All()
	expected := Lower + Upper + Digits + Symbols
	if all != expected {
		t.Errorf("All() = %q, want %q", all, expected)
	}
	if len(all) != 61 {
		t.Errorf("All() length = %d, want 61", len(all))
	}
}

// Round 4: Categories() 関数
func TestCategories(t *testing.T) {
	cats := Categories()
	if len(cats) != 4 {
		t.Fatalf("Categories() length = %d, want 4", len(cats))
	}
	expected := []struct {
		index int
		name  string
		value string
	}{
		{0, "Lower", Lower},
		{1, "Upper", Upper},
		{2, "Digits", Digits},
		{3, "Symbols", Symbols},
	}
	for _, e := range expected {
		if cats[e.index] != e.value {
			t.Errorf("Categories()[%d] = %q, want %s (%q)", e.index, cats[e.index], e.name, e.value)
		}
	}
}

// Categories() が毎回新しいスライスを返すことを検証
func TestCategoriesReturnsNewSlice(t *testing.T) {
	cats1 := Categories()
	cats2 := Categories()
	cats1[0] = "modified"
	if cats2[0] == "modified" {
		t.Error("Categories() should return a new slice each call")
	}
}

// Round 5: Exclude() 関数
func TestExclude(t *testing.T) {
	tests := []struct {
		name     string
		base     string
		excluded string
		expected string
	}{
		{"basic exclusion", "abcdef", "bd", "acef"},
		{"empty excluded", "abc", "", "abc"},
		{"empty base", "", "abc", ""},
		{"non-existent chars", "abc", "xyz", "abc"},
		{"all excluded", "abc", "abc", ""},
		{"duplicate in excluded", "abcabc", "a", "bcbc"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := Exclude(tt.base, tt.excluded)
			if got != tt.expected {
				t.Errorf("Exclude(%q, %q) = %q, want %q", tt.base, tt.excluded, got, tt.expected)
			}
		})
	}
}
