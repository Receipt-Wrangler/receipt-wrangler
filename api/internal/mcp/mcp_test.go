package mcp

import (
	"context"
	"net/http"
	"net/http/httptest"
	"receipt-wrangler/api/internal/models"
	"receipt-wrangler/api/internal/repositories"
	"receipt-wrangler/api/internal/services"
	"receipt-wrangler/api/internal/structs"
	"receipt-wrangler/api/internal/utils"
	"strings"
	"testing"
	"time"

	"github.com/modelcontextprotocol/go-sdk/auth"
	"github.com/shopspring/decimal"
	"gorm.io/gorm/clause"
)

func TestNewHandlerRejectsMissingToken(t *testing.T) {
	handler := NewHandler()

	recorder := httptest.NewRecorder()
	handler.ServeHTTP(recorder, httptest.NewRequest(http.MethodGet, "/mcp", nil))

	if recorder.Code != http.StatusUnauthorized {
		t.Fatalf("expected 401 without a bearer token, got %d", recorder.Code)
	}

	wwwAuth := recorder.Header().Get("WWW-Authenticate")
	if !strings.Contains(wwwAuth, "resource_metadata") {
		t.Errorf("expected WWW-Authenticate to advertise resource_metadata, got %q", wwwAuth)
	}
}

func TestVerifyTokenAcceptsValidJwt(t *testing.T) {
	defer repositories.TruncateTestDb()

	user := createUser(t, "tokenuser")
	accessToken, _, _, err := services.GenerateJWT(user.ID)
	if err != nil {
		t.Fatalf("failed to generate jwt: %v", err)
	}

	tokenValidator, err := services.InitTokenValidator()
	if err != nil {
		t.Fatalf("failed to init validator: %v", err)
	}

	info, err := verifyToken(context.Background(), tokenValidator, accessToken)
	if err != nil {
		t.Fatalf("verifyToken rejected a valid token: %v", err)
	}
	if info.UserID != utils.UintToString(user.ID) {
		t.Errorf("token UserID = %q, want %q", info.UserID, utils.UintToString(user.ID))
	}
	if _, ok := info.Extra[claimsKey]; !ok {
		t.Errorf("expected verified claims to be stored on the token info")
	}
}

func TestHandlerAcceptsIssuedToken(t *testing.T) {
	defer repositories.TruncateTestDb()

	user := createUser(t, "bearer")
	accessToken, _, _, err := services.GenerateJWT(user.ID)
	if err != nil {
		t.Fatalf("failed to generate jwt: %v", err)
	}

	handler := NewHandler()
	request := httptest.NewRequest(http.MethodGet, "/mcp", nil)
	request.Header.Set("Authorization", "Bearer "+accessToken)

	recorder := httptest.NewRecorder()
	handler.ServeHTTP(recorder, request)

	// A valid token must clear the auth layer. The transport may then reject a
	// bare GET for protocol reasons, but it must not be a 401.
	if recorder.Code == http.StatusUnauthorized {
		t.Fatalf("a valid issued token was rejected by the mcp auth layer")
	}
}

func TestVerifyTokenRejectsInvalidJwt(t *testing.T) {
	tokenValidator, err := services.InitTokenValidator()
	if err != nil {
		t.Fatalf("failed to init validator: %v", err)
	}

	_, err = verifyToken(context.Background(), tokenValidator, "not-a-real-token")
	if err == nil {
		t.Fatalf("expected an error for an invalid token")
	}
	if !strings.Contains(err.Error(), auth.ErrInvalidToken.Error()) {
		t.Errorf("expected error to wrap ErrInvalidToken, got %v", err)
	}
}

func createReceiptInGroup(t *testing.T, name string, groupId uint, paidByUserId uint) models.Receipt {
	t.Helper()
	receipt := models.Receipt{
		Name:         name,
		Amount:       decimal.NewFromInt(10),
		Date:         time.Now(),
		Status:       models.OPEN,
		GroupId:      groupId,
		PaidByUserID: paidByUserId,
	}
	if err := repositories.GetDB().Omit(clause.Associations).Create(&receipt).Error; err != nil {
		t.Fatalf("failed to create receipt: %v", err)
	}
	return receipt
}

func addGroupMember(t *testing.T, userId uint, groupId uint, role models.GroupRole) {
	t.Helper()
	member := models.GroupMember{UserID: userId, GroupID: groupId, GroupRole: role}
	if err := repositories.GetDB().Create(&member).Error; err != nil {
		t.Fatalf("failed to create group member: %v", err)
	}
}

func TestListCategoriesReturnsCategories(t *testing.T) {
	defer repositories.TruncateTestDb()

	user := createUser(t, "catuser")
	repositories.CreateTestCategories()

	_, out, err := handleListCategories(context.Background(), requestForUser(user.ID), emptyInput{})
	if err != nil {
		t.Fatalf("handleListCategories returned error: %v", err)
	}

	categories, ok := out.([]models.Category)
	if !ok {
		t.Fatalf("expected []models.Category, got %T", out)
	}
	if len(categories) != 3 {
		t.Errorf("expected 3 categories, got %d", len(categories))
	}
}

func TestGetReceiptEnforcesGroupAccess(t *testing.T) {
	defer repositories.TruncateTestDb()

	member := createUser(t, "member")
	outsider := createUser(t, "outsider")

	group := models.Group{Name: "g1"}
	if err := repositories.GetDB().Create(&group).Error; err != nil {
		t.Fatalf("failed to create group: %v", err)
	}
	addGroupMember(t, member.ID, group.ID, models.VIEWER)

	receipt := createReceiptInGroup(t, "Lunch", group.ID, member.ID)
	receiptId := utils.UintToString(receipt.ID)

	// A member of the group can read the receipt.
	_, out, err := handleGetReceipt(context.Background(), requestForUser(member.ID), getReceiptInput{Id: receiptId})
	if err != nil {
		t.Fatalf("group member could not read receipt: %v", err)
	}
	if got, ok := out.(models.Receipt); !ok || got.ID != receipt.ID {
		t.Errorf("expected the requested receipt to be returned")
	}

	// A non-member is denied, with an error that does not confirm existence.
	_, _, err = handleGetReceipt(context.Background(), requestForUser(outsider.ID), getReceiptInput{Id: receiptId})
	if err == nil {
		t.Fatalf("expected non-member to be denied access to the receipt")
	}
	if !strings.Contains(err.Error(), "not found") {
		t.Errorf("expected a non-leaking 'not found' error, got %v", err)
	}
}

func TestSearchReceiptsScopesToUserGroups(t *testing.T) {
	defer repositories.TruncateTestDb()

	user := createUser(t, "searcher")

	memberGroup := models.Group{Name: "mine"}
	otherGroup := models.Group{Name: "theirs"}
	if err := repositories.GetDB().Create(&memberGroup).Error; err != nil {
		t.Fatalf("failed to create member group: %v", err)
	}
	if err := repositories.GetDB().Create(&otherGroup).Error; err != nil {
		t.Fatalf("failed to create other group: %v", err)
	}
	addGroupMember(t, user.ID, memberGroup.ID, models.VIEWER)

	createReceiptInGroup(t, "Coffee shop", memberGroup.ID, user.ID)
	createReceiptInGroup(t, "Coffee beans", memberGroup.ID, user.ID)
	createReceiptInGroup(t, "Coffee elsewhere", otherGroup.ID, user.ID)

	_, out, err := handleSearchReceipts(context.Background(), requestForUser(user.ID), searchReceiptsInput{Query: "Coffee"})
	if err != nil {
		t.Fatalf("handleSearchReceipts returned error: %v", err)
	}

	results, ok := out.([]structs.SearchResult)
	if !ok {
		t.Fatalf("expected []structs.SearchResult, got %T", out)
	}
	if len(results) != 2 {
		t.Errorf("expected 2 receipts from the user's group, got %d", len(results))
	}
	for _, result := range results {
		if result.GroupID != memberGroup.ID {
			t.Errorf("search returned a receipt from group %d outside the user's groups", result.GroupID)
		}
	}
}
