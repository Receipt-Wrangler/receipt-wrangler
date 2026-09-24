package structs

// SystemTaskSourceFileView is the preview payload for the upload behind a failed
// activity.
//
// Deliberately not models.FileDataView: that one embeds BaseModel, whose swagger
// schema requires id and createdAt. A temp file has no database row, so it would
// ship id 0 and invite a client to treat it as a receipt image id.
type SystemTaskSourceFileView struct {
	Name         string `json:"name"`
	EncodedImage string `json:"encodedImage"`
}
