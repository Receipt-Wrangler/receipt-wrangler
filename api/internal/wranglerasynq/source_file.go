package wranglerasynq

import (
	"encoding/json"
	"errors"

	"github.com/hibiken/asynq"
	"receipt-wrangler/api/internal/logging"
	"receipt-wrangler/api/internal/models"
	"receipt-wrangler/api/internal/repositories"
	"receipt-wrangler/api/internal/structs"
	"receipt-wrangler/api/internal/utils"
)

var (
	// ErrSourceFileUnsupportedTask is returned for a system task type that never
	// owns an upload. It is a malformed request rather than a missing file.
	ErrSourceFileUnsupportedTask = errors.New("system task type has no source file")

	// ErrSourceFileUnavailable is returned when the task legitimately has no
	// upload, or the upload is gone — aged out of Redis, or swept once past the
	// retention window.
	ErrSourceFileUnavailable = errors.New("source file is no longer available")
)

// SystemTaskSourceFile locates the upload behind a quick scan or email upload,
// so a failed activity can hand the user their file back.
type SystemTaskSourceFile struct {
	// Path is the original upload, exactly as it arrived. The download endpoint
	// serves this and nothing else.
	Path string
	// PreviewPath is a server-made copy already converted for OCR, when the
	// pipeline produced one. The preview endpoint prefers it because converting
	// Path instead can mean rasterizing a multi-page PDF at 300 DPI inside a
	// request. Empty for quick scan, which has no such copy.
	PreviewPath string
	// FileName is what the upload was called, used for Content-Disposition.
	FileName string
	// GroupId is the group the activity belongs to, for the permission gate.
	GroupId string
}

// taskSourceFiles are the temp files one task's payload refers to.
type taskSourceFiles struct {
	// Primary is the original upload — what the user recognises. Empty when the
	// task legitimately has none, e.g. an email with a body but no attachment.
	Primary string
	// Preview is the converted copy, when the pipeline made one.
	Preview string
}

// expectsSourceFile reports whether this task should have an upload at all.
// Distinguishing "has none by design" from "had one and it is gone" is what lets
// a body-only email stay rerunnable.
func (files taskSourceFiles) expectsSourceFile() bool {
	return len(files.Primary) > 0
}

// rerunPaths are every file a rerun reads. Email needs both: HandleEmailProcessTask
// reads Primary for the FileData it persists and hands Preview to the OCR/vision
// pipeline, so one missing file is enough to make the rerun fail.
func (files taskSourceFiles) rerunPaths() []string {
	return cleanedTempPaths(files.Primary, files.Preview)
}

// activityFlags are the two computed booleans an activity carries.
type activityFlags struct {
	canBeRestarted bool
	hasSourceFile  bool
}

// taskInfoLookup fetches one task by id. Named so the flag resolution can be
// tested without Redis; in production it is always asynq.Inspector.GetTaskInfo.
type taskInfoLookup func(queue string, id string) (*asynq.TaskInfo, error)

// SetActivityFlags fills CanBeRestarted and HasSourceFile on each activity.
//
// It looks each task up by id rather than listing the archived set, as this used
// to. The listing was capped at asynq's default 30 per page, so the flag was
// simply wrong past 30 archived tasks, and it needed a GetSystemTaskById per row
// on top. A direct lookup is O(1), has no page cap, and reports the task's state
// whatever that state is — which HasSourceFile needs, since it is not restricted
// to archived tasks.
//
// Nothing here fails the request. The activity feed rendering is worth more than
// the two flags, so an unreachable Redis or an unreadable payload leaves them
// false, exactly as the previous implementation tolerated a fresh Redis instance.
func SetActivityFlags(activities *[]structs.Activity) error {
	if len(*activities) == 0 {
		return nil
	}

	inspector, err := GetAsynqInspector()
	if err != nil {
		return nil
	}
	defer inspector.Close()

	applyActivityFlags(activities, memoizeTaskInfoLookup(inspector.GetTaskInfo))

	return nil
}

// applyActivityFlags is the loop of SetActivityFlags with the inspector lifted
// out, so the rules can be tested without Redis.
func applyActivityFlags(activities *[]structs.Activity, lookup taskInfoLookup) {
	for i := range *activities {
		activity := &(*activities)[i]

		flags, err := resolveActivityFlags(lookup, activity.Type, activity.AsynqTaskId)
		if err != nil {
			// A hard lookup failure means Redis is unhappy, not that this one row
			// is odd. Stop rather than paying a dial timeout per remaining row.
			logging.LogStd(logging.LOG_LEVEL_ERROR, "Could not resolve activity flags: ", err.Error())
			return
		}

		activity.CanBeRestarted = flags.canBeRestarted
		activity.HasSourceFile = flags.hasSourceFile
	}
}

// SetSystemTaskHasSourceFile fills HasSourceFile on the system tasks of one page.
//
// canReadGroup decides per group whether the caller may reach the file. The
// system-task table is app-scoped, so an administrator can see rows from groups
// they are not a member of; the endpoints gate on the group permission, and
// without this the table would offer buttons that 403.
//
// It deliberately does NOT descend into ChildSystemTasks. Every EMAIL_UPLOAD task
// is chained under an EMAIL_READ parent, so "top level" cannot be tested with
// AssociatedSystemTaskId == nil — that reads as false for exactly the rows this
// feature exists for. Only the rows the page itself returned are hydrated.
func SetSystemTaskHasSourceFile(systemTasks []models.SystemTask, canReadGroup func(groupId uint) bool) error {
	if len(systemTasks) == 0 {
		return nil
	}

	inspector, err := GetAsynqInspector()
	if err != nil {
		return nil
	}
	defer inspector.Close()

	hydrateSystemTaskSourceFiles(systemTasks, memoizeTaskInfoLookup(inspector.GetTaskInfo), canReadGroup)

	return nil
}

// hydrateSystemTaskSourceFiles is the loop of SetSystemTaskHasSourceFile with the
// inspector lifted out, so the rules can be tested without Redis.
func hydrateSystemTaskSourceFiles(
	systemTasks []models.SystemTask,
	lookup taskInfoLookup,
	canReadGroup func(groupId uint) bool,
) {
	for i := range systemTasks {
		systemTask := &systemTasks[i]

		if systemTask.GroupId == nil || !canReadGroup(*systemTask.GroupId) {
			continue
		}

		flags, err := resolveActivityFlags(lookup, systemTask.Type, systemTask.AsynqTaskId)
		if err != nil {
			logging.LogStd(logging.LOG_LEVEL_ERROR, "Could not resolve system task flags: ", err.Error())
			return
		}

		systemTask.HasSourceFile = flags.hasSourceFile
	}
}

// resolveActivityFlags computes both flags for one task.
//
// CanBeRestarted stays pinned to Archived because Inspector.RunTask refuses a task
// in any other state — only the lookup became state-agnostic, not the rule.
// HasSourceFile is deliberately not gated on Archived: asynq backs its retries off
// exponentially, so a task takes minutes to archive, and the user should not watch
// an activity sit at FAILED with no way to retrieve their image.
func resolveActivityFlags(
	lookup taskInfoLookup,
	taskType models.SystemTaskType,
	asynqTaskId string,
) (activityFlags, error) {
	// RECEIPT_UPLOADED / RECEIPT_UPDATED rows, and anything predating the column,
	// have no task to look up. GetTaskInfo would reject the empty id.
	if len(asynqTaskId) == 0 {
		return activityFlags{}, nil
	}

	queueName, err := SystemTaskToQueueName(taskType)
	if err != nil {
		// Not a type that owns an asynq task; neither flag applies.
		return activityFlags{}, nil
	}

	taskInfo, err := lookup(queueName, asynqTaskId)
	if err != nil {
		if errors.Is(err, asynq.ErrTaskNotFound) || errors.Is(err, asynq.ErrQueueNotFound) {
			// The task aged out of Redis. Nothing to rerun, nothing to serve.
			return activityFlags{}, nil
		}

		return activityFlags{}, err
	}

	files, err := taskSourceFilesFromPayload(taskType, taskInfo.Payload)
	if err != nil {
		logging.LogStd(logging.LOG_LEVEL_ERROR, "Could not read source files from task ", asynqTaskId, ": ", err.Error())
		return activityFlags{}, nil
	}

	return activityFlags{
		canBeRestarted: taskInfo.State == asynq.TaskStateArchived && allFilesPresent(files.rerunPaths()),
		hasSourceFile:  files.expectsSourceFile() && utils.FileExists(files.Primary),
	}, nil
}

// RerunSourceFilesPresent reports whether every file a rerun of this task would
// read is still on disk. A task type that reads none — or an unreadable payload,
// which the rerun itself will surface — counts as present, so this only ever
// blocks the case it can actually prove.
func RerunSourceFilesPresent(taskType models.SystemTaskType, payload []byte) bool {
	files, err := taskSourceFilesFromPayload(taskType, payload)
	if err != nil {
		return true
	}

	return allFilesPresent(files.rerunPaths())
}

func allFilesPresent(paths []string) bool {
	for _, path := range paths {
		if !utils.FileExists(path) {
			return false
		}
	}

	return true
}

// memoizeTaskInfoLookup caches within a single call. A primary and fallback
// processing setting produce two system tasks sharing one asynq task id, and the
// system-task table shows parent and child rows, so the same id is looked up
// several times per page.
func memoizeTaskInfoLookup(lookup taskInfoLookup) taskInfoLookup {
	type result struct {
		taskInfo *asynq.TaskInfo
		err      error
	}

	cache := make(map[string]result)

	return func(queue string, id string) (*asynq.TaskInfo, error) {
		key := queue + "|" + id

		if cached, ok := cache[key]; ok {
			return cached.taskInfo, cached.err
		}

		taskInfo, err := lookup(queue, id)
		cache[key] = result{taskInfo: taskInfo, err: err}

		return taskInfo, err
	}
}

func taskSourceFilesFromPayload(taskType models.SystemTaskType, payload []byte) (taskSourceFiles, error) {
	switch taskType {
	case models.QUICK_SCAN:
		var parsedPayload QuickScanTaskPayload
		if err := json.Unmarshal(payload, &parsedPayload); err != nil {
			return taskSourceFiles{}, err
		}

		return taskSourceFiles{Primary: parsedPayload.TempPath}, nil

	case models.EMAIL_UPLOAD:
		var parsedPayload EmailProcessTaskPayload
		if err := json.Unmarshal(payload, &parsedPayload); err != nil {
			return taskSourceFiles{}, err
		}

		return taskSourceFiles{
			Primary: parsedPayload.TempFilePath,
			Preview: parsedPayload.ImageForOcrPath,
		}, nil
	}

	return taskSourceFiles{}, nil
}

// ResolveSystemTaskSourceFile locates the upload behind a system task, ready to
// preview or download.
//
// The path comes out of a Redis payload, so it is confined to temp/ before the
// caller is allowed near it.
func ResolveSystemTaskSourceFile(systemTask models.SystemTask) (SystemTaskSourceFile, error) {
	queueName, err := SystemTaskToQueueName(systemTask.Type)
	if err != nil {
		return SystemTaskSourceFile{}, ErrSourceFileUnsupportedTask
	}

	if len(systemTask.AsynqTaskId) == 0 {
		return SystemTaskSourceFile{}, ErrSourceFileUnavailable
	}

	inspector, err := GetAsynqInspector()
	if err != nil {
		return SystemTaskSourceFile{}, err
	}
	defer inspector.Close()

	taskInfo, err := inspector.GetTaskInfo(queueName, systemTask.AsynqTaskId)
	if err != nil {
		if errors.Is(err, asynq.ErrTaskNotFound) || errors.Is(err, asynq.ErrQueueNotFound) {
			return SystemTaskSourceFile{}, ErrSourceFileUnavailable
		}

		return SystemTaskSourceFile{}, err
	}

	var payload RerunTaskPayload
	if err := json.Unmarshal(taskInfo.Payload, &payload); err != nil {
		return SystemTaskSourceFile{}, err
	}

	groupId, err := ResolveActivityGroupId(payload)
	if err != nil {
		return SystemTaskSourceFile{}, err
	}

	sourceFile := SystemTaskSourceFile{GroupId: groupId}

	switch systemTask.Type {
	case models.QUICK_SCAN:
		sourceFile.Path = payload.TempPath
		sourceFile.FileName = payload.OriginalFileName
	case models.EMAIL_UPLOAD:
		sourceFile.Path = payload.TempFilePath
		sourceFile.PreviewPath = payload.ImageForOcrPath
		sourceFile.FileName = payload.Attachment.Filename
	}

	if len(sourceFile.Path) == 0 {
		return SystemTaskSourceFile{}, ErrSourceFileUnavailable
	}

	fileRepository := repositories.NewFileRepository(nil)
	if err := fileRepository.AssertWithinTempDirectory(sourceFile.Path); err != nil {
		return SystemTaskSourceFile{}, err
	}

	if len(sourceFile.PreviewPath) > 0 {
		if err := fileRepository.AssertWithinTempDirectory(sourceFile.PreviewPath); err != nil {
			return SystemTaskSourceFile{}, err
		}
	}

	if !utils.FileExists(sourceFile.Path) {
		return SystemTaskSourceFile{}, ErrSourceFileUnavailable
	}

	return sourceFile, nil
}

// ResolveActivityGroupId resolves the group an activity belongs to. Quick scan
// carries the group id directly; email carries the group settings it was polled
// for, whose group has to be looked up.
func ResolveActivityGroupId(payload RerunTaskPayload) (string, error) {
	if payload.GroupSettingsId == 0 {
		return utils.UintToString(payload.GroupId), nil
	}

	groupSettingsRepository := repositories.NewGroupSettingsRepository(nil)

	groupSettings, err := groupSettingsRepository.GetGroupSettingsById(utils.UintToString(payload.GroupSettingsId))
	if err != nil {
		return "", err
	}

	return utils.UintToString(groupSettings.GroupId), nil
}
