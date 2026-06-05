package handlers

import (
	"context"
	"encoding/json"
	"net/http/httptest"
	"receipt-wrangler/api/internal/models"
	"receipt-wrangler/api/internal/permissions"
	"receipt-wrangler/api/internal/structs"
	"receipt-wrangler/api/internal/utils"
	"strings"
	"testing"

	jwtmiddleware "github.com/auth0/go-jwt-middleware/v2"
	"github.com/auth0/go-jwt-middleware/v2/validator"
)

func TestShouldGetAllPermissions(t *testing.T) {
	defer tearDownGenericHandlerTest()
	reader := strings.NewReader("")
	w := httptest.NewRecorder()
	r := httptest.NewRequest("GET", "/api/permission", reader)

	newContext := context.WithValue(
		r.Context(),
		jwtmiddleware.ContextKey{},
		&validator.ValidatedClaims{CustomClaims: &structs.Claims{UserId: 1, UserRole: models.USER}},
	)
	r = r.WithContext(newContext)

	grantAppPerms(t, 1, permissions.AppRolesRead)

	GetPermissions(w, r)

	if w.Result().StatusCode != 200 {
		utils.PrintTestError(t, w.Result().StatusCode, 200)
		return
	}

	var got []permissions.Descriptor
	if err := json.Unmarshal(w.Body.Bytes(), &got); err != nil {
		utils.PrintTestError(t, err, nil)
		return
	}

	want := permissions.All()
	if len(got) != len(want) {
		utils.PrintTestError(t, len(got), len(want))
		return
	}

	for i, d := range want {
		if got[i] != d {
			utils.PrintTestError(t, got[i], d)
			return
		}
	}
}
