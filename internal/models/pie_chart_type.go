package models

import (
	"database/sql/driver"
	"errors"
)

type PieChartType string

const (
	CATEGORIES PieChartType = "CATEGORIES"
	TAGS       PieChartType = "TAGS"
	PAID_BY    PieChartType = "PAID_BY"
	STATUS     PieChartType = "STATUS"
)

func (pieChartType *PieChartType) Scan(value string) error {
	*pieChartType = PieChartType(value)
	return nil
}

func (pieChartType PieChartType) Value() (driver.Value, error) {
	if pieChartType != CATEGORIES &&
		pieChartType != TAGS &&
		pieChartType != PAID_BY &&
		pieChartType != STATUS {
		return nil, errors.New("invalid pie chart type")
	}
	return string(pieChartType), nil
}
