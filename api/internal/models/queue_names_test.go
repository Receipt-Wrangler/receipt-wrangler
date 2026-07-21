package models

import (
	"receipt-wrangler/api/internal/utils"
	"testing"
)

func TestQueueName_Value(t *testing.T) {
	valid := []QueueName{
		QuickScanQueue,
		EmailPollingQueue,
		EmailReceiptProcessingQueue,
		EmailReceiptImageCleanupQueue,
		SystemCleanUpQueue,
	}
	for _, v := range valid {
		assertValuerValid(t, string(v), v, string(v))
	}
}

func TestQueueName_Value_Invalid(t *testing.T) {
	// This type has no empty-string exception, so an empty value is invalid too.
	assertValuerInvalid(t, "empty", QueueName(""))
	assertValuerInvalid(t, "bogus", QueueName("bogus"))
}

func TestQueueName_Scan(t *testing.T) {
	var name QueueName
	err := name.Scan("quick_scan")
	if err != nil {
		utils.PrintTestError(t, err, nil)
	}
	if name != QuickScanQueue {
		utils.PrintTestError(t, name, QuickScanQueue)
	}
}

func TestGetQueueNames(t *testing.T) {
	names := GetQueueNames()

	expected := []QueueName{
		QuickScanQueue,
		EmailPollingQueue,
		EmailReceiptProcessingQueue,
		EmailReceiptImageCleanupQueue,
		SystemCleanUpQueue,
	}
	if len(names) != len(expected) {
		utils.PrintTestError(t, len(names), len(expected))
	}

	for i, v := range expected {
		if names[i] != v {
			utils.PrintTestError(t, names[i], v)
		}
	}
}

func TestGetDefaultQueueConfigurationMap(t *testing.T) {
	configMap := GetDefaultQueueConfigurationMap()

	expected := map[QueueName]int{
		QuickScanQueue:                4,
		EmailPollingQueue:             2,
		EmailReceiptProcessingQueue:   3,
		EmailReceiptImageCleanupQueue: 1,
		SystemCleanUpQueue:            5,
	}
	if len(configMap) != len(expected) {
		utils.PrintTestError(t, len(configMap), len(expected))
	}

	for name, priority := range expected {
		config, ok := configMap[name]
		if !ok {
			utils.PrintTestError(t, "missing queue "+string(name), "present")
			continue
		}
		if config.Name != name {
			utils.PrintTestError(t, config.Name, name)
		}
		if config.Priority != priority {
			utils.PrintTestError(t, config.Priority, priority)
		}
	}
}
