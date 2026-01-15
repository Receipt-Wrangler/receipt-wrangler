package services

import (
	"gorm.io/gorm"
	"receipt-wrangler/api/internal/commands"
	"receipt-wrangler/api/internal/models"
	"receipt-wrangler/api/internal/repositories"
	"receipt-wrangler/api/internal/structs"

	"github.com/shopspring/decimal"
)

type PieChartService struct {
	BaseService
}

func NewPieChartService(tx *gorm.DB) PieChartService {
	service := PieChartService{BaseService: BaseService{
		DB: repositories.GetDB(),
		TX: tx,
	}}
	return service
}

func (service PieChartService) GetPieChartData(
	userId uint,
	groupId string,
	command commands.PieChartDataCommand,
) (structs.PieChartData, error) {
	receiptRepository := repositories.NewReceiptRepository(service.TX)

	pagedRequest := commands.ReceiptPagedRequestCommand{
		PagedRequestCommand: commands.PagedRequestCommand{
			Page:          -1,
			PageSize:      -1,
			OrderBy:       "date",
			SortDirection: commands.DESCENDING,
		},
		Filter: command.Filter,
	}

	receipts, _, err := receiptRepository.GetPagedReceiptsByGroupId(
		userId,
		groupId,
		pagedRequest,
		[]string{"Categories", "Tags"},
	)
	if err != nil {
		return structs.PieChartData{}, err
	}

	pieChartData := structs.PieChartData{
		Data: []structs.PieChartDataPoint{},
	}

	switch command.ChartGrouping {
	case models.CHART_GROUPING_CATEGORIES:
		pieChartData.Data = service.groupByCategories(receipts)
	case models.CHART_GROUPING_TAGS:
		pieChartData.Data = service.groupByTags(receipts)
	case models.CHART_GROUPING_PAIDBY:
		pieChartData.Data, err = service.groupByPaidBy(receipts)
		if err != nil {
			return structs.PieChartData{}, err
		}
	}

	return pieChartData, nil
}

func (service PieChartService) groupByCategories(receipts []models.Receipt) []structs.PieChartDataPoint {
	categoryAmounts := make(map[string]decimal.Decimal)

	for _, receipt := range receipts {
		if len(receipt.Categories) == 0 {
			if _, exists := categoryAmounts["Uncategorized"]; !exists {
				categoryAmounts["Uncategorized"] = decimal.NewFromInt(0)
			}
			categoryAmounts["Uncategorized"] = categoryAmounts["Uncategorized"].Add(receipt.Amount)
		} else {
			for _, category := range receipt.Categories {
				if _, exists := categoryAmounts[category.Name]; !exists {
					categoryAmounts[category.Name] = decimal.NewFromInt(0)
				}
				categoryAmounts[category.Name] = categoryAmounts[category.Name].Add(receipt.Amount)
			}
		}
	}

	return service.convertToDataPoints(categoryAmounts)
}

func (service PieChartService) groupByTags(receipts []models.Receipt) []structs.PieChartDataPoint {
	tagAmounts := make(map[string]decimal.Decimal)

	for _, receipt := range receipts {
		if len(receipt.Tags) == 0 {
			if _, exists := tagAmounts["Untagged"]; !exists {
				tagAmounts["Untagged"] = decimal.NewFromInt(0)
			}
			tagAmounts["Untagged"] = tagAmounts["Untagged"].Add(receipt.Amount)
		} else {
			for _, tag := range receipt.Tags {
				if _, exists := tagAmounts[tag.Name]; !exists {
					tagAmounts[tag.Name] = decimal.NewFromInt(0)
				}
				tagAmounts[tag.Name] = tagAmounts[tag.Name].Add(receipt.Amount)
			}
		}
	}

	return service.convertToDataPoints(tagAmounts)
}

func (service PieChartService) groupByPaidBy(receipts []models.Receipt) ([]structs.PieChartDataPoint, error) {
	userAmounts := make(map[uint]decimal.Decimal)
	userNames := make(map[uint]string)
	userRepository := repositories.NewUserRepository(service.TX)

	for _, receipt := range receipts {
		if _, exists := userAmounts[receipt.PaidByUserID]; !exists {
			userAmounts[receipt.PaidByUserID] = decimal.NewFromInt(0)
		}
		userAmounts[receipt.PaidByUserID] = userAmounts[receipt.PaidByUserID].Add(receipt.Amount)

		if _, exists := userNames[receipt.PaidByUserID]; !exists {
			user, err := userRepository.GetUserById(receipt.PaidByUserID)
			if err != nil {
				userNames[receipt.PaidByUserID] = "Unknown User"
			} else {
				if len(user.DisplayName) > 0 {
					userNames[receipt.PaidByUserID] = user.DisplayName
				} else {
					userNames[receipt.PaidByUserID] = user.Username
				}
			}
		}
	}

	dataPoints := make([]structs.PieChartDataPoint, 0, len(userAmounts))
	for userId, amount := range userAmounts {
		floatVal, _ := amount.Float64()
		dataPoints = append(dataPoints, structs.PieChartDataPoint{
			Label: userNames[userId],
			Value: floatVal,
		})
	}

	return dataPoints, nil
}

func (service PieChartService) convertToDataPoints(amounts map[string]decimal.Decimal) []structs.PieChartDataPoint {
	dataPoints := make([]structs.PieChartDataPoint, 0, len(amounts))
	for name, amount := range amounts {
		floatVal, _ := amount.Float64()
		dataPoints = append(dataPoints, structs.PieChartDataPoint{
			Label: name,
			Value: floatVal,
		})
	}
	return dataPoints
}
