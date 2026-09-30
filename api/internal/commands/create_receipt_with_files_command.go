package commands

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"receipt-wrangler/api/internal/constants"
)

// ReceiptFileUpload is one image uploaded alongside a receipt create.
type ReceiptFileUpload struct {
	Name  string
	Bytes []byte
}

// CreateReceiptWithFilesCommand is the multipart body of POST /receipt/withFiles:
// the receipt (comments included) plus every image, so the receipt is created in
// one atomic call and the server sees the images at create time.
type CreateReceiptWithFilesCommand struct {
	Receipt UpsertReceiptCommand
	Files   []ReceiptFileUpload
}

// LoadDataFromRequest parses the multipart form. The `receipt` part is JSON, and
// is accepted either as a plain form value or as a file part (for example with
// an application/json content type), because generated clients differ in how
// they encode a model inside multipart/form-data. `files` may be absent or
// repeated. A malformed body returns an error the handler maps to a 400.
func (command *CreateReceiptWithFilesCommand) LoadDataFromRequest(r *http.Request) error {
	err := r.ParseMultipartForm(constants.MultipartFormMaxSize)
	if err != nil {
		return err
	}

	receiptJson, err := readMultipartJsonPart(r.MultipartForm, "receipt")
	if err != nil {
		return err
	}

	err = json.Unmarshal(receiptJson, &command.Receipt)
	if err != nil {
		return fmt.Errorf("invalid receipt: %w", err)
	}

	fileHeaders := r.MultipartForm.File["files"]
	command.Files = make([]ReceiptFileUpload, 0, len(fileHeaders))
	for _, fileHeader := range fileHeaders {
		fileBytes, err := readMultipartFile(fileHeader)
		if err != nil {
			return err
		}

		command.Files = append(command.Files, ReceiptFileUpload{
			Name:  fileHeader.Filename,
			Bytes: fileBytes,
		})
	}

	return nil
}

// readMultipartJsonPart returns the named part's raw bytes, whether the client
// sent it as a form value or as a file part.
func readMultipartJsonPart(form *multipart.Form, name string) ([]byte, error) {
	if values := form.Value[name]; len(values) > 0 && len(values[0]) > 0 {
		return []byte(values[0]), nil
	}

	if fileHeaders := form.File[name]; len(fileHeaders) > 0 {
		return readMultipartFile(fileHeaders[0])
	}

	return nil, errors.New(name + " is required")
}

func readMultipartFile(fileHeader *multipart.FileHeader) ([]byte, error) {
	file, err := fileHeader.Open()
	if err != nil {
		return nil, err
	}
	defer file.Close()

	return io.ReadAll(file)
}
