package wranglerasynq

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/hibiken/asynq"
	"receipt-wrangler/api/internal/commands"
	"receipt-wrangler/api/internal/logging"
	"receipt-wrangler/api/internal/models"
	"receipt-wrangler/api/internal/repositories"
)

// The temp/ sweeper reclaims files the ingest pipelines leave behind, without
// ever deleting one that something can still act on.
//
// It replaces a cleanup that was exactly inverted: it released a file once every
// referencing task was Completed OR Archived, and Archived is precisely the state
// that makes an activity rerunnable. Worse, no queue here sets asynq.Retention, so
// a successfully-processed task is dropped from Redis immediately and never enters
// the completed set — meaning the old Completed branch could never fire and the
// only files it ever deleted were the ones it had to keep.
//
// Two invariants keep the orphan branch safe, because the failure mode inverted
// with it. Previously a task the scan missed meant "we fail to delete" (harmless);
// now it means the file looks unreferenced and is deleted once aged. So:
//
//  1. every inspector list call pages (asynq defaults to 30 per page), and
//  2. any listing or payload error stands the orphan branch down for that run.
//
// Rules 1-3 below rest on positive evidence — a state we actually observed — so
// they stay live even when the reference map is known-incomplete.
const (
	// tempFileListPageSize is the page size for every inspector listing. Asynq's
	// default is 30; a queue with more tasks than that would silently truncate.
	tempFileListPageSize = 100

	// tempFileMaxListPages bounds the paging loop. A well-behaved lister ends it
	// with a short page; this is only so a misbehaving one cannot spin forever
	// inside a background job.
	tempFileMaxListPages = 1000
)

const defaultTempFileRetention = commands.DefaultTempFileRetentionHours * time.Hour

// tempFileAction is what the sweeper decided to do with one file in temp/.
type tempFileAction int

const (
	tempFileKeep tempFileAction = iota
	tempFileRemove
)

// tempFileTaskLister lists one state of one queue. It is a named type purely so
// the reference scan can be driven without Redis in tests; in production every
// value is an asynq.Inspector method.
type tempFileTaskLister func(queue string, opts ...asynq.ListOption) ([]*asynq.TaskInfo, error)

// tempFileQueue pairs an ingest queue with the extractor for its payload shape.
type tempFileQueue struct {
	name    models.QueueName
	extract func(payload []byte) ([]string, error)
}

// tempFileQueues are the only queues whose payloads reference temp/ files.
// SystemCleanUpQueue is deliberately absent: its tasks carry a nil payload, which
// would fail every unmarshal and drop the completeness flag on every run.
func tempFileQueues() []tempFileQueue {
	return []tempFileQueue{
		{name: models.QuickScanQueue, extract: tempPathsFromQuickScanPayload},
		{name: models.EmailReceiptProcessingQueue, extract: tempPathsFromEmailPayload},
	}
}

func HandleTempFileCleanUpTask(context context.Context, task *asynq.Task) error {
	inspector, err := GetAsynqInspector()
	if err != nil {
		return err
	}
	defer inspector.Close()

	references, referencesComplete, err := listTempFileReferences(inspectorTaskListers(inspector))
	if err != nil {
		return err
	}

	fileRepository := repositories.NewFileRepository(nil)
	removed, err := sweepTempDirectory(
		fileRepository.GetTempDirectoryPath(),
		references,
		referencesComplete,
		time.Now(),
		tempFileRetention(),
	)
	if err != nil {
		return err
	}

	if removed > 0 {
		logging.LogStd(logging.LOG_LEVEL_INFO, fmt.Sprintf("Temp file cleanup removed %d file(s)", removed))
	}

	return nil
}

func inspectorTaskListers(inspector *asynq.Inspector) []tempFileTaskLister {
	return []tempFileTaskLister{
		inspector.ListPendingTasks,
		inspector.ListActiveTasks,
		inspector.ListScheduledTasks,
		inspector.ListRetryTasks,
		inspector.ListArchivedTasks,
		inspector.ListCompletedTasks,
	}
}

// listTempFileReferences maps every temp path an ingest task still refers to onto
// the states of the tasks referring to it.
//
// The second return value reports whether the map is known-complete. It goes
// false when a payload cannot be read: that task's paths are missing from the map,
// so treating their files as unreferenced would delete them. Only the orphan
// branch consults it.
//
// A listing error is fatal to the whole sweep rather than partial, so a Redis
// hiccup can never be mistaken for "nothing references these files".
func listTempFileReferences(listers []tempFileTaskLister) (map[string][]asynq.TaskState, bool, error) {
	references := make(map[string][]asynq.TaskState)
	referencesComplete := true

	for _, queue := range tempFileQueues() {
		for _, list := range listers {
			tasks, err := listAllTasks(list, string(queue.name))
			if err != nil {
				return nil, false, err
			}

			for _, task := range tasks {
				paths, err := queue.extract(task.Payload)
				if err != nil {
					logging.LogStd(
						logging.LOG_LEVEL_ERROR,
						"Could not read temp file references from task ", task.ID, ": ", err.Error(),
					)
					referencesComplete = false
					continue
				}

				for _, path := range paths {
					references[path] = append(references[path], task.State)
				}
			}
		}
	}

	return references, referencesComplete, nil
}

// listAllTasks pages one listing to exhaustion. A queue that has never been used
// does not exist in Redis yet, which is not an error worth aborting a sweep over.
func listAllTasks(list tempFileTaskLister, queue string) ([]*asynq.TaskInfo, error) {
	allTasks := make([]*asynq.TaskInfo, 0)

	for page := 1; page <= tempFileMaxListPages; page++ {
		tasks, err := list(queue, asynq.PageSize(tempFileListPageSize), asynq.Page(page))
		if err != nil {
			if errors.Is(err, asynq.ErrQueueNotFound) {
				return allTasks, nil
			}
			return nil, err
		}

		allTasks = append(allTasks, tasks...)

		if len(tasks) < tempFileListPageSize {
			return allTasks, nil
		}
	}

	return allTasks, nil
}

// classifyTempFile decides the fate of a single file. Pure — no Redis, no
// filesystem — so the precedence rules can be table-tested exhaustively.
//
// Precedence, and the order matters:
//
//  1. any referencing task is still actionable  -> keep
//  2. else any referencing task is archived     -> keep until aged
//  3. else every referencing task succeeded     -> remove now
//  4. unreferenced                              -> keep until aged, and only
//     when the reference map is complete
//
// Rule 2 must be tested before rule 3. One email attachment fans out to a sibling
// task per group, so a file can be referenced by both a succeeded task and a
// permanently failed one; releasing it because something succeeded is the exact
// bug this sweeper replaces.
func classifyTempFile(
	states []asynq.TaskState,
	age time.Duration,
	retention time.Duration,
	referencesComplete bool,
) tempFileAction {
	if len(states) == 0 {
		if !referencesComplete {
			return tempFileKeep
		}

		return agedTempFileAction(age, retention)
	}

	hasArchived := false
	for _, state := range states {
		switch state {
		case asynq.TaskStateArchived:
			hasArchived = true
		case asynq.TaskStateCompleted:
			// Terminal success. On its own it releases the file, but only if
			// nothing else below outranks it.
		default:
			// Pending, active, scheduled, retry, aggregating — and anything a
			// future asynq adds. Still actionable, so the file stays.
			return tempFileKeep
		}
	}

	if hasArchived {
		return agedTempFileAction(age, retention)
	}

	return tempFileRemove
}

// agedTempFileAction releases a file once it is at least retention old. The
// window doubles as the grace period for the race where a file is written just
// before the task referencing it becomes visible in Redis.
func agedTempFileAction(age time.Duration, retention time.Duration) tempFileAction {
	if age >= retention {
		return tempFileRemove
	}

	return tempFileKeep
}

// sweepTempDirectory applies the classification to every file in directory and
// reports how many it removed.
//
// A missing directory is not an error and is deliberately not created: temp/ is
// made on demand by the writers, and utils.MakeDirectory is os.Mkdir rather than
// os.MkdirAll. Per-file failures are logged and skipped so one unreadable or
// undeletable entry cannot strand the rest.
func sweepTempDirectory(
	directory string,
	references map[string][]asynq.TaskState,
	referencesComplete bool,
	now time.Time,
	retention time.Duration,
) (int, error) {
	entries, err := os.ReadDir(directory)
	if err != nil {
		if os.IsNotExist(err) {
			return 0, nil
		}

		return 0, err
	}

	removed := 0
	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}

		path := filepath.Join(directory, entry.Name())

		info, err := entry.Info()
		if err != nil {
			logging.LogStd(logging.LOG_LEVEL_ERROR, "Could not stat temp file ", path, ": ", err.Error())
			continue
		}

		action := classifyTempFile(
			references[path],
			now.Sub(info.ModTime()),
			retention,
			referencesComplete,
		)
		if action != tempFileRemove {
			continue
		}

		if err := os.Remove(path); err != nil {
			logging.LogStd(logging.LOG_LEVEL_ERROR, "Could not remove temp file ", path, ": ", err.Error())
			continue
		}

		removed++
	}

	return removed, nil
}

func tempPathsFromQuickScanPayload(payload []byte) ([]string, error) {
	var parsedPayload QuickScanTaskPayload
	if err := json.Unmarshal(payload, &parsedPayload); err != nil {
		return nil, err
	}

	return cleanedTempPaths(parsedPayload.TempPath), nil
}

// tempPathsFromEmailPayload yields both files an email task owns: the original
// attachment and the converted copy fed to OCR. A body-only task has neither,
// which is not an error.
func tempPathsFromEmailPayload(payload []byte) ([]string, error) {
	var parsedPayload EmailProcessTaskPayload
	if err := json.Unmarshal(payload, &parsedPayload); err != nil {
		return nil, err
	}

	return cleanedTempPaths(parsedPayload.TempFilePath, parsedPayload.ImageForOcrPath), nil
}

// cleanedTempPaths drops empty entries and normalizes the rest, so a path read
// out of a payload compares equal to the one filepath.Join produces while
// walking the directory.
func cleanedTempPaths(paths ...string) []string {
	cleaned := make([]string, 0, len(paths))

	for _, path := range paths {
		if len(path) == 0 {
			continue
		}

		cleaned = append(cleaned, filepath.Clean(path))
	}

	return cleaned
}

// tempFileRetention resolves the configured retention window, falling back to the
// built-in default when it is unset, out of range, or unreadable.
//
// Same clamp-with-fallback shape as repositories.pdfRasterizationDpi: the clamp,
// not the command validator, is the real safety net, because it also covers a
// value that predates the bounds or a settings row that cannot be read at all.
func tempFileRetention() time.Duration {
	systemSettingsRepository := repositories.NewSystemSettingsRepository(nil)

	systemSettings, err := systemSettingsRepository.GetSystemSettings()
	if err != nil {
		logging.LogStd(
			logging.LOG_LEVEL_ERROR,
			"Could not read temp file retention, using default: "+err.Error(),
		)

		return defaultTempFileRetention
	}

	return clampTempFileRetention(systemSettings.TempFileRetentionHours)
}

func clampTempFileRetention(hours int) time.Duration {
	if hours < commands.MinTempFileRetentionHours || hours > commands.MaxTempFileRetentionHours {
		return defaultTempFileRetention
	}

	return time.Duration(hours) * time.Hour
}
