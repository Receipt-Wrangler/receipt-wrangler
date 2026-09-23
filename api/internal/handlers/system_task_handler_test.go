package handlers

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"receipt-wrangler/api/internal/commands"
	"receipt-wrangler/api/internal/models"
	"receipt-wrangler/api/internal/repositories"
	"receipt-wrangler/api/internal/structs"
	"receipt-wrangler/api/internal/utils"
	"strings"
	"testing"
	"time"

	jwtmiddleware "github.com/auth0/go-jwt-middleware/v2"
	"github.com/auth0/go-jwt-middleware/v2/validator"
)

func tearDownSystemTaskTest() {
	repositories.TruncateTestDb()
}

func TestShouldNotAllowUserToGetSystemTasks(t *testing.T) {
	defer tearDownSystemTaskTest()
	reader := strings.NewReader("")
	w := httptest.NewRecorder()
	r := httptest.NewRequest("POST", "/api", reader)

	newContext := context.WithValue(r.Context(), jwtmiddleware.ContextKey{}, &validator.ValidatedClaims{CustomClaims: &structs.Claims{UserId: 1}})
	r = r.WithContext(newContext)

	GetSystemTasks(w, r)

	if w.Result().StatusCode != http.StatusForbidden {
		utils.PrintTestError(t, w.Result().StatusCode, http.StatusForbidden)
	}
}

func TestShouldNotAllowAdminToGetSystemTasksWithInvalidCommand(t *testing.T) {
	defer tearDownSystemTaskTest()

	tests := map[string]struct {
		input  commands.GetSystemTaskCommand
		expect int
	}{
		"empty body": {
			expect: http.StatusBadRequest,
		},
		"empty command": {
			input:  commands.GetSystemTaskCommand{},
			expect: http.StatusBadRequest,
		},
		"missing count": {
			input: commands.GetSystemTaskCommand{
				PagedRequestCommand: commands.PagedRequestCommand{
					Page:          1,
					PageSize:      100,
					OrderBy:       "type",
					SortDirection: "asc",
				},
				AssociatedEntityId:   1,
				AssociatedEntityType: models.SYSTEM_EMAIL,
			},
			expect: http.StatusOK,
		},
		"missing associated entityId": {
			input: commands.GetSystemTaskCommand{
				PagedRequestCommand: commands.PagedRequestCommand{
					Page:          1,
					PageSize:      100,
					OrderBy:       "associated_entity_type",
					SortDirection: "asc",
				},
				AssociatedEntityType: models.SYSTEM_EMAIL,
			},
			expect: http.StatusOK,
		},
		"bad order by": {
			input: commands.GetSystemTaskCommand{
				PagedRequestCommand: commands.PagedRequestCommand{
					Page:          1,
					PageSize:      100,
					OrderBy:       "terrible order by",
					SortDirection: "asc",
				},
				AssociatedEntityType: models.SYSTEM_EMAIL,
			},
			expect: http.StatusInternalServerError,
		},
		"missing associated entityType": {
			input: commands.GetSystemTaskCommand{
				AssociatedEntityId: 1,
			},
			expect: http.StatusBadRequest,
		},
	}

	grantAllAppPerms(t, 1)

	for _, test := range tests {
		bytes, _ := json.Marshal(test.input)
		reader := strings.NewReader(string(bytes))
		w := httptest.NewRecorder()
		r := httptest.NewRequest("POST", "/api", reader)

		newContext := context.WithValue(r.Context(), jwtmiddleware.ContextKey{}, &validator.ValidatedClaims{CustomClaims: &structs.Claims{UserId: 1}})
		r = r.WithContext(newContext)

		GetSystemTasks(w, r)

		if w.Result().StatusCode != test.expect {
			utils.PrintTestError(t, w.Result().StatusCode, test.expect)
		}
	}
}

// An "Updated Receipt" row written before the version marker is served with its
// incomplete "before" replaced by the receipt's earlier complete copy (here the
// one stored when it was created). services.UpcastReceiptUpdateDescriptions
// covers the matching rules; this pins that the listing applies it.
func TestGetSystemTasksRebuildsVersionOneReceiptUpdates(t *testing.T) {
	defer tearDownSystemTaskTest()
	grantAllAppPerms(t, 1)

	db := repositories.GetDB()
	receiptId := uint(7)
	created := `{"id":7,"name":"Created"}`
	uploaded := models.SystemTask{
		Type:                 models.RECEIPT_UPLOADED,
		Status:               models.SYSTEM_TASK_SUCCEEDED,
		AssociatedEntityType: models.RECEIPT,
		AssociatedEntityId:   receiptId,
		ReceiptId:            &receiptId,
		StartedAt:            time.Now(),
		ResultDescription:    created,
	}
	updated := models.SystemTask{
		Type:                 models.RECEIPT_UPDATED,
		Status:               models.SYSTEM_TASK_SUCCEEDED,
		AssociatedEntityType: models.RECEIPT,
		AssociatedEntityId:   receiptId,
		ReceiptId:            &receiptId,
		StartedAt:            time.Now(),
		ResultDescription:    `{"before":"{\"id\":7}","after":"{\"id\":7,\"name\":\"One\"}"}`,
	}
	for _, task := range []*models.SystemTask{&uploaded, &updated} {
		if err := db.Create(task).Error; err != nil {
			t.Fatalf("creating system task: %v", err)
		}
	}

	body, _ := json.Marshal(commands.GetSystemTaskCommand{
		PagedRequestCommand: commands.PagedRequestCommand{
			Page:          1,
			PageSize:      10,
			OrderBy:       "started_at",
			SortDirection: "desc",
		},
	})
	w := httptest.NewRecorder()
	r := httptest.NewRequest("POST", "/api", strings.NewReader(string(body)))
	r = r.WithContext(context.WithValue(r.Context(), jwtmiddleware.ContextKey{}, &validator.ValidatedClaims{CustomClaims: &structs.Claims{UserId: 1}}))

	GetSystemTasks(w, r)

	if w.Result().StatusCode != http.StatusOK {
		t.Fatalf("status = %d, want 200; body=%s", w.Result().StatusCode, w.Body.String())
	}
	var paged struct {
		Data []models.SystemTask `json:"data"`
	}
	if err := json.Unmarshal(w.Body.Bytes(), &paged); err != nil {
		t.Fatalf("reading response: %v", err)
	}
	if len(paged.Data) != 1 || paged.Data[0].ID != updated.ID {
		t.Fatalf("expected only the update row, got %+v", paged.Data)
	}

	var description struct {
		Before       string `json:"before"`
		Version      int    `json:"version"`
		BeforeSource struct {
			SystemTaskId uint `json:"systemTaskId"`
		} `json:"beforeSource"`
	}
	if err := json.Unmarshal([]byte(paged.Data[0].ResultDescription), &description); err != nil {
		t.Fatalf("reading description: %v", err)
	}
	if description.Before != created || description.Version != 1 || description.BeforeSource.SystemTaskId != uploaded.ID {
		t.Errorf("description = %+v, want before %s from task %d at version 1", description, created, uploaded.ID)
	}
}
