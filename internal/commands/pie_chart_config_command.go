package commands

import "receipt-wrangler/api/internal/models"

type PieChartConfigCommand struct {
	ReceiptFilter ReceiptPagedRequestFilter
	PieChartType  models.PieChartType
}
