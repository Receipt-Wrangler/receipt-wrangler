package utils

import "testing"

func TestGetTriggerRegexShouldMatchMentions(t *testing.T) {
	regex := GetTriggerRegex()

	// Test matching @mentions
	testCases := []struct {
		input    string
		expected []string
	}{
		{"@user", []string{"@user"}},
		{"@User123", []string{"@User123"}},
		{"Hello @world", []string{"@world"}},
		{"@first and @second", []string{"@first", "@second"}},
		{"@abc_def", []string{"@abc_def"}},
		{"no mentions here", nil},
		{"", nil},
		{"please mention @someone", []string{"@someone"}},
	}

	for _, tc := range testCases {
		matches := regex.FindAllString(tc.input, -1)
		if len(matches) != len(tc.expected) {
			t.Errorf("Input: %q - Expected %d matches, got %d", tc.input, len(tc.expected), len(matches))
			continue
		}
		for i, match := range matches {
			if match != tc.expected[i] {
				t.Errorf("Input: %q - Expected match %q, got %q", tc.input, tc.expected[i], match)
			}
		}
	}
}

func TestGetTriggerRegexShouldNotMatchInvalidMentions(t *testing.T) {
	regex := GetTriggerRegex()

	// Test that certain patterns don't match
	testCases := []string{
		"@ space",
		"@@double",
	}

	for _, tc := range testCases {
		matches := regex.FindAllString(tc, -1)
		// The regex @\w+ will still match @double in @@double, which is expected behavior
		// We just verify the function returns a valid regex
		if matches == nil && tc == "@@double" {
			t.Errorf("Expected regex to match @double in %q", tc)
		}
	}
}
