package llm

import (
	"strings"
	"testing"
	"time"
)

func TestCleanJSONContent(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "Markdown code block with json language tag",
			input:    "```json\n{\"action\": \"create\", \"amount\": 350}\n```",
			expected: "{\"action\": \"create\", \"amount\": 350}",
		},
		{
			name:     "Markdown code block without language tag",
			input:    "```\n{\"action\": \"list\"}\n```",
			expected: "{\"action\": \"list\"}",
		},
		{
			name:     "Clean JSON without code blocks",
			input:    "{\"action\": \"delete\", \"target\": \"такси\"}",
			expected: "{\"action\": \"delete\", \"target\": \"такси\"}",
		},
		{
			name:     "JSON with surrounding spaces and newlines",
			input:    "   \n\n```json\n  {\"action\": \"update\"}  \n```  \n ",
			expected: "{\"action\": \"update\"}",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			result := cleanJSONContent(tc.input)
			if result != tc.expected {
				t.Errorf("cleanJSONContent() = %q; want %q", result, tc.expected)
			}
		})
	}
}

func TestBuildSystemPrompt(t *testing.T) {
	fixedDate := time.Date(2026, time.October, 19, 12, 0, 0, 0, time.UTC)
	prompt := BuildSystemPrompt(fixedDate)

	// Verify current and yesterday dates are present
	if !strings.Contains(prompt, "2026-10-19") {
		t.Errorf("expected prompt to contain current date 2026-10-19, got: %s", prompt)
	}
	if !strings.Contains(prompt, "2026-10-18") {
		t.Errorf("expected prompt to contain yesterday date 2026-10-18, got: %s", prompt)
	}

	// Verify all 7 required categories are mentioned
	requiredCategories := []string{"еда", "транспорт", "жильё", "развлечения", "здоровье", "покупки", "другое"}
	for _, cat := range requiredCategories {
		if !strings.Contains(prompt, cat) {
			t.Errorf("expected prompt to mention category %q", cat)
		}
	}

	// Verify actions are described
	requiredActions := []string{"create", "list", "update", "delete"}
	for _, act := range requiredActions {
		if !strings.Contains(prompt, act) {
			t.Errorf("expected prompt to mention action %q", act)
		}
	}
}
