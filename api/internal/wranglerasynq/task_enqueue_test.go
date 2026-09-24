package wranglerasynq

import (
	"testing"
	"time"

	"github.com/hibiken/asynq"
	"receipt-wrangler/api/internal/models"
)

// asynq.Option exposes Type() and Value(), so which queues retain can be pinned
// without a Redis instance — and which queues retain is the whole decision here.
func TestEnqueueOptions(t *testing.T) {
	tests := []struct {
		name              string
		queue             models.QueueName
		expectedRetention time.Duration
	}{
		{
			// The only ingest queue whose success path leaves a temp file behind:
			// one attachment fans out to a sibling task per groupSettingsId, so
			// the handler cannot delete it and the sweep has to decide.
			name:              "email receipt processing retains its completed tasks",
			queue:             models.EmailReceiptProcessingQueue,
			expectedRetention: completedTaskRetention,
		},
		{
			// ReceiptService.QuickScan removes its own upload on success and is
			// 1:1 task-to-file, so retention here would pin a payload naming a
			// path that is already gone.
			name:  "quick scan does not",
			queue: models.QuickScanQueue,
		},
		{
			name:  "email polling does not, it owns no temp files",
			queue: models.EmailPollingQueue,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			options := enqueueOptions(test.queue)

			retention, retentionSet := time.Duration(0), false
			maxRetry, maxRetrySet := 0, false
			queue, queueSet := "", false

			for _, option := range options {
				switch option.Type() {
				case asynq.RetentionOpt:
					retention, retentionSet = option.Value().(time.Duration), true
				case asynq.MaxRetryOpt:
					maxRetry, maxRetrySet = option.Value().(int), true
				case asynq.QueueOpt:
					queue, queueSet = option.Value().(string), true
				default:
					t.Errorf("unexpected option %v on %s", option, test.queue)
				}
			}

			if !maxRetrySet || maxRetry != 3 {
				t.Errorf("MaxRetry = %d (set %v), want 3", maxRetry, maxRetrySet)
			}
			if !queueSet || queue != string(test.queue) {
				t.Errorf("Queue = %q (set %v), want %q", queue, queueSet, test.queue)
			}

			if test.expectedRetention == 0 {
				if retentionSet {
					t.Errorf("%s carries Retention(%v), want none", test.queue, retention)
				}
				return
			}

			if !retentionSet {
				t.Fatalf("%s carries no Retention, want %v", test.queue, test.expectedRetention)
			}
			if retention != test.expectedRetention {
				t.Errorf("Retention = %v, want %v", retention, test.expectedRetention)
			}
		})
	}
}

// Retention shorter than the sweep's own interval would never be observed, which
// would make the whole option a no-op.
func TestCompletedTaskRetentionOutlivesTheSweepInterval(t *testing.T) {
	if completedTaskRetention <= time.Hour {
		t.Errorf("completedTaskRetention = %v, must exceed the @every 1h sweep interval", completedTaskRetention)
	}
}
