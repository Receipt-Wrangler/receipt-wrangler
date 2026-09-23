package wranglerasynq

import (
	"time"

	"github.com/hibiken/asynq"
	"receipt-wrangler/api/internal/models"
	"receipt-wrangler/api/internal/repositories"
)

// completedTaskRetention keeps a successfully-processed email task in Redis long
// enough for the hourly temp sweep to observe its Completed state and release the
// attachment and its OCR copy. Without it asynq drops the task the instant it
// succeeds, so the sweep only ever sees those files as unreferenced and they wait
// out the whole user-configured window (720h by default) instead.
//
// It has to exceed the sweep's interval or it is a no-op. It is a handful of
// intervals rather than a day because EmailProcessTaskPayload carries the email
// body TWICE — Metadata.Body and Metadata.BodyHtml — and email.go copies that
// metadata per groupSettingsId, so one message consumed by N groups with M
// attachments retains N*M copies of a body that is routinely 50-500 KB.
// Overshooting degrades benignly: the files simply fall back to waiting out the
// configured window. It also survives ~5h of downtime, since the schedule is
// re-registered on boot and the first post-restart sweep is an hour out.
const completedTaskRetention = 6 * time.Hour

// enqueueOptions are the asynq options a task enqueued onto this queue carries.
//
// Only EmailReceiptProcessingQueue retains its completed tasks, because it is the
// only queue whose success path leaves a temp file behind for the sweep to
// classify. ReceiptService.QuickScan removes its own upload as the last step of a
// successful scan and is 1:1 task-to-file, so retention there would only pin a
// payload naming a path that no longer exists; email cannot do the same because
// one attachment fans out to a sibling task per groupSettingsId and the first to
// finish would delete a file the others still need. EmailPollingQueue owns no
// temp files at all.
func enqueueOptions(queue models.QueueName) []asynq.Option {
	options := []asynq.Option{asynq.MaxRetry(3), asynq.Queue(string(queue))}

	if queue == models.EmailReceiptProcessingQueue {
		options = append(options, asynq.Retention(completedTaskRetention))
	}

	return options
}

func EnqueueTask(task *asynq.Task, queue models.QueueName) (*asynq.TaskInfo, error) {
	client := repositories.GetAsynqClient()
	return client.Enqueue(task, enqueueOptions(queue)...)
}

func RegisterTask(cronspec string, task *asynq.Task, queue models.QueueName, maxRetry int) (string, error) {
	return scheduler.Register(cronspec, task, asynq.MaxRetry(maxRetry), asynq.Queue(string(queue)))
}
