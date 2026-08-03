package commands

import (
	"encoding/json"
	"fmt"
	"net/http"
	"receipt-wrangler/api/internal/permissions"
	"receipt-wrangler/api/internal/structs"
	"receipt-wrangler/api/internal/utils"
)

type UpsertRoleCommand struct {
	Name        string            `json:"name"`
	Description string            `json:"description"`
	Scope       permissions.Scope `json:"scope"`
	Permissions []string          `json:"permissions"`
	// CategoryGrants / TagGrants restrict which categories/tags members of a
	// group role may read and use. They are only valid on GROUP-scoped roles. An
	// empty set means "unrestricted" (see every category/tag) — restriction is
	// opt-in. The values are category/tag ids; their existence is validated
	// against the database in the role service.
	CategoryGrants []uint `json:"categoryGrants"`
	TagGrants      []uint `json:"tagGrants"`
	// PaidByUserGrants / IncludeOwnPaidReceipts restrict which receipts members of
	// a group role may see, by the receipt's "paid by" user. Group-scoped only and
	// opt-in: an empty set with IncludeOwnPaidReceipts == false means unrestricted
	// (see every payer's receipts). PaidByUserGrants are user ids (validated to
	// exist in the role service); IncludeOwnPaidReceipts adds the current member's
	// own receipts. Selecting only specific users (IncludeOwnPaidReceipts false)
	// therefore hides the member's own receipts too.
	PaidByUserGrants       []uint `json:"paidByUserGrants"`
	IncludeOwnPaidReceipts bool   `json:"includeOwnPaidReceipts"`
	// SeesAllMembers is the "supervisor" exemption for member-presence isolation:
	// holders of this group role see every member of an isolated group and are
	// visible to every member. Group-scoped only; ignored (rejected) on APP scope,
	// mirroring the paid-by flags. Default false ⇒ no effect on existing roles.
	SeesAllMembers bool `json:"seesAllMembers"`
	// SkipDefaultGroupCreation suppresses the personal "My Receipts" group that is
	// otherwise created for every new user assigned this role. App-scoped only;
	// rejected on GROUP scope (a group role is assigned to a membership, long
	// after the user was created). Default false ⇒ no effect on existing roles.
	SkipDefaultGroupCreation bool `json:"skipDefaultGroupCreation"`
	// ReportTemplateGrants restrict which report templates members of a group role
	// may act on, per action. Group-scoped only and opt-in: an empty set means
	// unrestricted (act on every template the role's group access reaches). Each
	// entry names a template id and the subset of actions
	// (read/generate/update/delete/duplicate) the role may perform on it. Template
	// existence is validated against the database in the role service.
	ReportTemplateGrants []ReportTemplateGrantCommand `json:"reportTemplateGrants"`
}

// ReportTemplateGrantCommand is one row of the report-template access matrix: the
// set of actions a group role may perform on a single template.
type ReportTemplateGrantCommand struct {
	ReportTemplateId uint     `json:"reportTemplateId"`
	Permissions      []string `json:"permissions"`
}

// scopableReportTemplateActions is the set of report actions that can be granted
// per template. "create" is excluded: it makes a new template, so there is no
// existing template to scope it against.
var scopableReportTemplateActions = map[string]bool{
	"read":      true,
	"generate":  true,
	"update":    true,
	"delete":    true,
	"duplicate": true,
}

func (command *UpsertRoleCommand) LoadDataFromRequest(w http.ResponseWriter, r *http.Request) error {
	bytes, err := utils.GetBodyData(w, r)
	if err != nil {
		return err
	}

	err = json.Unmarshal(bytes, &command)
	if err != nil {
		return err
	}

	return nil
}

func (command *UpsertRoleCommand) Validate() structs.ValidatorError {
	errors := make(map[string]string)
	vErr := structs.ValidatorError{}

	if len(command.Name) == 0 {
		errors["name"] = "Name is required"
	}

	if command.Scope != permissions.ScopeApp && command.Scope != permissions.ScopeGroup {
		errors["scope"] = "Scope must be either APP or GROUP"
	}

	seen := make(map[string]bool)
	for _, permission := range command.Permissions {
		if seen[permission] {
			errors["permissions"] = fmt.Sprintf("Duplicate permission: %s", permission)
			break
		}
		seen[permission] = true

		descriptor, ok := permissions.Get(permission)
		if !ok {
			errors["permissions"] = fmt.Sprintf("Unknown permission: %s", permission)
			break
		}

		if descriptor.Scope != command.Scope {
			errors["permissions"] = fmt.Sprintf("Permission %s does not belong to scope %s", permission, command.Scope)
			break
		}
	}

	// Category/tag/paid-by grants and member-visibility settings are a group-role
	// concept (they slice the global pool / narrow visibility per group role); they
	// make no sense on an app role.
	if command.Scope == permissions.ScopeApp &&
		(len(command.CategoryGrants) > 0 || len(command.TagGrants) > 0 ||
			len(command.PaidByUserGrants) > 0 || command.IncludeOwnPaidReceipts ||
			command.SeesAllMembers) {
		errors["grants"] = "Category, tag, paid-by, and member-visibility settings are only valid on group roles"
	}

	// Skipping personal-group creation is an app-role concept: it acts when the
	// user account is created, whereas a group role is assigned to an existing
	// membership long after that.
	if command.Scope == permissions.ScopeGroup && command.SkipDefaultGroupCreation {
		errors["skipDefaultGroupCreation"] = "Skipping default group creation is only valid on application roles"
	}

	if hasDuplicateUint(command.CategoryGrants) {
		errors["categoryGrants"] = "Duplicate category grant"
	}

	if hasDuplicateUint(command.TagGrants) {
		errors["tagGrants"] = "Duplicate tag grant"
	}

	if hasDuplicateUint(command.PaidByUserGrants) {
		errors["paidByUserGrants"] = "Duplicate paid-by user grant"
	}

	// Report template grants are likewise a group-role concept; on an app role the
	// scope violation is the primary problem, so it wins over structural checks.
	if command.Scope == permissions.ScopeApp && len(command.ReportTemplateGrants) > 0 {
		errors["reportTemplateGrants"] = "Report template grants are only valid on group roles"
	} else if msg := validateReportTemplateGrants(command.ReportTemplateGrants); msg != "" {
		errors["reportTemplateGrants"] = msg
	}

	vErr.Errors = errors
	return vErr
}

// validateReportTemplateGrants checks a group role's report-template grant entries:
// no duplicate template id, each entry names at least one action, no duplicate
// action within an entry, and every action is one that can be scoped per template.
// It returns an empty string when the grants are well-formed.
func validateReportTemplateGrants(grants []ReportTemplateGrantCommand) string {
	seenTemplate := make(map[uint]bool, len(grants))
	for _, grant := range grants {
		if seenTemplate[grant.ReportTemplateId] {
			return fmt.Sprintf("Duplicate report template grant: %d", grant.ReportTemplateId)
		}
		seenTemplate[grant.ReportTemplateId] = true

		if len(grant.Permissions) == 0 {
			return fmt.Sprintf("Report template grant %d must list at least one action", grant.ReportTemplateId)
		}

		seenAction := make(map[string]bool, len(grant.Permissions))
		for _, action := range grant.Permissions {
			if seenAction[action] {
				return fmt.Sprintf("Duplicate action %q on report template grant %d", action, grant.ReportTemplateId)
			}
			seenAction[action] = true

			if !scopableReportTemplateActions[action] {
				return fmt.Sprintf("Unknown report template action: %s", action)
			}
		}
	}
	return ""
}

// hasDuplicateUint reports whether ids contains the same value more than once.
func hasDuplicateUint(ids []uint) bool {
	seen := make(map[uint]bool, len(ids))
	for _, id := range ids {
		if seen[id] {
			return true
		}
		seen[id] = true
	}
	return false
}
