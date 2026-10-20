package models

import (
	"testing"
)

func TestIsValidCategory(t *testing.T) {
	tests := []struct {
		name     string
		category string
		expected bool
	}{
		// Valid categories from specification
		{"Valid food", "еда", true},
		{"Valid transport", "транспорт", true},
		{"Valid housing", "жильё", true},
		{"Valid entertainment", "развлечения", true},
		{"Valid health", "здоровье", true},
		{"Valid shopping", "покупки", true},
		{"Valid other", "другое", true},

		// Invalid categories
		{"Invalid empty", "", false},
		{"Invalid crypto", "криптовалюта", false},
		{"Invalid travel", "путешествия", false},
		{"Invalid english", "food", false},
		{"Invalid random text", "непонятно_что", false},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			result := IsValidCategory(tc.category)
			if result != tc.expected {
				t.Errorf("IsValidCategory(%q) = %v; want %v", tc.category, result, tc.expected)
			}
		})
	}
}

func TestValidCategoriesCount(t *testing.T) {
	expectedCount := 7
	if len(ValidCategories) != expectedCount {
		t.Errorf("ValidCategories length = %d; want %d", len(ValidCategories), expectedCount)
	}
}
