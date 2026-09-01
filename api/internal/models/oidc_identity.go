package models

import "time"

// OidcIdentity links a local User to one identity at one OidcProvider.
//
// This is the identity anchor for the whole feature. The link is keyed on the
// ID token's `sub` claim, which is the ONLY claim OIDC guarantees is both stable
// and unique within an issuer and never reassigned. Every login after the first
// is a lookup on (OidcProviderId, Subject) and consults no other claim — so a
// user renaming themselves at the IdP, or an IdP recycling a released username,
// can never re-point an existing link.
type OidcIdentity struct {
	BaseModel

	// Subject is the ID token's `sub`. Capped at 255 so the composite unique
	// index below stays inside MySQL's 3072-byte index limit; that is also why
	// PreferredUsername is deliberately NOT part of any index.
	Subject        string `gorm:"not null;size:255;uniqueIndex:idx_oidc_identity_provider_subject" json:"subject"`
	OidcProviderId uint   `gorm:"not null;uniqueIndex:idx_oidc_identity_provider_subject;uniqueIndex:idx_oidc_identity_provider_user" json:"oidcProviderId"`

	// UserId, together with the provider, is unique: one local user holds at most
	// one identity per provider, so the profile page can never show duplicates and
	// an unlink/relink cycle cannot leave a stale row that still logs in.
	UserId uint `gorm:"not null;index;uniqueIndex:idx_oidc_identity_provider_user" json:"userId"`

	OidcProvider *OidcProvider `gorm:"foreignKey:OidcProviderId;constraint:OnDelete:CASCADE" json:"-"`
	User         *User         `gorm:"foreignKey:UserId;constraint:OnDelete:CASCADE" json:"-"`

	// PreferredUsername and Email are the last-seen values of those claims, kept
	// for display on the Connected Accounts row only. They are refreshed on every
	// login and are NEVER used to resolve an identity once the link exists.
	PreferredUsername string `json:"preferredUsername"`
	Email             string `json:"email"`

	// ProvisionedUser records that this identity created the local account. Such
	// an account only ever had a random, discarded password, so unlinking its last
	// identity would strand it with no way back in — see the unlink lockout guard
	// in the OIDC service. Kept here rather than as a User column so the User model
	// does not have to widen.
	ProvisionedUser bool `gorm:"not null;default:false" json:"provisionedUser"`

	LastLoginAt *time.Time `json:"lastLoginAt"`
}
