package wranglerasynq

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/hibiken/asynq"
	"receipt-wrangler/api/internal/commands"
	"receipt-wrangler/api/internal/models"
	"receipt-wrangler/api/internal/structs"
)

const testRetention = 30 * 24 * time.Hour

var (
	youngEnough = testRetention - time.Hour
	oldEnough   = testRetention + time.Hour
)

// TestClassifyTempFile walks the precedence rules exhaustively. Each row names
// the behaviour it protects, because most of them look interchangeable until one
// regresses.
func TestClassifyTempFile(t *testing.T) {
	tests := []struct {
		name               string
		states             []asynq.TaskState
		age                time.Duration
		referencesComplete bool
		expected           tempFileAction
	}{
		{
			name:               "pending keeps even when aged",
			states:             []asynq.TaskState{asynq.TaskStatePending},
			age:                oldEnough,
			referencesComplete: true,
			expected:           tempFileKeep,
		},
		{
			name:               "active keeps even when aged",
			states:             []asynq.TaskState{asynq.TaskStateActive},
			age:                oldEnough,
			referencesComplete: true,
			expected:           tempFileKeep,
		},
		{
			name:               "scheduled keeps even when aged",
			states:             []asynq.TaskState{asynq.TaskStateScheduled},
			age:                oldEnough,
			referencesComplete: true,
			expected:           tempFileKeep,
		},
		{
			// A quick scan between retries must not lose the upload it is about
			// to read again.
			name:               "retry keeps even when aged",
			states:             []asynq.TaskState{asynq.TaskStateRetry},
			age:                oldEnough,
			referencesComplete: true,
			expected:           tempFileKeep,
		},
		{
			name:               "aggregating keeps even when aged",
			states:             []asynq.TaskState{asynq.TaskStateAggregating},
			age:                oldEnough,
			referencesComplete: true,
			expected:           tempFileKeep,
		},
		{
			// The bug this sweeper replaces: archived is the rerunnable state,
			// so its file has to survive the retention window.
			name:               "archived keeps while young",
			states:             []asynq.TaskState{asynq.TaskStateArchived},
			age:                youngEnough,
			referencesComplete: true,
			expected:           tempFileKeep,
		},
		{
			name:               "archived releases once aged",
			states:             []asynq.TaskState{asynq.TaskStateArchived},
			age:                oldEnough,
			referencesComplete: true,
			expected:           tempFileRemove,
		},
		{
			// Live since enqueueOptions put asynq.Retention on the email queue:
			// a succeeded upload's attachment and OCR copy are released on the
			// next sweep instead of waiting out the whole window.
			name:               "all completed releases immediately",
			states:             []asynq.TaskState{asynq.TaskStateCompleted, asynq.TaskStateCompleted},
			age:                youngEnough,
			referencesComplete: true,
			expected:           tempFileRemove,
		},
		{
			// One attachment, two groups, one of them permanently failed. Rule 2
			// must outrank rule 3 or the successful sibling releases the file the
			// failed one still needs.
			name:               "completed plus archived keeps while young",
			states:             []asynq.TaskState{asynq.TaskStateCompleted, asynq.TaskStateArchived},
			age:                youngEnough,
			referencesComplete: true,
			expected:           tempFileKeep,
		},
		{
			name:               "completed plus archived releases once aged",
			states:             []asynq.TaskState{asynq.TaskStateCompleted, asynq.TaskStateArchived},
			age:                oldEnough,
			referencesComplete: true,
			expected:           tempFileRemove,
		},
		{
			name:               "completed plus pending keeps, actionable outranks everything",
			states:             []asynq.TaskState{asynq.TaskStateCompleted, asynq.TaskStatePending},
			age:                oldEnough,
			referencesComplete: true,
			expected:           tempFileKeep,
		},
		{
			// The grace window that covers a file written moments before the task
			// referencing it becomes visible in Redis.
			name:               "unreferenced keeps while young",
			states:             nil,
			age:                youngEnough,
			referencesComplete: true,
			expected:           tempFileKeep,
		},
		{
			// Where a successful quick scan's file would land if it had not already
			// been deleted by the handler, and where an email attachment lands once
			// its completed task has outlived completedTaskRetention — plus the
			// DebugOcr dumps, which no task ever references.
			name:               "unreferenced releases once aged",
			states:             nil,
			age:                oldEnough,
			referencesComplete: true,
			expected:           tempFileRemove,
		},
		{
			// Fail closed: an incomplete reference map cannot tell "nothing refers
			// to this" from "we could not read what refers to this".
			name:               "unreferenced keeps when the reference map is incomplete",
			states:             nil,
			age:                oldEnough,
			referencesComplete: false,
			expected:           tempFileKeep,
		},
		{
			// Rules 1-3 rest on a state we actually observed, so completeness
			// does not gate them.
			name:               "archived still releases when the reference map is incomplete",
			states:             []asynq.TaskState{asynq.TaskStateArchived},
			age:                oldEnough,
			referencesComplete: false,
			expected:           tempFileRemove,
		},
		{
			name:               "age exactly at retention releases",
			states:             []asynq.TaskState{asynq.TaskStateArchived},
			age:                testRetention,
			referencesComplete: true,
			expected:           tempFileRemove,
		},
		{
			name:               "age one tick under retention keeps",
			states:             []asynq.TaskState{asynq.TaskStateArchived},
			age:                testRetention - time.Nanosecond,
			referencesComplete: true,
			expected:           tempFileKeep,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got := classifyTempFile(test.states, test.age, testRetention, test.referencesComplete)
			if got != test.expected {
				t.Errorf("classifyTempFile() = %v, want %v", got, test.expected)
			}
		})
	}
}

func TestSweepTempDirectory_RemovesAgedOrphanKeepsYoungOne(t *testing.T) {
	directory := t.TempDir()
	aged := writeTempFile(t, directory, "aged.jpg", oldEnough)
	young := writeTempFile(t, directory, "young.jpg", youngEnough)

	removed := sweep(t, directory, nil, true)

	if removed != 1 {
		t.Errorf("expected 1 removal, got %d", removed)
	}
	assertMissing(t, aged)
	assertExists(t, young)
}

func TestSweepTempDirectory_KeepsFileWhoseTaskIsPending(t *testing.T) {
	directory := t.TempDir()
	path := writeTempFile(t, directory, "in-flight.jpg", oldEnough)

	references := map[string][]asynq.TaskState{path: {asynq.TaskStatePending}}
	removed := sweep(t, directory, references, true)

	if removed != 0 {
		t.Errorf("expected no removals, got %d", removed)
	}
	assertExists(t, path)
}

// The original attachment and its OCR copy are two independent files. The old
// cleanup keyed on the pair and could only ever release both together.
func TestSweepTempDirectory_TreatsAttachmentAndOcrCopyIndependently(t *testing.T) {
	directory := t.TempDir()
	attachment := writeTempFile(t, directory, "receipt.pdf", oldEnough)
	ocrCopy := writeTempFile(t, directory, "image-receipt.pdf", oldEnough)

	references := map[string][]asynq.TaskState{attachment: {asynq.TaskStatePending}}
	removed := sweep(t, directory, references, true)

	if removed != 1 {
		t.Errorf("expected 1 removal, got %d", removed)
	}
	assertExists(t, attachment)
	assertMissing(t, ocrCopy)
}

// A payload path is cleaned on the way into the map so it matches the path the
// directory walk builds. Without that, an un-normalized payload path would look
// like a different file and its image would be orphan-deleted.
func TestSweepTempDirectory_MatchesUnnormalizedPayloadPaths(t *testing.T) {
	directory := t.TempDir()
	path := writeTempFile(t, directory, "receipt.jpg", oldEnough)

	unnormalized := filepath.Join(directory, "nested", "..", "receipt.jpg")
	paths, err := tempPathsFromQuickScanPayload(quickScanPayload(t, unnormalized))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	references := map[string][]asynq.TaskState{paths[0]: {asynq.TaskStateRetry}}
	removed := sweep(t, directory, references, true)

	if removed != 0 {
		t.Errorf("expected no removals, got %d", removed)
	}
	assertExists(t, path)
}

func TestSweepTempDirectory_MissingDirectoryIsNotAnErrorAndIsNotCreated(t *testing.T) {
	directory := filepath.Join(t.TempDir(), "absent")

	removed, err := sweepTempDirectory(directory, nil, true, time.Now(), testRetention)
	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}
	if removed != 0 {
		t.Errorf("expected no removals, got %d", removed)
	}
	if _, err := os.Stat(directory); !os.IsNotExist(err) {
		t.Error("sweep must not create the temp directory")
	}
}

func TestSweepTempDirectory_SkipsSubdirectories(t *testing.T) {
	directory := t.TempDir()
	nested := filepath.Join(directory, "nested")
	if err := os.Mkdir(nested, os.ModePerm); err != nil {
		t.Fatalf("failed to create directory: %v", err)
	}
	if err := os.Chtimes(nested, time.Now().Add(-oldEnough), time.Now().Add(-oldEnough)); err != nil {
		t.Fatalf("failed to age directory: %v", err)
	}

	removed := sweep(t, directory, nil, true)

	if removed != 0 {
		t.Errorf("expected no removals, got %d", removed)
	}
	if _, err := os.Stat(nested); err != nil {
		t.Error("expected the subdirectory to survive")
	}
}

// DebugOcr writes these next to the real uploads and never removes them, so the
// orphan branch is the only thing that ever reclaims them.
func TestSweepTempDirectory_ReclaimsAgedDebugOcrDumps(t *testing.T) {
	directory := t.TempDir()
	text := writeTempFile(t, directory, "receipt.jpg.txt", oldEnough)
	image := writeTempFile(t, directory, "receipt.jpg.jpg", oldEnough)
	fresh := writeTempFile(t, directory, "fresh.jpg.txt", youngEnough)

	removed := sweep(t, directory, nil, true)

	if removed != 2 {
		t.Errorf("expected 2 removals, got %d", removed)
	}
	assertMissing(t, text)
	assertMissing(t, image)
	assertExists(t, fresh)
}

// A file the sweeper cannot delete is logged and skipped, never propagated as an
// error — one stuck entry must not abort an hourly janitorial run.
func TestSweepTempDirectory_SwallowsRemovalFailures(t *testing.T) {
	if os.Geteuid() == 0 {
		t.Skip("running as root, which ignores directory write permissions")
	}

	protectedDirectory := filepath.Join(t.TempDir(), "protected")
	if err := os.Mkdir(protectedDirectory, os.ModePerm); err != nil {
		t.Fatalf("failed to create directory: %v", err)
	}

	first := writeTempFile(t, protectedDirectory, "first.jpg", oldEnough)
	second := writeTempFile(t, protectedDirectory, "second.jpg", oldEnough)

	// Unlink permission is a property of the directory, so this denies removal of
	// every entry at once.
	if err := os.Chmod(protectedDirectory, 0o500); err != nil {
		t.Fatalf("failed to chmod: %v", err)
	}
	t.Cleanup(func() { _ = os.Chmod(protectedDirectory, os.ModePerm) })

	removed, err := sweepTempDirectory(protectedDirectory, nil, true, time.Now(), testRetention)
	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}
	if removed != 0 {
		t.Errorf("expected no successful removals, got %d", removed)
	}
	assertExists(t, first)
	assertExists(t, second)
}

func TestListTempFileReferences_MapsBothPayloadShapes(t *testing.T) {
	quickScanTask := taskInfo(t, "quick-scan-1", asynq.TaskStateArchived, QuickScanTaskPayload{
		TempPath: "/base/temp/scan.jpg",
	})
	emailTask := taskInfo(t, "email-1", asynq.TaskStatePending, EmailProcessTaskPayload{
		TempFilePath:    "/base/temp/attachment.pdf",
		ImageForOcrPath: "/base/temp/image-attachment.jpg",
	})

	references, complete, err := listTempFileReferences([]tempFileTaskLister{
		staticLister(map[models.QueueName][]*asynq.TaskInfo{
			models.QuickScanQueue:              {quickScanTask},
			models.EmailReceiptProcessingQueue: {emailTask},
		}),
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !complete {
		t.Error("expected the reference map to be complete")
	}

	assertStates(t, references, "/base/temp/scan.jpg", asynq.TaskStateArchived)
	assertStates(t, references, "/base/temp/attachment.pdf", asynq.TaskStatePending)
	// An email task owns two files, and both must be protected.
	assertStates(t, references, "/base/temp/image-attachment.jpg", asynq.TaskStatePending)
}

func TestListTempFileReferences_BodyOnlyEmailContributesNothing(t *testing.T) {
	bodyOnly := taskInfo(t, "email-body", asynq.TaskStateCompleted, EmailProcessTaskPayload{})

	references, complete, err := listTempFileReferences([]tempFileTaskLister{
		staticLister(map[models.QueueName][]*asynq.TaskInfo{
			models.EmailReceiptProcessingQueue: {bodyOnly},
		}),
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !complete {
		t.Error("a body-only task is not a completeness failure")
	}
	if len(references) != 0 {
		t.Errorf("expected no references, got %v", references)
	}
}

// Asynq's inspector pages at 30 by default. A queue holding more than one page
// would silently truncate, and every unseen task's file would then look
// unreferenced and be deleted once aged.
func TestListTempFileReferences_PagesBeyondTheFirstPage(t *testing.T) {
	total := tempFileListPageSize + 7
	tasks := make([]*asynq.TaskInfo, 0, total)
	for i := 0; i < total; i++ {
		tasks = append(tasks, taskInfo(t, fmt.Sprintf("task-%d", i), asynq.TaskStateArchived, QuickScanTaskPayload{
			TempPath: fmt.Sprintf("/base/temp/scan-%d.jpg", i),
		}))
	}

	references, complete, err := listTempFileReferences([]tempFileTaskLister{
		pagedLister(map[models.QueueName][]*asynq.TaskInfo{models.QuickScanQueue: tasks}),
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !complete {
		t.Error("expected the reference map to be complete")
	}
	if len(references) != total {
		t.Fatalf("expected %d referenced paths, got %d", total, len(references))
	}
}

func TestListTempFileReferences_ToleratesAQueueThatDoesNotExistYet(t *testing.T) {
	// A queue that has never been used does not exist in Redis. That must not
	// abort the sweep, and the other queue must still be scanned.
	lister := func(queue string, opts ...asynq.ListOption) ([]*asynq.TaskInfo, error) {
		if queue == string(models.QuickScanQueue) {
			return nil, fmt.Errorf("asynq: %w", asynq.ErrQueueNotFound)
		}

		return []*asynq.TaskInfo{
			taskInfo(t, "email-1", asynq.TaskStateArchived, EmailProcessTaskPayload{
				TempFilePath: "/base/temp/attachment.pdf",
			}),
		}, nil
	}

	references, complete, err := listTempFileReferences([]tempFileTaskLister{lister})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !complete {
		t.Error("a never-used queue is not a completeness failure")
	}
	assertStates(t, references, "/base/temp/attachment.pdf", asynq.TaskStateArchived)
}

// Anything other than a missing queue means we cannot trust the map at all, so
// the sweep aborts rather than deleting on partial information.
func TestListTempFileReferences_AbortsOnListingError(t *testing.T) {
	listingErr := errors.New("redis unreachable")
	lister := func(queue string, opts ...asynq.ListOption) ([]*asynq.TaskInfo, error) {
		return nil, listingErr
	}

	references, complete, err := listTempFileReferences([]tempFileTaskLister{lister})
	if !errors.Is(err, listingErr) {
		t.Fatalf("expected the listing error to propagate, got: %v", err)
	}
	if references != nil {
		t.Error("expected no references on a listing failure")
	}
	if complete {
		t.Error("expected the reference map to be marked incomplete")
	}
}

// The page-by-page extraction builds the reference map as it goes, so a failure
// PART WAY THROUGH must still discard what it had. A partial map is the dangerous
// shape: the files it did not reach look unreferenced, which is the orphan branch
// deleting them once aged.
func TestListTempFileReferences_AbortsOnAnErrorPartWayThroughPaging(t *testing.T) {
	listingErr := errors.New("redis went away mid-scan")

	firstPage := make([]*asynq.TaskInfo, 0, tempFileListPageSize)
	for i := 0; i < tempFileListPageSize; i++ {
		firstPage = append(firstPage, taskInfo(t, fmt.Sprintf("quick-scan-%d", i), asynq.TaskStateArchived, QuickScanTaskPayload{
			TempPath: fmt.Sprintf("/base/temp/scan-%d.jpg", i),
		}))
	}

	calls := 0
	lister := func(queue string, opts ...asynq.ListOption) ([]*asynq.TaskInfo, error) {
		calls++
		if calls == 1 {
			// A full page, so the loop asks for another.
			return firstPage, nil
		}
		return nil, listingErr
	}

	references, complete, err := listTempFileReferences([]tempFileTaskLister{lister})
	if !errors.Is(err, listingErr) {
		t.Fatalf("expected the listing error to propagate, got: %v", err)
	}
	if references != nil {
		t.Error("expected the partially built reference map to be discarded")
	}
	if complete {
		t.Error("expected the reference map to be marked incomplete")
	}
	if calls < 2 {
		t.Errorf("lister called %d time(s); the first page was full, so a second was expected", calls)
	}
}

// One unreadable payload used to abort every future sweep. Now it drops the
// completeness flag, which stands the orphan branch down without losing the
// references we could read.
func TestListTempFileReferences_MalformedPayloadSkipsTaskAndDropsCompleteness(t *testing.T) {
	good := taskInfo(t, "quick-scan-1", asynq.TaskStateArchived, QuickScanTaskPayload{
		TempPath: "/base/temp/scan.jpg",
	})
	malformed := &asynq.TaskInfo{ID: "quick-scan-2", State: asynq.TaskStateArchived, Payload: []byte("{not json")}

	references, complete, err := listTempFileReferences([]tempFileTaskLister{
		staticLister(map[models.QueueName][]*asynq.TaskInfo{
			models.QuickScanQueue: {good, malformed},
		}),
	})
	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}
	if complete {
		t.Error("expected the reference map to be marked incomplete")
	}
	assertStates(t, references, "/base/temp/scan.jpg", asynq.TaskStateArchived)
}

// SystemCleanUpQueue tasks carry a nil payload, which would fail every unmarshal
// and drop the completeness flag on every single run.
func TestListTempFileReferences_NeverScansTheCleanUpQueue(t *testing.T) {
	scanned := make(map[string]bool)
	lister := func(queue string, opts ...asynq.ListOption) ([]*asynq.TaskInfo, error) {
		scanned[queue] = true
		return nil, nil
	}

	if _, _, err := listTempFileReferences([]tempFileTaskLister{lister}); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if scanned[string(models.SystemCleanUpQueue)] {
		t.Error("the clean up queue must not be scanned")
	}
	if !scanned[string(models.QuickScanQueue)] || !scanned[string(models.EmailReceiptProcessingQueue)] {
		t.Error("expected both ingest queues to be scanned")
	}
}

func TestClampTempFileRetention(t *testing.T) {
	tests := []struct {
		name     string
		hours    int
		expected time.Duration
	}{
		{"unset falls back to the default", 0, defaultTempFileRetention},
		{"below the floor falls back", commands.MinTempFileRetentionHours - 1, defaultTempFileRetention},
		{"above the ceiling falls back", commands.MaxTempFileRetentionHours + 1, defaultTempFileRetention},
		{"negative falls back", -1, defaultTempFileRetention},
		{"the floor is honoured", commands.MinTempFileRetentionHours, 24 * time.Hour},
		{"the ceiling is honoured", commands.MaxTempFileRetentionHours, 8760 * time.Hour},
		{"a configured value is honoured", 48, 48 * time.Hour},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if got := clampTempFileRetention(test.hours); got != test.expected {
				t.Errorf("clampTempFileRetention(%d) = %v, want %v", test.hours, got, test.expected)
			}
		})
	}
}

func TestTempPathsFromQuickScanPayload(t *testing.T) {
	paths, err := tempPathsFromQuickScanPayload(quickScanPayload(t, "/base/temp/scan.jpg"))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(paths) != 1 || paths[0] != "/base/temp/scan.jpg" {
		t.Errorf("unexpected paths: %v", paths)
	}

	if _, err := tempPathsFromQuickScanPayload([]byte("{not json")); err == nil {
		t.Error("expected a malformed payload to error")
	}
}

func TestTempPathsFromEmailPayload(t *testing.T) {
	payload, err := json.Marshal(EmailProcessTaskPayload{
		TempFilePath:    "/base/temp/attachment.pdf",
		ImageForOcrPath: "/base/temp/image-attachment.jpg",
		Attachment:      structs.Attachment{Filename: "receipt.pdf"},
	})
	if err != nil {
		t.Fatalf("failed to marshal payload: %v", err)
	}

	paths, err := tempPathsFromEmailPayload(payload)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(paths) != 2 {
		t.Fatalf("expected both files, got %v", paths)
	}

	bodyOnly, err := tempPathsFromEmailPayload(quickScanPayload(t, ""))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(bodyOnly) != 0 {
		t.Errorf("expected no paths for a body-only task, got %v", bodyOnly)
	}
}

// --- helpers ---

func sweep(
	t *testing.T,
	directory string,
	references map[string][]asynq.TaskState,
	referencesComplete bool,
) int {
	t.Helper()

	removed, err := sweepTempDirectory(directory, references, referencesComplete, time.Now(), testRetention)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	return removed
}

func writeTempFile(t *testing.T, directory string, name string, age time.Duration) string {
	t.Helper()

	if len(name) == 0 {
		return ""
	}

	path := filepath.Join(directory, name)
	if err := os.WriteFile(path, []byte("contents"), 0o644); err != nil {
		t.Fatalf("failed to write %s: %v", path, err)
	}

	modifiedAt := time.Now().Add(-age)
	if err := os.Chtimes(path, modifiedAt, modifiedAt); err != nil {
		t.Fatalf("failed to age %s: %v", path, err)
	}

	return path
}

func assertExists(t *testing.T, path string) {
	t.Helper()

	if _, err := os.Stat(path); err != nil {
		t.Errorf("expected %s to still exist", path)
	}
}

func assertMissing(t *testing.T, path string) {
	t.Helper()

	if _, err := os.Stat(path); !os.IsNotExist(err) {
		t.Errorf("expected %s to be removed", path)
	}
}

func assertStates(t *testing.T, references map[string][]asynq.TaskState, path string, expected ...asynq.TaskState) {
	t.Helper()

	states, ok := references[path]
	if !ok {
		t.Fatalf("expected %s to be referenced, have %v", path, references)
	}
	if len(states) != len(expected) {
		t.Fatalf("expected %d states for %s, got %v", len(expected), path, states)
	}
	for i, state := range expected {
		if states[i] != state {
			t.Errorf("state %d for %s = %v, want %v", i, path, states[i], state)
		}
	}
}

func taskInfo(t *testing.T, id string, state asynq.TaskState, payload any) *asynq.TaskInfo {
	t.Helper()

	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		t.Fatalf("failed to marshal payload: %v", err)
	}

	return &asynq.TaskInfo{ID: id, State: state, Payload: payloadBytes}
}

func quickScanPayload(t *testing.T, tempPath string) []byte {
	t.Helper()

	payload, err := json.Marshal(QuickScanTaskPayload{TempPath: tempPath})
	if err != nil {
		t.Fatalf("failed to marshal payload: %v", err)
	}

	return payload
}

// staticLister returns the same single page for every call, which is all a test
// that is not about paging needs.
func staticLister(tasksByQueue map[models.QueueName][]*asynq.TaskInfo) tempFileTaskLister {
	return func(queue string, opts ...asynq.ListOption) ([]*asynq.TaskInfo, error) {
		if pageNumberFrom(opts) > 1 {
			return nil, nil
		}

		return tasksByQueue[models.QueueName(queue)], nil
	}
}

// pagedLister slices its tasks by the requested page, so the caller's paging loop
// is genuinely exercised rather than assumed.
func pagedLister(tasksByQueue map[models.QueueName][]*asynq.TaskInfo) tempFileTaskLister {
	return func(queue string, opts ...asynq.ListOption) ([]*asynq.TaskInfo, error) {
		tasks := tasksByQueue[models.QueueName(queue)]

		page := pageNumberFrom(opts)
		start := (page - 1) * tempFileListPageSize
		if start >= len(tasks) {
			return nil, nil
		}

		end := min(start+tempFileListPageSize, len(tasks))

		return tasks[start:end], nil
	}
}

// pageNumberFrom recovers the requested page. asynq's option types are
// unexported, but Page(n) returns a comparable value, so matching against it is
// the only way to assert from outside the package that paging really happens.
func pageNumberFrom(opts []asynq.ListOption) int {
	for _, opt := range opts {
		for page := 1; page <= tempFileMaxListPages; page++ {
			if opt == asynq.Page(page) {
				return page
			}
		}
	}

	return 1
}
