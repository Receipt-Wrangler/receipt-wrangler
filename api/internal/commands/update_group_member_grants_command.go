package commands

import (
	"encoding/json"
	"net/http"
	"receipt-wrangler/api/internal/structs"
	"receipt-wrangler/api/internal/utils"
)

// UpdateGroupMemberGrantsCommand is the body of
// PUT /group/{groupId}/member/{userId}/grants. It carries only the id sets — the
// membership being edited comes from the URL, never the body, so a request cannot
// aim a grant write at a different member or group than the one the caller was
// authorized against.
//
// An empty CategoryIds (or TagIds) means "not restricted for that resource": the
// member falls back to whatever their group role allows. There is no separate
// "clear restriction" verb.
type UpdateGroupMemberGrantsCommand struct {
	CategoryIds []uint `json:"categoryGrants"`
	TagIds      []uint `json:"tagGrants"`
}

func (command *UpdateGroupMemberGrantsCommand) LoadDataFromRequest(w http.ResponseWriter, r *http.Request) error {
	bytes, err := utils.GetBodyData(w, r)
	if err != nil {
		return err
	}

	return json.Unmarshal(bytes, &command)
}

// Validate rejects duplicate ids in either set. Whether the ids EXIST, and whether
// they fall within the member's group-role ceiling, is checked in the service
// layer (it needs the database).
func (command *UpdateGroupMemberGrantsCommand) Validate() structs.ValidatorError {
	errors := structs.ValidatorError{Errors: make(map[string]string)}

	// Keyed by the REQUEST field names, not the Go field names, so a client can map
	// a validation error back to the control the user typed into.
	if hasDuplicateIds(command.CategoryIds) {
		errors.Errors["categoryGrants"] = "Duplicate category ids are not allowed"
	}

	if hasDuplicateIds(command.TagIds) {
		errors.Errors["tagGrants"] = "Duplicate tag ids are not allowed"
	}

	return errors
}

func (command *UpdateGroupMemberGrantsCommand) LoadDataFromRequestAndValidate(w http.ResponseWriter, r *http.Request) (structs.ValidatorError, error) {
	err := command.LoadDataFromRequest(w, r)
	if err != nil {
		return structs.ValidatorError{}, err
	}

	return command.Validate(), nil
}

func hasDuplicateIds(ids []uint) bool {
	seen := make(map[uint]struct{}, len(ids))
	for _, id := range ids {
		if _, ok := seen[id]; ok {
			return true
		}
		seen[id] = struct{}{}
	}

	return false
}
