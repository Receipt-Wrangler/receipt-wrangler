package structs

// ReceiptRequirements is what a caller must supply on a group's receipts, resolved
// from their group role's RequireReceiptComment / RequireReceiptImage flags with
// the group's waivers already applied (see ReceiptService.ResolveReceiptRequirements).
type ReceiptRequirements struct {
	CommentRequired bool `json:"commentRequired"`
	ImageRequired   bool `json:"imageRequired"`
}

// Any reports whether anything is required at all.
func (requirements ReceiptRequirements) Any() bool {
	return requirements.CommentRequired || requirements.ImageRequired
}
