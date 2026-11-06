package wranglerasynq

import (
	"errors"
	"slices"
	"testing"

	"receipt-wrangler/api/internal/models"
)

type fakeRetiredQueueInspector struct {
	queues    []string
	queuesErr error
	deleted   []string
}

func (f *fakeRetiredQueueInspector) Queues() ([]string, error) {
	return f.queues, f.queuesErr
}

func (f *fakeRetiredQueueInspector) DeleteAllScheduledTasks(queue string) (int, error) {
	f.deleted = append(f.deleted, "scheduled:"+queue)
	return 0, nil
}

func (f *fakeRetiredQueueInspector) DeleteAllPendingTasks(queue string) (int, error) {
	f.deleted = append(f.deleted, "pending:"+queue)
	return 0, nil
}

func TestRetireEmailReceiptImageCleanupQueue_DrainsScheduledAndPending(t *testing.T) {
	retiredQueue := string(models.EmailReceiptImageCleanupQueue)
	inspector := &fakeRetiredQueueInspector{
		queues: []string{string(models.QuickScanQueue), retiredQueue},
	}

	retireEmailReceiptImageCleanupQueue(inspector)

	expected := []string{"scheduled:" + retiredQueue, "pending:" + retiredQueue}
	if !slices.Equal(inspector.deleted, expected) {
		t.Fatalf("expected %v, got %v", expected, inspector.deleted)
	}
}

// Nothing enqueues to the retired queue any more, so on a Redis that never ran
// the old cron it does not exist. asynq's delete calls fail with NOT_FOUND for a
// missing queue, which used to log on every startup and settings save.
func TestRetireEmailReceiptImageCleanupQueue_SkipsAQueueThatDoesNotExist(t *testing.T) {
	inspector := &fakeRetiredQueueInspector{
		queues: []string{string(models.QuickScanQueue), string(models.SystemCleanUpQueue)},
	}

	retireEmailReceiptImageCleanupQueue(inspector)

	if len(inspector.deleted) != 0 {
		t.Fatalf("expected no deletes, got %v", inspector.deleted)
	}
}

func TestRetireEmailReceiptImageCleanupQueue_SkipsWhenQueuesCannotBeListed(t *testing.T) {
	inspector := &fakeRetiredQueueInspector{queuesErr: errors.New("redis unreachable")}

	retireEmailReceiptImageCleanupQueue(inspector)

	if len(inspector.deleted) != 0 {
		t.Fatalf("expected no deletes, got %v", inspector.deleted)
	}
}
