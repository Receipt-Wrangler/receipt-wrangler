package wranglerasynq

import (
	"errors"
	"receipt-wrangler/api/internal/models"
)

func SystemTaskToQueueName(taskType models.SystemTaskType) (string, error) {
	if string(taskType) == string(models.QUICK_SCAN) {
		return string(models.QuickScanQueue), nil
	}

	if string(taskType) == string(models.EMAIL_UPLOAD) {
		return string(models.EmailReceiptProcessingQueue), nil
	}

	return "", errors.New("unsupported task type")
}
