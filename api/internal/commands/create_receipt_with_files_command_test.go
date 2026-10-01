package commands

import (
	"bytes"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"net/textproto"
	"strings"
	"testing"
)

const withFilesReceiptJson = `{"name":"Lunch","amount":"12.5","date":"2026-09-01T00:00:00Z","groupId":3,"paidByUserId":7,"status":"OPEN","comments":[{"comment":"Team lunch","userId":7}]}`

type withFilesPart struct {
	receiptValue    *string // written as a plain form value when non-nil
	receiptFilePart *string // written as an application/json file part when non-nil
	files           []withFilesUpload
}

type withFilesUpload struct {
	name  string
	bytes []byte
}

func strPtr(s string) *string { return &s }

// withFilesRequest builds a multipart CreateReceiptWithFiles request. The form
// value is written before the file part, mirroring how a client sending both
// would order them.
func withFilesRequest(t *testing.T, parts withFilesPart) *http.Request {
	t.Helper()
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	if parts.receiptValue != nil {
		if err := writer.WriteField("receipt", *parts.receiptValue); err != nil {
			t.Fatalf("write receipt field: %v", err)
		}
	}
	if parts.receiptFilePart != nil {
		header := make(textproto.MIMEHeader)
		header.Set("Content-Disposition", `form-data; name="receipt"; filename="receipt.json"`)
		header.Set("Content-Type", "application/json")
		part, err := writer.CreatePart(header)
		if err != nil {
			t.Fatalf("create receipt part: %v", err)
		}
		part.Write([]byte(*parts.receiptFilePart))
	}
	for _, file := range parts.files {
		part, err := writer.CreateFormFile("files", file.name)
		if err != nil {
			t.Fatalf("create file part: %v", err)
		}
		part.Write(file.bytes)
	}
	if err := writer.Close(); err != nil {
		t.Fatalf("close writer: %v", err)
	}

	r := httptest.NewRequest("POST", "/api/receipt/withFiles", body)
	r.Header.Set("Content-Type", writer.FormDataContentType())
	return r
}

func loadWithFiles(t *testing.T, r *http.Request) (CreateReceiptWithFilesCommand, error) {
	t.Helper()
	command := CreateReceiptWithFilesCommand{}
	err := command.LoadDataFromRequest(r)
	return command, err
}

func assertParsedReceipt(t *testing.T, command CreateReceiptWithFilesCommand, wantName string) {
	t.Helper()
	if command.Receipt.Name != wantName {
		t.Errorf("name = %q, want %q", command.Receipt.Name, wantName)
	}
	if command.Receipt.GroupId != 3 || command.Receipt.PaidByUserID != 7 {
		t.Errorf("groupId/paidBy = %d/%d, want 3/7", command.Receipt.GroupId, command.Receipt.PaidByUserID)
	}
	if command.Receipt.Amount.String() != "12.5" {
		t.Errorf("amount = %s, want 12.5", command.Receipt.Amount.String())
	}
	if len(command.Receipt.Comments) != 1 || command.Receipt.Comments[0].Comment != "Team lunch" {
		t.Errorf("comments = %+v, want the one submitted", command.Receipt.Comments)
	}
}

func TestCreateReceiptWithFilesCommand_ReceiptAsFormValue(t *testing.T) {
	command, err := loadWithFiles(t, withFilesRequest(t, withFilesPart{receiptValue: strPtr(withFilesReceiptJson)}))
	if err != nil {
		t.Fatalf("LoadDataFromRequest: %v", err)
	}
	assertParsedReceipt(t, command, "Lunch")
}

func TestCreateReceiptWithFilesCommand_ReceiptAsJsonFilePart(t *testing.T) {
	command, err := loadWithFiles(t, withFilesRequest(t, withFilesPart{receiptFilePart: strPtr(withFilesReceiptJson)}))
	if err != nil {
		t.Fatalf("LoadDataFromRequest: %v", err)
	}
	assertParsedReceipt(t, command, "Lunch")
}

// The documented order: a non-empty form value is read before the file part.
func TestCreateReceiptWithFilesCommand_FormValueWinsOverFilePart(t *testing.T) {
	fromFile := strings.Replace(withFilesReceiptJson, `"Lunch"`, `"From file part"`, 1)
	command, err := loadWithFiles(t, withFilesRequest(t, withFilesPart{
		receiptValue:    strPtr(withFilesReceiptJson),
		receiptFilePart: strPtr(fromFile),
	}))
	if err != nil {
		t.Fatalf("LoadDataFromRequest: %v", err)
	}
	assertParsedReceipt(t, command, "Lunch")
}

func TestCreateReceiptWithFilesCommand_EmptyFormValueFallsBackToFilePart(t *testing.T) {
	command, err := loadWithFiles(t, withFilesRequest(t, withFilesPart{
		receiptValue:    strPtr(""),
		receiptFilePart: strPtr(withFilesReceiptJson),
	}))
	if err != nil {
		t.Fatalf("LoadDataFromRequest: %v", err)
	}
	assertParsedReceipt(t, command, "Lunch")
}

func TestCreateReceiptWithFilesCommand_Errors(t *testing.T) {
	cases := []struct {
		name    string
		request func(t *testing.T) *http.Request
		wantErr string
	}{
		{"missing receipt", func(t *testing.T) *http.Request {
			return withFilesRequest(t, withFilesPart{})
		}, "receipt is required"},
		{"empty receipt and no file part", func(t *testing.T) *http.Request {
			return withFilesRequest(t, withFilesPart{receiptValue: strPtr("")})
		}, "receipt is required"},
		{"malformed receipt value", func(t *testing.T) *http.Request {
			return withFilesRequest(t, withFilesPart{receiptValue: strPtr(`{"name":`)})
		}, "invalid receipt"},
		{"malformed receipt file part", func(t *testing.T) *http.Request {
			return withFilesRequest(t, withFilesPart{receiptFilePart: strPtr(`not json`)})
		}, "invalid receipt"},
		{"not multipart", func(t *testing.T) *http.Request {
			r := httptest.NewRequest("POST", "/api/receipt/withFiles", strings.NewReader(withFilesReceiptJson))
			r.Header.Set("Content-Type", "application/json")
			return r
		}, ""},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			_, err := loadWithFiles(t, tc.request(t))
			if err == nil {
				t.Fatal("expected an error")
			}
			if tc.wantErr != "" && !strings.Contains(err.Error(), tc.wantErr) {
				t.Errorf("error = %q, want it to contain %q", err.Error(), tc.wantErr)
			}
		})
	}
}

// No files is a valid create; the slice is empty, not nil, so callers can range
// over it and a count of zero means "no images".
func TestCreateReceiptWithFilesCommand_NoFilesIsEmptyNotNil(t *testing.T) {
	command, err := loadWithFiles(t, withFilesRequest(t, withFilesPart{receiptValue: strPtr(withFilesReceiptJson)}))
	if err != nil {
		t.Fatalf("LoadDataFromRequest: %v", err)
	}
	if command.Files == nil {
		t.Fatal("files = nil, want an empty slice")
	}
	if len(command.Files) != 0 {
		t.Errorf("files = %d, want 0", len(command.Files))
	}
}

func TestCreateReceiptWithFilesCommand_FilesKeepOrderNamesAndBytes(t *testing.T) {
	uploads := []withFilesUpload{
		{"front.jpg", []byte{0xFF, 0xD8, 0x01, 0x02}},
		{"back.png", []byte("second file bytes")},
		{"extra.pdf", bytes.Repeat([]byte{0x42}, 4096)},
	}
	command, err := loadWithFiles(t, withFilesRequest(t, withFilesPart{
		receiptValue: strPtr(withFilesReceiptJson),
		files:        uploads,
	}))
	if err != nil {
		t.Fatalf("LoadDataFromRequest: %v", err)
	}

	if len(command.Files) != len(uploads) {
		t.Fatalf("files = %d, want %d", len(command.Files), len(uploads))
	}
	for i, want := range uploads {
		got := command.Files[i]
		if got.Name != want.name {
			t.Errorf("files[%d].name = %q, want %q", i, got.Name, want.name)
		}
		if !bytes.Equal(got.Bytes, want.bytes) {
			t.Errorf("files[%d] bytes differ (%d vs %d bytes)", i, len(got.Bytes), len(want.bytes))
		}
	}
}
