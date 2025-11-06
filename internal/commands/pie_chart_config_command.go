package commands

type PieChartConfigCommand struct {
	ReceiptFilter ReceiptPagedRequestFilter
	Labels        []string
	Colors        []string
}
