package structs

type PieChartDataPoint struct {
	Label string  `json:"label"`
	Value float64 `json:"value"`
}

type PieChartData struct {
	Data []PieChartDataPoint `json:"data"`
}
