package wranglerasynq

import (
	"github.com/hibiken/asynq"
	"receipt-wrangler/api/internal/logging"
	"receipt-wrangler/api/internal/models"
	"slices"
)

func StartSystemCleanUpTasks() error {
	inspector, err := GetAsynqInspector()
	if err != nil {
		return err
	}
	defer inspector.Close()

	cleanUpQueue := models.SystemCleanUpQueue

	inspector.DeleteAllScheduledTasks(string(cleanUpQueue))
	retireEmailReceiptImageCleanupQueue(inspector)

	refreshTokenTask := asynq.NewTask(RefreshTokenCleanUp, nil)
	_, err = RegisterTask("@every 24h", refreshTokenTask, cleanUpQueue, 0)
	if err != nil {
		return err
	}

	oidcSessionTask := asynq.NewTask(OidcSessionCleanUp, nil)
	_, err = RegisterTask("@every 24h", oidcSessionTask, cleanUpQueue, 0)
	if err != nil {
		return err
	}

	// The temp sweep lives here rather than in StartEmailPolling, where it used
	// to be, for two reasons. main.go only calls StartEmailPolling when email
	// polling is configured, so an install that never set it up ran no temp
	// cleanup at all — including for quick-scan uploads, which have nothing to do
	// with email. And StartEmailPolling runs again on every polling-interval
	// change, while scheduler.Register adds an entry without removing the
	// previous one, so the sweep used to accumulate duplicate crons.
	tempFileCleanUpTask := asynq.NewTask(TempFileCleanUp, nil)
	_, err = RegisterTask("@every 1h", tempFileCleanUpTask, cleanUpQueue, 0)

	return err
}

// retiredQueueInspector is the part of *asynq.Inspector the drain uses, so a test
// can stand in for Redis.
type retiredQueueInspector interface {
	Queues() ([]string, error)
	DeleteAllScheduledTasks(queue string) (int, error)
	DeleteAllPendingTasks(queue string) (int, error)
}

// retireEmailReceiptImageCleanupQueue drains what the pre-upgrade cleanup cron
// left on its old queue. Scheduled and pending are both needed: DeleteAllScheduled
// only touches tasks waiting on a future time, while anything the cron had already
// fired sits in pending.
//
// The queue itself stays declared in models.QueueName — persisted
// TaskQueueConfigurations reference it, QueueName.Value() whitelists it, and
// UpsertSystemSettingsCommand.Validate requires a configuration for every name.
func retireEmailReceiptImageCleanupQueue(inspector retiredQueueInspector) {
	retiredQueue := string(models.EmailReceiptImageCleanupQueue)

	// asynq records a queue in Redis only when a task is first enqueued to it, and
	// nothing enqueues here any more. On a Redis that never ran the old cron — a
	// fresh install, or a Redis that keeps no data across restarts — the queue does
	// not exist, there is nothing to drain, and both deletes would fail with
	// NOT_FOUND on every call. They report that through asynq's internal error
	// type rather than asynq.ErrQueueNotFound, so it is checked for up front.
	queues, err := inspector.Queues()
	if err != nil {
		logging.LogStd(logging.LOG_LEVEL_INFO, err.Error())
		return
	}
	if !slices.Contains(queues, retiredQueue) {
		return
	}

	if _, err := inspector.DeleteAllScheduledTasks(retiredQueue); err != nil {
		logging.LogStd(logging.LOG_LEVEL_INFO, err.Error())
	}

	if _, err := inspector.DeleteAllPendingTasks(retiredQueue); err != nil {
		logging.LogStd(logging.LOG_LEVEL_INFO, err.Error())
	}
}
