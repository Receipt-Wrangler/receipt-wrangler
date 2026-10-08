package services

import (
	"receipt-wrangler/api/internal/models"
	"receipt-wrangler/api/internal/reporting/render"
	"receipt-wrangler/api/internal/repositories"
	"receipt-wrangler/api/internal/structs"
	"receipt-wrangler/api/internal/utils"
	"strings"
	"time"
)

type ReceiptCsvService struct {
	CsvService

	// Location is the app time zone Added At and Resolved Date are printed in.
	// Nil reads it from System Settings when the CSV is built. A receipt's own
	// Date is a calendar day and is always printed in UTC.
	Location *time.Location
}

func NewReceiptCsvService() ReceiptCsvService {
	service := ReceiptCsvService{
		CsvService: NewCsvService(),
	}
	return service
}

func (service *ReceiptCsvService) BuildReceiptCsv(receipts []models.Receipt) (structs.ReceiptCsvResult, error) {
	csvResult := structs.ReceiptCsvResult{}

	items := make([]models.Item, 0)

	headers := []string{
		"Id",
		"Added At",
		"Receipt Date",
		"Name",
		"Paid By",
		"Amount",
		"Status",
		"Categories",
		"Tags",
		"Resolved Date",
	}
	rowData := make([][]string, 0, len(receipts))
	dateFormat := "2006-01-02"
	location := service.Location
	if location == nil {
		location = repositories.GetAppLocation()
	}

	for _, receipt := range receipts {
		resolvedDateString := ""
		if receipt.ResolvedDate != nil {
			resolvedDateString = receipt.ResolvedDate.In(location).Format(dateFormat)
		}

		for _, item := range receipt.ReceiptItems {
			items = append(items, item)
		}
		newRow := []string{
			utils.UintToString(receipt.ID),
			receipt.CreatedAt.In(location).Format(dateFormat),
			receipt.Date.UTC().Format(dateFormat),
			render.SanitizeCSVField(receipt.Name),
			render.SanitizeCSVField(receipt.PaidByUser.DisplayName),
			receipt.Amount.String(),
			string(receipt.Status),
			render.SanitizeCSVField(service.BuildCategoryString(receipt.Categories)),
			render.SanitizeCSVField(service.BuildTagString(receipt.Tags)),
			resolvedDateString,
		}
		rowData = append(rowData, newRow)
	}

	buffer, err := service.CsvService.BuildCsv(headers, rowData)
	if err != nil {
		return structs.ReceiptCsvResult{}, err
	}
	csvResult.ReceiptCsvBytes = buffer.Bytes()

	csvResult.ReceiptItemCsvBytes, err = service.BuildItemCsv(items)
	if err != nil {
		return structs.ReceiptCsvResult{}, err
	}

	return csvResult, nil
}

func (service *ReceiptCsvService) BuildItemCsv(items []models.Item) ([]byte, error) {
	headers := []string{
		"Id",
		"Receipt Id",
		"Receipt Name",
		"Receipt Date",
		"Name",
		"Charged to User",
		"Amount",
		"Status",
		"Categories",
		"Tags",
	}
	rowData := make([][]string, 0, len(items))
	dateFormat := "2006-01-02"

	for _, item := range items {
		newRow := []string{
			utils.UintToString(item.ID),
			utils.UintToString(item.ReceiptId),
			render.SanitizeCSVField(item.Receipt.Name),
			item.Receipt.Date.UTC().Format(dateFormat),
			render.SanitizeCSVField(item.Name),
			render.SanitizeCSVField(item.ChargedToUser.DisplayName),
			item.Amount.String(),
			string(item.Status),
			render.SanitizeCSVField(service.BuildCategoryString(item.Categories)),
			render.SanitizeCSVField(service.BuildTagString(item.Tags)),
		}
		rowData = append(rowData, newRow)
	}

	csvService := NewCsvService()
	buffer, err := csvService.BuildCsv(headers, rowData)
	if err != nil {
		return nil, err
	}

	return buffer.Bytes(), nil
}

func (service *ReceiptCsvService) BuildCategoryString(categories []models.Category) string {
	categoryNames := make([]string, 0, len(categories))
	for _, category := range categories {
		categoryNames = append(categoryNames, category.Name)
	}

	return strings.Join(categoryNames, ",")
}

func (service *ReceiptCsvService) BuildTagString(tags []models.Tag) string {
	tagNames := make([]string, 0, len(tags))
	for _, tag := range tags {
		tagNames = append(tagNames, tag.Name)
	}

	return strings.Join(tagNames, ",")
}

func (service *ReceiptCsvService) GetZippedCsvFiles(receipts []models.Receipt) ([]byte, error) {
	csvResult, err := service.BuildReceiptCsv(receipts)
	if err != nil {
		return nil, err
	}

	fileRepository := repositories.NewFileRepository(nil)
	zip, err := fileRepository.ZipFiles(
		[]string{"receipts.csv", "items.csv"},
		[][]byte{csvResult.ReceiptCsvBytes, csvResult.ReceiptItemCsvBytes},
	)
	if err != nil {
		return nil, err
	}

	return zip, nil
}
