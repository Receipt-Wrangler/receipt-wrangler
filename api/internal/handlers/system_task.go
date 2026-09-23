package handlers

import (
	"encoding/json"
	"errors"
	"github.com/go-chi/chi/v5"
	"net/http"
	"os"
	"path/filepath"
	"receipt-wrangler/api/internal/commands"
	"receipt-wrangler/api/internal/constants"
	"receipt-wrangler/api/internal/logging"
	"receipt-wrangler/api/internal/models"
	"receipt-wrangler/api/internal/permissions"
	"receipt-wrangler/api/internal/repositories"
	"receipt-wrangler/api/internal/services"
	"receipt-wrangler/api/internal/structs"
	"receipt-wrangler/api/internal/utils"
	"receipt-wrangler/api/internal/wranglerasynq"
)

func GetSystemTasks(w http.ResponseWriter, r *http.Request) {
	handler := structs.Handler{
		ErrorMessage:   "Error getting system tasks",
		Writer:         w,
		Request:        r,
		AppPermissions: []string{permissions.AppSystemTasksRead},
		ResponseType:   constants.ApplicationJson,
		HandlerFunction: func(w http.ResponseWriter, r *http.Request) (int, error) {
			command := commands.GetSystemTaskCommand{}
			err := command.LoadDataFromRequest(w, r)
			if err != nil {
				return http.StatusInternalServerError, err
			}

			vErr := command.Validate()
			if len(vErr.Errors) > 0 {
				structs.WriteValidatorErrorResponse(w, vErr, http.StatusBadRequest)
				return 0, nil
			}

			systemTaskRepository := repositories.NewSystemTaskRepository(nil)
			systemTasks, count, err := systemTaskRepository.GetPagedSystemTasks(command)
			if err != nil {
				return http.StatusInternalServerError, err
			}

			// Hydrate before the copy loop below, which takes each task by value.
			// This table is app-scoped, so it lists groups the caller may not be a
			// member of; the flag is resolved per caller so it never advertises a
			// file the source-file endpoints would refuse to serve them.
			token := structs.GetClaims(r)
			err = wranglerasynq.SetSystemTaskHasSourceFile(systemTasks, groupSourceFileReader(token.UserId))
			if err != nil {
				return http.StatusInternalServerError, err
			}

			pagedData := structs.PagedData{}
			data := make([]any, 0)

			for i := 0; i < len(systemTasks); i++ {
				data = append(data, systemTasks[i])
			}

			pagedData.Data = data
			pagedData.TotalCount = count

			responseBytes, err := utils.MarshalResponseData(pagedData)
			if err != nil {
				return http.StatusInternalServerError, err
			}

			w.WriteHeader(http.StatusOK)
			w.Write(responseBytes)

			return 0, nil
		},
	}

	HandleRequest(handler)
}

func GetActivitiesForGroups(w http.ResponseWriter, r *http.Request) {
	errorMsg := "Error getting group activities"
	command := commands.PagedActivityRequestCommand{}
	err := command.LoadDataFromRequest(w, r)
	if err != nil {
		utils.WriteCustomErrorResponse(w, errorMsg, http.StatusInternalServerError)
		return
	}

	stringGroupIds := make([]string, 0)
	for _, groupId := range command.GroupIds {
		stringGroupIds = append(stringGroupIds, utils.UintToString(groupId))
	}

	handler := structs.Handler{
		ErrorMessage:     errorMsg,
		Writer:           w,
		Request:          r,
		GroupIds:         stringGroupIds,
		GroupPermissions: []string{permissions.GroupActivitiesRead},
		ResponseType:     constants.ApplicationJson,
		HandlerFunction: func(w http.ResponseWriter, r *http.Request) (int, error) {

			vErr := command.Validate()
			if len(vErr.Errors) > 0 {
				structs.WriteValidatorErrorResponse(w, vErr, http.StatusBadRequest)
				return 0, nil
			}

			systemTaskRepository := repositories.NewSystemTaskRepository(nil)
			token := structs.GetClaims(r)

			// Member isolation: activities run by a user the caller may not see in that
			// activity's group are filtered IN THE QUERY (before Count + pagination), so
			// TotalCount and the returned page both reflect only visible rows and DB-side
			// LIMIT/OFFSET is preserved. See applyActivityVisibilityDisjunction (mirrors
			// the paid-by disjunction).
			permissionService := services.NewPermissionService(nil)
			resolver, err := permissionService.ActivityVisibilityResolver(token.UserId, command.GroupIds)
			if err != nil {
				return http.StatusInternalServerError, err
			}
			activities, count, err := systemTaskRepository.GetPagedActivities(command, resolver)
			if err != nil {
				return http.StatusInternalServerError, err
			}

			err = wranglerasynq.SetActivityFlags(&activities)
			if err != nil {
				return http.StatusInternalServerError, err
			}

			pagedData := structs.PagedData{}
			data := make([]any, 0)

			for i := 0; i < len(activities); i++ {
				data = append(data, activities[i])
			}

			pagedData.Data = data
			pagedData.TotalCount = count

			responseBytes, err := utils.MarshalResponseData(pagedData)
			if err != nil {
				return http.StatusInternalServerError, err
			}

			w.WriteHeader(http.StatusOK)
			w.Write(responseBytes)

			return 0, nil
		},
	}

	HandleRequest(handler)
}

func RerunActivity(w http.ResponseWriter, r *http.Request) {
	errorMsg := "Error rerunning activity"
	systemTaskRepository := repositories.NewSystemTaskRepository(nil)
	inspector, err := wranglerasynq.GetAsynqInspector()
	if err != nil {
		logging.LogStd(logging.LOG_LEVEL_ERROR, err.Error())
		utils.WriteCustomErrorResponse(w, errorMsg, http.StatusInternalServerError)
		return
	}
	defer inspector.Close()

	systemTaskId := chi.URLParam(r, "id")
	systemTaskUintId, err := utils.StringToUint(systemTaskId)
	if err != nil {
		logging.LogStd(logging.LOG_LEVEL_ERROR, err.Error())
		utils.WriteCustomErrorResponse(w, errorMsg, http.StatusInternalServerError)
		return
	}

	systemTask, err := systemTaskRepository.GetSystemTaskById(systemTaskUintId)
	if err != nil {
		logging.LogStd(logging.LOG_LEVEL_ERROR, err.Error())
		utils.WriteCustomErrorResponse(w, errorMsg, http.StatusInternalServerError)
		return
	}

	if systemTask.Type != models.QUICK_SCAN && systemTask.Type != models.EMAIL_UPLOAD {
		logging.LogStd(logging.LOG_LEVEL_ERROR, "Only quick scan and email upload activities can be rerun")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	queueName, err := wranglerasynq.SystemTaskToQueueName(systemTask.Type)
	if err != nil {
		logging.LogStd(logging.LOG_LEVEL_ERROR, err.Error())
		utils.WriteCustomErrorResponse(w, errorMsg, http.StatusInternalServerError)
		return
	}

	taskInfo, err := inspector.GetTaskInfo(queueName, systemTask.AsynqTaskId)
	if err != nil {
		logging.LogStd(logging.LOG_LEVEL_ERROR, err.Error())
		utils.WriteCustomErrorResponse(w, errorMsg, http.StatusInternalServerError)
		return
	}

	var payload wranglerasynq.RerunTaskPayload
	err = json.Unmarshal(taskInfo.Payload, &payload)
	if err != nil {
		logging.LogStd(logging.LOG_LEVEL_ERROR, err.Error())
		utils.WriteCustomErrorResponse(w, errorMsg, http.StatusInternalServerError)
		return
	}

	stringGroupId, err := wranglerasynq.ResolveActivityGroupId(payload)
	if err != nil {
		logging.LogStd(logging.LOG_LEVEL_ERROR, err.Error())
		utils.WriteCustomErrorResponse(w, errorMsg, http.StatusInternalServerError)
		return
	}

	handler := structs.Handler{
		ErrorMessage:     errorMsg,
		Writer:           w,
		Request:          r,
		GroupId:          stringGroupId,
		GroupPermissions: []string{permissions.GroupActivitiesRerun},
		HandlerFunction: func(w http.ResponseWriter, r *http.Request) (int, error) {
			// Inspector.RunTask does not refuse a task by state — it pushes back
			// anything that is not active or pending, a succeeded task included —
			// so the rule lives here. The email queue retains its completed tasks
			// for the temp sweep, so without this a rerun of one would re-process
			// an email that already produced a receipt. See CanRerunTask.
			if !wranglerasynq.CanRerunTask(taskInfo) {
				utils.WriteCustomErrorResponse(w, "Only a failed activity can be rerun.", http.StatusBadRequest)
				return 0, nil
			}

			// Defence in depth: the client hides the control once the source file
			// is gone, but the endpoint stays callable. A rerun reads its upload
			// first, so without this it would fail deep in the pipeline instead of
			// saying what is wrong.
			if !wranglerasynq.RerunSourceFilesPresent(systemTask.Type, taskInfo.Payload) {
				utils.WriteCustomErrorResponse(w, "The file this activity needs is no longer available.", http.StatusBadRequest)
				return 0, nil
			}

			err = inspector.RunTask(queueName, systemTask.AsynqTaskId)
			if err != nil {
				return http.StatusInternalServerError, err
			}

			return 0, nil
		},
	}

	HandleRequest(handler)
}

// groupSourceFileReader answers "may this user reach source files in that group",
// memoized per group for the life of one request. It is the same gate the
// source-file endpoints apply, so a flag can never promise a button that 403s.
func groupSourceFileReader(userId uint) func(groupId uint) bool {
	permissionService := services.NewPermissionService(nil)
	cache := make(map[uint]bool)

	return func(groupId uint) bool {
		if allowed, ok := cache[groupId]; ok {
			return allowed
		}

		allowed, err := permissionService.HasGroupPermissions(userId, groupId, permissions.GroupActivitiesRead)
		if err != nil {
			logging.LogStd(logging.LOG_LEVEL_ERROR, err.Error())
			allowed = false
		}

		cache[groupId] = allowed

		return allowed
	}
}

// GetSystemTaskSourceFile returns the upload behind a failed activity, converted
// for display.
func GetSystemTaskSourceFile(w http.ResponseWriter, r *http.Request) {
	errorMsg := "Error getting activity source file"

	sourceFile, systemTaskErr := resolveSourceFileFromRequest(w, r, errorMsg)
	if systemTaskErr {
		return
	}

	handler := structs.Handler{
		ErrorMessage:     errorMsg,
		Writer:           w,
		Request:          r,
		GroupId:          sourceFile.GroupId,
		GroupPermissions: []string{permissions.GroupActivitiesRead},
		ResponseType:     constants.ApplicationJson,
		HandlerFunction: func(w http.ResponseWriter, r *http.Request) (int, error) {
			fileRepository := repositories.NewFileRepository(nil)

			// Prefer the already-converted OCR copy. Converting the original
			// instead can mean rasterizing a multi-page PDF at the configured DPI
			// inside this request.
			pathToRead := sourceFile.PreviewPath
			if len(pathToRead) == 0 || !utils.FileExists(pathToRead) {
				pathToRead = sourceFile.Path
			}

			// os.ReadFile, never utils.ReadFile: that one returns (nil, nil) on a
			// read error, which would serve an empty image as a success.
			fileBytes, err := os.ReadFile(pathToRead)
			if err != nil {
				return http.StatusInternalServerError, err
			}

			encodedImage, err := fileRepository.BuildEncodedImageString(fileBytes)
			if err != nil {
				return http.StatusInternalServerError, err
			}

			responseBytes, err := utils.MarshalResponseData(structs.SystemTaskSourceFileView{
				Name:         sourceFile.FileName,
				EncodedImage: encodedImage,
			})
			if err != nil {
				return http.StatusInternalServerError, err
			}

			w.WriteHeader(http.StatusOK)
			w.Write(responseBytes)

			return 0, nil
		},
	}

	HandleRequest(handler)
}

// DownloadSystemTaskSourceFile serves the upload behind a failed activity
// verbatim, so the user can file the receipt by hand.
func DownloadSystemTaskSourceFile(w http.ResponseWriter, r *http.Request) {
	errorMsg := "Error downloading activity source file"

	sourceFile, systemTaskErr := resolveSourceFileFromRequest(w, r, errorMsg)
	if systemTaskErr {
		return
	}

	handler := structs.Handler{
		ErrorMessage:     errorMsg,
		Writer:           w,
		Request:          r,
		GroupId:          sourceFile.GroupId,
		GroupPermissions: []string{permissions.GroupActivitiesRead},
		// Deliberately unset so http.ServeFile derives the Content-Type itself.
		ResponseType: "",
		HandlerFunction: func(w http.ResponseWriter, r *http.Request) (int, error) {
			// Quoted and sanitized: an email attachment's name comes from a MIME
			// header, so it can carry spaces, commas or separators.
			fileName := utils.SanitizeFileName(sourceFile.FileName)
			if len(fileName) == 0 {
				fileName = filepath.Base(sourceFile.Path)
			}

			w.Header().Set("Content-Disposition", "attachment; filename=\""+fileName+"\"")
			http.ServeFile(w, r, sourceFile.Path)

			// Streaming has begun, so returning an error here would write an error
			// body over the file.
			return 0, nil
		},
	}

	HandleRequest(handler)
}

// resolveSourceFileFromRequest loads the system task named in the URL and locates
// its upload. It reports true when it has already written a response, in which
// case the caller must return without building a handler.
//
// This runs before the structs.Handler is built because the group it gates on is
// only knowable from the task's asynq payload — the same shape RerunActivity uses.
func resolveSourceFileFromRequest(
	w http.ResponseWriter,
	r *http.Request,
	errorMsg string,
) (wranglerasynq.SystemTaskSourceFile, bool) {
	systemTaskRepository := repositories.NewSystemTaskRepository(nil)

	systemTaskId, err := utils.StringToUint(chi.URLParam(r, "id"))
	if err != nil {
		logging.LogStd(logging.LOG_LEVEL_ERROR, err.Error())
		utils.WriteCustomErrorResponse(w, errorMsg, http.StatusBadRequest)
		return wranglerasynq.SystemTaskSourceFile{}, true
	}

	systemTask, err := systemTaskRepository.GetSystemTaskById(systemTaskId)
	if err != nil {
		logging.LogStd(logging.LOG_LEVEL_ERROR, err.Error())
		utils.WriteCustomErrorResponse(w, errorMsg, http.StatusInternalServerError)
		return wranglerasynq.SystemTaskSourceFile{}, true
	}

	sourceFile, err := wranglerasynq.ResolveSystemTaskSourceFile(systemTask)
	if err != nil {
		if errors.Is(err, wranglerasynq.ErrSourceFileUnsupportedTask) {
			utils.WriteCustomErrorResponse(w, "This activity type has no source file.", http.StatusBadRequest)
			return wranglerasynq.SystemTaskSourceFile{}, true
		}

		if errors.Is(err, wranglerasynq.ErrSourceFileUnavailable) {
			utils.WriteCustomErrorResponse(w, "The source file is no longer available.", http.StatusNotFound)
			return wranglerasynq.SystemTaskSourceFile{}, true
		}

		logging.LogStd(logging.LOG_LEVEL_ERROR, err.Error())
		utils.WriteCustomErrorResponse(w, errorMsg, http.StatusInternalServerError)
		return wranglerasynq.SystemTaskSourceFile{}, true
	}

	return sourceFile, false
}
