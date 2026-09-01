package oidc

import (
	"errors"
	"testing"

	"receipt-wrangler/api/internal/models"
	"receipt-wrangler/api/internal/repositories"
	"receipt-wrangler/api/internal/utils"

	"gorm.io/gorm"
)

func claims(subject string, preferredUsername string) idTokenClaims {
	return idTokenClaims{
		Subject:           subject,
		PreferredUsername: preferredUsername,
		Name:              preferredUsername,
		Email:             preferredUsername + "@example.com",
	}
}

func TestResolveUserRefusesUnknownIdentityWhenNothingIsEnabled(t *testing.T) {
	defer teardownOidcTest()
	_, provider := setupOidcTest(t, oidcTestOptions{})

	createTestUser(t, "existing")

	// Neither provisioning nor username linking is on, so an unseen identity has no
	// way in -- not even one whose username happens to match.
	_, err := resolveUser(provider, claims("unseen-subject", "existing"))
	if !errors.Is(err, ErrNoAccount) {
		t.Errorf("expected ErrNoAccount, got %v", err)
	}
}

func TestResolveUserProvisionsWhenEnabled(t *testing.T) {
	defer teardownOidcTest()
	_, provider := setupOidcTest(t, oidcTestOptions{allowProvisioning: true})

	user, err := resolveUser(provider, claims("new-subject", "newperson"))
	if err != nil {
		t.Fatalf("expected provisioning to succeed, got %v", err)
	}

	if user.Username != "newperson" {
		t.Errorf("expected the username to come from preferred_username, got %q", user.Username)
	}

	identity, err := repositories.NewOidcIdentityRepository(nil).GetIdentityBySubject(provider.ID, "new-subject")
	if err != nil {
		t.Fatalf("expected an identity to be linked: %v", err)
	}

	if !identity.ProvisionedUser {
		t.Error("a provisioned account's identity must be marked as such, or the unlink lockout guard cannot fire")
	}
}

// TestProvisionedPasswordIsUnusable is the single most important test in this
// file.
//
// UserRepository.CreateUser bcrypts whatever it is handed, so a sentinel password
// would become a WORKING password for anyone who knew the sentinel -- through the
// normal login form and through the MCP OAuth login form, which shares
// services.LoginUser. The value handed over must be real randomness, and must be
// discarded.
func TestProvisionedPasswordIsUnusable(t *testing.T) {
	defer teardownOidcTest()
	_, provider := setupOidcTest(t, oidcTestOptions{allowProvisioning: true})

	subject := "provisioned-subject"
	user, err := resolveUser(provider, claims(subject, "provisioned"))
	if err != nil {
		t.Fatalf("expected provisioning to succeed, got %v", err)
	}

	var stored models.User
	if err := repositories.GetDB().Where("id = ?", user.ID).First(&stored).Error; err != nil {
		t.Fatalf("failed to load the provisioned user: %v", err)
	}

	if len(stored.Password) == 0 {
		t.Fatal("the password column must not be empty")
	}

	for _, guess := range []string{
		"",
		"password",
		"oidc",
		"!oidc",
		subject,
		stored.Username,
		provider.Name,
	} {
		if utils.VerifyPassword(stored.Password, guess) == nil {
			t.Errorf("a provisioned account must not accept the password %q", guess)
		}
	}
}

func TestResolveUserLinksByUsernameOnlyWhenEnabled(t *testing.T) {
	defer teardownOidcTest()
	_, provider := setupOidcTest(t, oidcTestOptions{linkByUsername: true})

	existing := createTestUser(t, "alice")

	user, err := resolveUser(provider, claims("alice-subject", "alice"))
	if err != nil {
		t.Fatalf("expected the username match to link, got %v", err)
	}

	if user.ID != existing.ID {
		t.Errorf("expected to land on the existing account %d, got %d", existing.ID, user.ID)
	}

	identity, err := repositories.NewOidcIdentityRepository(nil).GetIdentityBySubject(provider.ID, "alice-subject")
	if err != nil {
		t.Fatalf("expected an identity to be linked: %v", err)
	}

	if identity.ProvisionedUser {
		t.Error("linking to an existing account must not mark it as provisioned -- that account has its own password")
	}
}

// TestSecondLoginUsesSubjectNotUsername is the guarantee that makes
// linkByUsername survivable: it only ever applies to a FIRST login. Once the
// subject is stored, a rename at the identity provider cannot re-point anything.
func TestSecondLoginUsesSubjectNotUsername(t *testing.T) {
	defer teardownOidcTest()
	_, provider := setupOidcTest(t, oidcTestOptions{linkByUsername: true, allowProvisioning: true})

	existing := createTestUser(t, "bob")
	createTestUser(t, "carol")

	first, err := resolveUser(provider, claims("bob-subject", "bob"))
	if err != nil {
		t.Fatalf("first login failed: %v", err)
	}

	// The same person renames themselves at the identity provider to a name that
	// belongs to somebody else here. The subject is unchanged, so the login must
	// still land on the original account.
	second, err := resolveUser(provider, claims("bob-subject", "carol"))
	if err != nil {
		t.Fatalf("second login failed: %v", err)
	}

	if second.ID != first.ID || second.ID != existing.ID {
		t.Errorf("a rename at the identity provider re-pointed the account: first %d, second %d, expected %d", first.ID, second.ID, existing.ID)
	}
}

func TestResolveUserRefusesUsernameCollisionRatherThanSuffixing(t *testing.T) {
	defer teardownOidcTest()
	_, provider := setupOidcTest(t, oidcTestOptions{allowProvisioning: true})

	createTestUser(t, "dave")

	// Username linking is OFF, so this identity cannot attach to the existing
	// "dave". Provisioning must refuse rather than quietly create "dave-2", which
	// would look like data loss to the user.
	_, err := resolveUser(provider, claims("dave-subject", "dave"))
	if !errors.Is(err, ErrAccountExists) {
		t.Fatalf("expected ErrAccountExists, got %v", err)
	}

	var count int64
	if err := repositories.GetDB().Model(&models.User{}).Count(&count).Error; err != nil {
		t.Fatalf("failed to count users: %v", err)
	}

	if count != 1 {
		t.Errorf("expected no second account to be created, found %d users", count)
	}
}

func TestResolveUserSkipsDummyUsersOnUsernameMatch(t *testing.T) {
	defer teardownOidcTest()
	_, provider := setupOidcTest(t, oidcTestOptions{linkByUsername: true})

	dummy := models.User{Username: "placeholder", DisplayName: "placeholder", Password: "x", IsDummyUser: true}
	if err := repositories.GetDB().Create(&dummy).Error; err != nil {
		t.Fatalf("failed to create the dummy user: %v", err)
	}

	// A dummy user is blocked from logging in anyway, so attaching an identity to
	// one would only produce a confusing failure later.
	_, err := resolveUser(provider, claims("placeholder-subject", "placeholder"))
	if !errors.Is(err, ErrNoAccount) {
		t.Errorf("expected ErrNoAccount for a dummy-user match, got %v", err)
	}
}

func TestResolveUserRefusesDummyUserOnReturningLogin(t *testing.T) {
	defer teardownOidcTest()
	_, provider := setupOidcTest(t, oidcTestOptions{})

	user := createTestUser(t, "willbecomedummy")
	if err := createIdentity(provider, user.ID, claims("dummy-subject", "willbecomedummy"), false); err != nil {
		t.Fatalf("failed to seed the identity: %v", err)
	}

	if err := repositories.GetDB().Model(&models.User{}).Where("id = ?", user.ID).Update("is_dummy_user", true).Error; err != nil {
		t.Fatalf("failed to flip the dummy flag: %v", err)
	}

	_, err := resolveUser(provider, claims("dummy-subject", "willbecomedummy"))
	if !errors.Is(err, ErrUserIsDummy) {
		t.Errorf("expected ErrUserIsDummy, got %v", err)
	}
}

func TestResolveLinkAttachesWithoutMatchingOrProvisioning(t *testing.T) {
	defer teardownOidcTest()
	// Both toggles OFF: linking works anyway, because the session already proved
	// who the caller is. This is what makes linkByUsername safe to default off.
	_, provider := setupOidcTest(t, oidcTestOptions{})

	user := createTestUser(t, "erin")

	err := resolveLink(provider, claims("erin-subject", "somethingelse"), user.ID)
	if err != nil {
		t.Fatalf("expected linking to succeed, got %v", err)
	}

	identity, err := repositories.NewOidcIdentityRepository(nil).GetIdentityBySubject(provider.ID, "erin-subject")
	if err != nil {
		t.Fatalf("expected an identity: %v", err)
	}

	if identity.UserId != user.ID {
		t.Errorf("expected the identity on user %d, got %d", user.ID, identity.UserId)
	}
}

func TestResolveLinkRefusesAnIdentityOwnedByAnotherUser(t *testing.T) {
	defer teardownOidcTest()
	_, provider := setupOidcTest(t, oidcTestOptions{})

	owner := createTestUser(t, "owner")
	intruder := createTestUser(t, "intruder")

	if err := createIdentity(provider, owner.ID, claims("shared-subject", "owner"), false); err != nil {
		t.Fatalf("failed to seed the identity: %v", err)
	}

	// Re-pointing would silently transfer the account.
	err := resolveLink(provider, claims("shared-subject", "owner"), intruder.ID)
	if !errors.Is(err, ErrIdentityLinkedElsewhere) {
		t.Errorf("expected ErrIdentityLinkedElsewhere, got %v", err)
	}
}

func TestResolveLinkRefusesASecondIdentityAtTheSameProvider(t *testing.T) {
	defer teardownOidcTest()
	_, provider := setupOidcTest(t, oidcTestOptions{})

	user := createTestUser(t, "frank")
	if err := createIdentity(provider, user.ID, claims("frank-subject", "frank"), false); err != nil {
		t.Fatalf("failed to seed the identity: %v", err)
	}

	err := resolveLink(provider, claims("frank-other-subject", "frank"), user.ID)
	if !errors.Is(err, ErrAlreadyLinked) {
		t.Errorf("expected ErrAlreadyLinked, got %v", err)
	}
}

func TestUnlinkRefusesToStrandAProvisionedAccount(t *testing.T) {
	defer teardownOidcTest()
	_, provider := setupOidcTest(t, oidcTestOptions{allowProvisioning: true})

	user, err := resolveUser(provider, claims("stranded-subject", "stranded"))
	if err != nil {
		t.Fatalf("expected provisioning to succeed, got %v", err)
	}

	// This account has only a random, discarded password, so removing its last
	// identity would leave it with no way in at all.
	err = UnlinkIdentity(user.ID, provider.Name)
	if !errors.Is(err, ErrWouldLockOut) {
		t.Errorf("expected ErrWouldLockOut, got %v", err)
	}

	if _, err := repositories.NewOidcIdentityRepository(nil).GetIdentityForUser(provider.ID, user.ID); err != nil {
		t.Errorf("the identity should still exist after a refused unlink: %v", err)
	}
}

func TestUnlinkSucceedsForAnAccountThatHasAPassword(t *testing.T) {
	defer teardownOidcTest()
	_, provider := setupOidcTest(t, oidcTestOptions{})

	user := createTestUser(t, "grace")
	if err := createIdentity(provider, user.ID, claims("grace-subject", "grace"), false); err != nil {
		t.Fatalf("failed to seed the identity: %v", err)
	}

	if err := UnlinkIdentity(user.ID, provider.Name); err != nil {
		t.Fatalf("expected the unlink to succeed, got %v", err)
	}

	_, err := repositories.NewOidcIdentityRepository(nil).GetIdentityForUser(provider.ID, user.ID)
	if !errors.Is(err, gorm.ErrRecordNotFound) {
		t.Errorf("expected the identity to be gone, got %v", err)
	}
}

func TestDeriveUsernameFallsBackThroughTheClaims(t *testing.T) {
	provider := models.OidcProvider{Name: "acme"}

	tests := []struct {
		name     string
		claims   idTokenClaims
		expected string
	}{
		{"prefers preferred_username", idTokenClaims{Subject: "s", PreferredUsername: "Hank", Email: "other@x.com"}, "hank"},
		{"falls back to the email local part", idTokenClaims{Subject: "s", Email: "Ivy.Jones@x.com"}, "ivy.jones"},
		{"falls back to the subject", idTokenClaims{Subject: "abc123"}, "acme-abc123"},
		{"strips characters a username may not contain", idTokenClaims{Subject: "s", PreferredUsername: "a b/c!d"}, "abcd"},
		{"skips a candidate that sanitizes too short", idTokenClaims{Subject: "s", PreferredUsername: "a", Email: "longenough@x.com"}, "longenough"},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if actual := deriveUsername(provider, test.claims); actual != test.expected {
				t.Errorf("expected %q, got %q", test.expected, actual)
			}
		})
	}
}
