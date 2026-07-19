package utils

import "testing"

func TestGetTriggerRegexShouldMatchTriggers(t *testing.T) {
	regex := GetTriggerRegex()

	input := "hey @alice and @bob, ping @charlie"
	matches := regex.FindAllString(input, -1)

	expected := []string{"@alice", "@bob", "@charlie"}
	if len(matches) != len(expected) {
		t.Fatalf("Expected %d matches, got %d (%v)", len(expected), len(matches), matches)
	}

	for i, match := range matches {
		if match != expected[i] {
			PrintTestError(t, match, expected[i])
		}
	}
}

func TestGetTriggerRegexShouldNotMatchWithoutTrigger(t *testing.T) {
	regex := GetTriggerRegex()

	if regex.MatchString("no triggers here") {
		t.Errorf("Expected no match for a string without an @-trigger")
	}
}
