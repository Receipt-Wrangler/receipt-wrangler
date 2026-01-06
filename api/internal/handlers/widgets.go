package handlers

import (
	"net/http"
	"receipt-wrangler/api/internal/commands"
	"receipt-wrangler/api/internal/constants"
	"receipt-wrangler/api/internal/models"
	"receipt-wrangler/api/internal/repositories"
	"receipt-wrangler/api/internal/structs"
	"receipt-wrangler/api/internal/utils"

	"github.com/go-chi/chi/v5"
	"github.com/shopspring/decimal"
)

func GetPieChartData(w http.ResponseWriter, r *http.Request) {
	groupId := chi.URLParam(r, "groupId")

	handler := structs.Handler{
		ErrorMessage: "Error getting pie chart data",
		Writer:       w,
		Request:      r,
		GroupId:      groupId,
		GroupRole:    models.VIEWER,
		ResponseType: constants.ApplicationJson,
		HandlerFunction: func(w http.ResponseWriter, r *http.Request) (int, error) {
			command := commands.PieChartDataCommand{}
			err := command.LoadDataFromRequest(w, r)
			if err != nil {
				return http.StatusInternalServerError, err
			}

			vErr := command.Validate()
			if len(vErr.Errors) > 0 {
				structs.WriteValidatorErrorResponse(w, vErr, http.StatusBadRequest)
				return 0, nil
			}

			token := structs.GetClaims(r)
			receiptRepository := repositories.NewReceiptRepository(nil)

			// Build paged request from the filter with page -1 to get all receipts
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
				token.UserId,
				groupId,
				pagedRequest,
				[]string{"Categories", "Tags"},
			)
			if err != nil {
				return http.StatusInternalServerError, err
			}

			pieChartData := structs.PieChartData{
				Data: []structs.PieChartDataPoint{},
			}

			switch command.ChartGrouping {
			case models.CHART_GROUPING_CATEGORIES:
				pieChartData.Data = groupByCategories(receipts)
			case models.CHART_GROUPING_TAGS:
				pieChartData.Data = groupByTags(receipts)
			case models.CHART_GROUPING_PAIDBY:
				pieChartData.Data, err = groupByPaidBy(receipts)
				if err != nil {
					return http.StatusInternalServerError, err
				}
			}

			bytes, err := utils.MarshalResponseData(pieChartData)
			if err != nil {
				return http.StatusInternalServerError, err
			}

			w.WriteHeader(http.StatusOK)
			w.Write(bytes)

			return 0, nil
		},
	}

	HandleRequest(handler)
}

func groupByCategories(receipts []models.Receipt) []structs.PieChartDataPoint {
	categoryAmounts := make(map[string]decimal.Decimal)

	for _, receipt := range receipts {
		if len(receipt.Categories) == 0 {
			// Add to "Uncategorized" if no categories
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

	return convertToDataPoints(categoryAmounts)
}

func groupByTags(receipts []models.Receipt) []structs.PieChartDataPoint {
	tagAmounts := make(map[string]decimal.Decimal)

	for _, receipt := range receipts {
		if len(receipt.Tags) == 0 {
			// Add to "Untagged" if no tags
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

	return convertToDataPoints(tagAmounts)
}

func groupByPaidBy(receipts []models.Receipt) ([]structs.PieChartDataPoint, error) {
	userAmounts := make(map[uint]decimal.Decimal)
	userNames := make(map[uint]string)
	userRepository := repositories.NewUserRepository(nil)

	for _, receipt := range receipts {
		if _, exists := userAmounts[receipt.PaidByUserID]; !exists {
			userAmounts[receipt.PaidByUserID] = decimal.NewFromInt(0)
		}
		userAmounts[receipt.PaidByUserID] = userAmounts[receipt.PaidByUserID].Add(receipt.Amount)

		// Get user display name if not cached
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

func convertToDataPoints(amounts map[string]decimal.Decimal) []structs.PieChartDataPoint {
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
