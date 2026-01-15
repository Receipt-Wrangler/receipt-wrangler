package models

import (
	"database/sql/driver"
	"errors"
)

type ChartGrouping string

const (
	CHART_GROUPING_CATEGORIES ChartGrouping = "CATEGORIES"
	CHART_GROUPING_TAGS       ChartGrouping = "TAGS"
	CHART_GROUPING_PAIDBY     ChartGrouping = "PAIDBY"
)

func (chartGrouping *ChartGrouping) Scan(value string) error {
	*chartGrouping = ChartGrouping(value)
	return nil
}

func (chartGrouping ChartGrouping) Value() (driver.Value, error) {
	if chartGrouping != CHART_GROUPING_CATEGORIES &&
		chartGrouping != CHART_GROUPING_TAGS &&
		chartGrouping != CHART_GROUPING_PAIDBY {
		return nil, errors.New("invalid chart grouping")
	}
	return string(chartGrouping), nil
}
