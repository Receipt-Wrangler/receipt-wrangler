package structs

import "time"

type Attachment struct {
	Filename string `json:"filename"`
	FileType string `json:"fileType"`
	Size     uint   `json:"size"`
}

type EmailMetadata struct {
	// Inbox is the IMAP login (the System Email's username) of the mailbox the
	// email was read from. Jobs queued before it existed carry it empty.
	Inbox            string       `json:"inbox"`
	Date             time.Time    `json:"date"`
	Subject          string       `json:"subject"`
	To               string       `json:"to"`
	FromName         string       `json:"fromName"`
	FromEmail        string       `json:"fromEmail"`
	Body             string       `json:"body"`
	BodyHtml         string       `json:"bodyHtml"`
	Attachments      []Attachment `json:"attachments"`
	GroupSettingsIds []uint       `json:"groupSettingsIds"`
}
