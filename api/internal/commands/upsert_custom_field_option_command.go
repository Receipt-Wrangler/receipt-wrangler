package commands

type UpsertCustomFieldOptionCommand struct {
	// Id of an existing option to update. Zero means "add a new option".
	Id uint `json:"id"`
	// Value is the option's display text.
	Value string `json:"value"`
	// CustomFieldId is carried by the contract but never trusted: create derives
	// the foreign key from the parent association, and update derives it from the
	// request URL, so an option can never be re-parented by a client.
	CustomFieldId uint `json:"customFieldId"`
}
