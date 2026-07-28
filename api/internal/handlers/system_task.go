package handlers

import (
	"encoding/json"
	"github.com/go-chi/chi/v5"
	"net/http"
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

			// Fetch ALL matching activities unpaged so member-isolation filtering runs
			// BEFORE count + pagination. Filtering an already-paged result would leak
			// hidden-peer presence through TotalCount and could return a short page,
			// since the DB limited/offset before the invisible rows were known.
			unpagedCommand := command
			unpagedCommand.Page = -1
			unpagedCommand.PageSize = -1
			activities, _, err := systemTaskRepository.GetPagedActivities(unpagedCommand)
			if err != nil {
				return http.StatusInternalServerError, err
			}

			// Member isolation: drop activities run by a user the caller may not see
			// IN THAT ACTIVITY'S GROUP (a nil RanByUserId is a system action and is
			// always kept). Resolved per group, so an isolated group hides its
			// non-visible actors regardless of any open group the caller also shares.
			token := structs.GetClaims(r)
			activities, err = filterActivitiesByVisibility(services.NewPermissionService(nil), token.UserId, activities)
			if err != nil {
				return http.StatusInternalServerError, err
			}

			// Count and paginate the FILTERED set, so both the total and the returned
			// page reflect only what the caller may see.
			totalCount := int64(len(activities))
			activities = paginateActivities(activities, command.Page, command.PageSize)

			err = wranglerasynq.SetActivityCanBeRestarted(&activities)
			if err != nil {
				return http.StatusInternalServerError, err
			}

			pagedData := structs.PagedData{}
			data := make([]any, 0)

			for i := 0; i < len(activities); i++ {
				data = append(data, activities[i])
			}

			pagedData.Data = data
			pagedData.TotalCount = totalCount

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

// filterActivitiesByVisibility drops activities run by a user the viewer may not see IN
// THAT ACTIVITY'S GROUP, under member-presence isolation. An activity with a nil
// RanByUserId (system action) or a nil GroupId (no group to scope against) is always
// kept. Visibility is resolved per group and memoized, so an isolated group hides its
// non-visible actors regardless of any open group the viewer also shares; a group where
// the viewer is unrestricted (admin, supervisor, non-isolated) keeps every actor.
func filterActivitiesByVisibility(permissionService services.PermissionService, viewerId uint, activities []structs.Activity) ([]structs.Activity, error) {
	if len(activities) == 0 {
		return activities, nil
	}

	type groupVis struct {
		visible      map[uint]struct{}
		unrestricted bool
	}
	cache := map[uint]groupVis{}

	visible := make([]structs.Activity, 0, len(activities))
	for _, activity := range activities {
		if activity.RanByUserId == nil || activity.GroupId == nil {
			visible = append(visible, activity)
			continue
		}
		groupId := *activity.GroupId
		v, ok := cache[groupId]
		if !ok {
			set, unrestricted, err := permissionService.GetVisibleUserIdsForUserInGroup(viewerId, groupId)
			if err != nil {
				return nil, err
			}
			v = groupVis{visible: set, unrestricted: unrestricted}
			cache[groupId] = v
		}
		if v.unrestricted {
			visible = append(visible, activity)
			continue
		}
		if _, actorOk := v.visible[*activity.RanByUserId]; actorOk {
			visible = append(visible, activity)
		}
	}
	return visible, nil
}

// paginateActivities slices an already-filtered activity list to the requested page,
// mirroring repositories.BaseRepository.Paginate so in-memory pagination matches the DB
// semantics (1-indexed page; pageSize -1 = all; pageSize clamped to [1,100] with a
// default of 10; page <= 0 treated as 1).
func paginateActivities(activities []structs.Activity, page int, pageSize int) []structs.Activity {
	if pageSize == -1 {
		return activities
	}
	if page <= 0 {
		page = 1
	}
	switch {
	case pageSize > 100:
		pageSize = 100
	case pageSize <= 0:
		pageSize = 10
	}

	offset := (page - 1) * pageSize
	if offset >= len(activities) {
		return []structs.Activity{}
	}
	end := offset + pageSize
	if end > len(activities) {
		end = len(activities)
	}
	return activities[offset:end]
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

	stringGroupId := utils.UintToString(payload.GroupId)
	if payload.GroupSettingsId > 0 {
		groupSettingsRepository := repositories.NewGroupSettingsRepository(nil)

		stringId := utils.UintToString(payload.GroupSettingsId)
		groupSettings, err := groupSettingsRepository.GetGroupSettingsById(stringId)
		if err != nil {
			logging.LogStd(logging.LOG_LEVEL_ERROR, err.Error())
			utils.WriteCustomErrorResponse(w, errorMsg, http.StatusInternalServerError)
			return
		}

		stringGroupId = utils.UintToString(groupSettings.GroupId)
	}

	handler := structs.Handler{
		ErrorMessage:     errorMsg,
		Writer:           w,
		Request:          r,
		GroupId:          stringGroupId,
		GroupPermissions: []string{permissions.GroupActivitiesRerun},
		HandlerFunction: func(w http.ResponseWriter, r *http.Request) (int, error) {
			err = inspector.RunTask(queueName, systemTask.AsynqTaskId)
			if err != nil {
				return http.StatusInternalServerError, err
			}

			return 0, nil
		},
	}

	HandleRequest(handler)
}
