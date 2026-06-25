package mcp

import (
	"context"
	"errors"
	"receipt-wrangler/api/internal/permissions"
	"receipt-wrangler/api/internal/repositories"
	"receipt-wrangler/api/internal/services"
	"receipt-wrangler/api/internal/structs"
	"receipt-wrangler/api/internal/utils"

	mcpsdk "github.com/modelcontextprotocol/go-sdk/mcp"
)

// maxSearchResults bounds how many receipts a single search returns, matching
// the limit used by the REST search handler.
const maxSearchResults = 100

// emptyInput is the parameter type for tools that take no arguments. It yields
// an empty JSON object schema.
type emptyInput struct{}

type searchReceiptsInput struct {
	Query      string `json:"query"`
	MaxResults int    `json:"maxResults,omitempty"`
}

type getReceiptInput struct {
	Id string `json:"id"`
}

type listDashboardsInput struct {
	GroupId string `json:"groupId"`
}

// registerTools wires the read-only MCP tools onto the server. Every handler
// derives the acting user from the verified token and enforces the same
// group-scoped authorization the REST handlers do.
func registerTools(server *mcpsdk.Server) {
	mcpsdk.AddTool(server, &mcpsdk.Tool{
		Name:        "search_receipts",
		Description: "Search the authenticated user's receipts by name. Returns receipts from groups the user belongs to, most recent first.",
	}, handleSearchReceipts)

	mcpsdk.AddTool(server, &mcpsdk.Tool{
		Name:        "get_receipt",
		Description: "Get a single receipt by id, including its items, categories, and tags. Only receipts in the user's groups are accessible.",
	}, handleGetReceipt)

	mcpsdk.AddTool(server, &mcpsdk.Tool{
		Name:        "list_groups",
		Description: "List the groups the authenticated user belongs to.",
	}, handleListGroups)

	mcpsdk.AddTool(server, &mcpsdk.Tool{
		Name:        "list_categories",
		Description: "List all receipt categories available in the system.",
	}, handleListCategories)

	mcpsdk.AddTool(server, &mcpsdk.Tool{
		Name:        "list_tags",
		Description: "List all receipt tags available in the system.",
	}, handleListTags)

	mcpsdk.AddTool(server, &mcpsdk.Tool{
		Name:        "list_dashboards",
		Description: "List the authenticated user's dashboards for a given group. Requires a groupId the user can access.",
	}, handleListDashboards)
}

func handleSearchReceipts(ctx context.Context, req *mcpsdk.CallToolRequest, in searchReceiptsInput) (*mcpsdk.CallToolResult, any, error) {
	claims, err := claimsFromRequest(req)
	if err != nil {
		return nil, nil, err
	}

	limit := in.MaxResults
	if limit <= 0 || limit > maxSearchResults {
		limit = maxSearchResults
	}

	// Scope to the user's groups exactly like handlers.Search, then delegate
	// the query to the repository layer.
	groupMemberRepository := repositories.NewGroupMemberRepository(nil)
	groupIds, err := groupMemberRepository.GetGroupIdsByUserId(utils.UintToString(claims.UserId))
	if err != nil {
		return nil, nil, err
	}

	receiptRepository := repositories.NewReceiptRepository(nil)
	receipts, err := receiptRepository.SearchReceiptsByGroupIds(groupIds, in.Query, limit)
	if err != nil {
		return nil, nil, err
	}

	results := make([]structs.SearchResult, 0, len(receipts))
	for _, receipt := range receipts {
		results = append(results, structs.SearchResult{
			ID:            receipt.ID,
			GroupID:       receipt.GroupId,
			Name:          receipt.Name,
			Date:          receipt.Date,
			Type:          "Receipt",
			Amount:        receipt.Amount,
			ReceiptStatus: receipt.Status,
			PaidByUserId:  receipt.PaidByUserID,
			CreatedAt:     receipt.CreatedAt,
		})
	}

	return nil, results, nil
}

func handleGetReceipt(ctx context.Context, req *mcpsdk.CallToolRequest, in getReceiptInput) (*mcpsdk.CallToolResult, any, error) {
	claims, err := claimsFromRequest(req)
	if err != nil {
		return nil, nil, err
	}

	if len(in.Id) == 0 {
		return nil, nil, errors.New("id is required")
	}

	receiptRepository := repositories.NewReceiptRepository(nil)
	receipt, err := receiptRepository.GetFullyLoadedReceiptById(in.Id)
	if err != nil {
		return nil, nil, errors.New("receipt not found")
	}

	// Enforce the same group.receipts.read check the REST handler relies on.
	// Use an identical "not found" error for missing and unauthorized so we
	// don't leak the existence of receipts in other users' groups.
	permissionService := services.NewPermissionService(nil)
	hasAccess, err := permissionService.HasGroupPermissions(claims.UserId, receipt.GroupId, permissions.GroupReceiptsRead)
	if err != nil || !hasAccess {
		return nil, nil, errors.New("receipt not found")
	}

	return nil, receipt, nil
}

func handleListGroups(ctx context.Context, req *mcpsdk.CallToolRequest, _ emptyInput) (*mcpsdk.CallToolResult, any, error) {
	claims, err := claimsFromRequest(req)
	if err != nil {
		return nil, nil, err
	}

	groupService := services.NewGroupService(nil)
	groups, err := groupService.GetGroupsForUser(utils.UintToString(claims.UserId))
	if err != nil {
		return nil, nil, err
	}

	return nil, groups, nil
}

func handleListCategories(ctx context.Context, req *mcpsdk.CallToolRequest, _ emptyInput) (*mcpsdk.CallToolResult, any, error) {
	if _, err := claimsFromRequest(req); err != nil {
		return nil, nil, err
	}

	categoryRepository := repositories.NewCategoryRepository(nil)
	categories, err := categoryRepository.GetAllCategories("*")
	if err != nil {
		return nil, nil, err
	}

	return nil, categories, nil
}

func handleListTags(ctx context.Context, req *mcpsdk.CallToolRequest, _ emptyInput) (*mcpsdk.CallToolResult, any, error) {
	if _, err := claimsFromRequest(req); err != nil {
		return nil, nil, err
	}

	tagsRepository := repositories.NewTagsRepository(nil)
	tags, err := tagsRepository.GetAllTags("*")
	if err != nil {
		return nil, nil, err
	}

	return nil, tags, nil
}

func handleListDashboards(ctx context.Context, req *mcpsdk.CallToolRequest, in listDashboardsInput) (*mcpsdk.CallToolResult, any, error) {
	claims, err := claimsFromRequest(req)
	if err != nil {
		return nil, nil, err
	}

	if len(in.GroupId) == 0 {
		return nil, nil, errors.New("groupId is required")
	}

	groupId, err := utils.StringToUint(in.GroupId)
	if err != nil {
		return nil, nil, errors.New("invalid groupId")
	}

	permissionService := services.NewPermissionService(nil)
	hasAccess, err := permissionService.HasGroupPermissions(claims.UserId, groupId, permissions.GroupDashboardsRead)
	if err != nil || !hasAccess {
		return nil, nil, errors.New("unauthorized to access this group")
	}

	dashboardRepository := repositories.NewDashboardRepository(nil)
	dashboards, err := dashboardRepository.GetDashboardsForUserByGroup(claims.UserId, groupId)
	if err != nil {
		return nil, nil, err
	}

	return nil, dashboards, nil
}

// claimsFromRequest retrieves the verified user claims that verifyToken stashed
// on the bearer TokenInfo.
func claimsFromRequest(req *mcpsdk.CallToolRequest) (*structs.Claims, error) {
	if req.Extra == nil || req.Extra.TokenInfo == nil {
		return nil, errors.New("missing authentication context")
	}

	claims, ok := req.Extra.TokenInfo.Extra[claimsKey].(*structs.Claims)
	if !ok || claims == nil {
		return nil, errors.New("invalid authentication context")
	}

	return claims, nil
}
