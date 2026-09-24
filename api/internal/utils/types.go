package utils

import (
	"strconv"
	"strings"
)

func UintToString(v uint) string {
	return strconv.FormatUint(uint64(v), 10)
}

func StringToUint(v string) (uint, error) {
	vTrimmed := strings.Trim(v, " ")
	result, err := strconv.ParseUint(vTrimmed, 10, 32)
	if err != nil {
		return 0, err
	}

	return uint(result), nil
}

func StringToUint64(v string) (uint64, error) {
	vTrimmed := strings.Trim(v, " ")
	result, err := strconv.ParseUint(vTrimmed, 10, 64)
	if err != nil {
		return 0, err
	}

	return result, nil
}

func StringToInt(v string) (int, error) {
	result, err := strconv.Atoi(v)
	if err != nil {
		return 0, err
	}

	return result, nil
}

// FilterValueToUint coerces a JSON-decoded filter id (numbers decode to float64)
// to uint. Shared by the receipt grant filter (services) and the per-group
// category/tag filter disjunction (repositories), which cannot import services.
func FilterValueToUint(value interface{}) (uint, bool) {
	switch typed := value.(type) {
	case float64:
		return uint(typed), true
	case int:
		return uint(typed), true
	case int64:
		return uint(typed), true
	case uint:
		return typed, true
	default:
		return 0, false
	}
}
