package permissions

import "testing"

func TestMatches(t *testing.T) {
	tests := []struct {
		name     string
		granted  string
		required string
		want     bool
	}{
		{"exact match", "app.users.read", "app.users.read", true},
		{"exact mismatch", "app.users.read", "app.users.create", false},
		{"global wildcard", "*", "group.receipts.read", true},
		{"scope wildcard matches deeper", "group.*", "group.receipts.read", true},
		{"scope wildcard matches single", "group.*", "group.view", true},
		{"scope wildcard does not cross scope", "group.*", "app.users.read", false},
		{"scope wildcard does not match bare prefix", "group.*", "group", false},
		{"domain wildcard matches", "group.receipts.*", "group.receipts.read", true},
		{"domain wildcard does not match sibling", "group.receipts.*", "group.dashboards.read", false},
		{"domain wildcard covers sub-scope suffix", "group.receipts.*", "group.receipts.read:any", true},
		{"mid wildcard matches one segment", "group.*.read", "group.receipts.read", true},
		{"mid wildcard requires the trailing literal", "group.*.read", "group.receipts.create", false},
		{"mid wildcard matches exactly one segment", "group.*.read", "group.a.b.read", false},
		{"granted more specific than required", "group.receipts.read", "group.receipts", false},
		{"required more specific than granted", "group.receipts", "group.receipts.read", false},
		{"literal suffix exact", "group.receipts.read:any", "group.receipts.read:any", true},
		{"literal suffix no superset", "group.receipts.read:any", "group.receipts.read:paid_by", false},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if got := matches(test.granted, test.required); got != test.want {
				t.Errorf("matches(%q, %q) = %v, want %v", test.granted, test.required, got, test.want)
			}
		})
	}
}

func TestHasAll(t *testing.T) {
	tests := []struct {
		name     string
		granted  []string
		required []string
		want     bool
	}{
		{"single granted", []string{"app.users.read"}, []string{"app.users.read"}, true},
		{"single not granted", []string{"app.users.read"}, []string{"app.users.create"}, false},
		{"all granted", []string{"app.users.read", "app.users.create"}, []string{"app.users.read", "app.users.create"}, true},
		{"one missing fails AND", []string{"app.users.read"}, []string{"app.users.read", "app.users.create"}, false},
		{"wildcard grant satisfies all", []string{"group.*"}, []string{"group.receipts.read", "group.view"}, true},
		{"empty granted denies", []string{}, []string{"app.users.read"}, false},
		{"no required denies", []string{"app.users.read"}, []string{}, false},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if got := HasAll(test.granted, test.required...); got != test.want {
				t.Errorf("HasAll(%v, %v) = %v, want %v", test.granted, test.required, got, test.want)
			}
		})
	}
}

func TestHasAny(t *testing.T) {
	tests := []struct {
		name     string
		granted  []string
		required []string
		want     bool
	}{
		{"one of many granted", []string{"app.users.read"}, []string{"app.users.create", "app.users.read"}, true},
		{"none granted", []string{"app.users.read"}, []string{"app.users.create", "app.users.delete"}, false},
		{"wildcard satisfies one", []string{"group.receipts.*"}, []string{"group.view", "group.receipts.read"}, true},
		{"empty granted denies", []string{}, []string{"app.users.read"}, false},
		{"no required denies", []string{"app.users.read"}, []string{}, false},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if got := HasAny(test.granted, test.required...); got != test.want {
				t.Errorf("HasAny(%v, %v) = %v, want %v", test.granted, test.required, got, test.want)
			}
		})
	}
}
