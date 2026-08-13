package models

// MaxCommentLength is the maximum length of a comment, matching the column definition below. Writing
// a longer comment fails at the database on MySQL/Postgres, so callers that accept a comment from a
// user should reject an over-length one up front rather than letting it blow up mid-write.
const MaxCommentLength = 500

type Comment struct {
	BaseModel
	Comment        string    `gorm:"type:varchar(500); not null" json:"comment"`
	Receipt        Receipt   `json:"-"`
	ReceiptId      uint      `json:"receiptId"`
	User           User      `json:"-"`
	UserId         *uint     `json:"userId"`
	AdditionalInfo string    `gorm:"type:varchar(500)" json:"additionalInfo"`
	CommentId      *uint     `json:"commentId"`
	Replies        []Comment `json:"replies"`
}
