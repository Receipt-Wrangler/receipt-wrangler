package utils

import "reflect"

func Contains(slice []interface{}, target interface{}) bool {
	for _, value := range slice {
		if reflect.DeepEqual(value, target) {
			return true
		}
	}

	return false
}
