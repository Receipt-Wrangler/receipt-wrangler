package repositories

import (
	"receipt-wrangler/api/internal/commands"
	"receipt-wrangler/api/internal/models"
	"receipt-wrangler/api/internal/permissions"
	"receipt-wrangler/api/internal/structs"

	"gorm.io/gorm"
)

type RoleRepository struct {
	BaseRepository
}

func NewRoleRepository(tx *gorm.DB) RoleRepository {
	repository := RoleRepository{BaseRepository: BaseRepository{
		DB: GetDB(),
		TX: tx,
	}}
	return repository
}

func (repository RoleRepository) CreateAppRole(name string, description string, perms []string) (models.AppRole, error) {
	db := repository.GetDB()

	rolePermissions := make([]models.AppRolePermission, 0, len(perms))
	for _, permission := range perms {
		rolePermissions = append(rolePermissions, models.AppRolePermission{Permission: permission})
	}

	role := models.AppRole{
		Name:        name,
		Description: description,
		Permissions: rolePermissions,
	}

	err := db.Create(&role).Error
	if err != nil {
		return models.AppRole{}, err
	}

	return role, nil
}

func (repository RoleRepository) CreateGroupRole(name string, description string, perms []string) (models.GroupRoleDefinition, error) {
	db := repository.GetDB()

	rolePermissions := make([]models.GroupRolePermission, 0, len(perms))
	for _, permission := range perms {
		rolePermissions = append(rolePermissions, models.GroupRolePermission{Permission: permission})
	}

	role := models.GroupRoleDefinition{
		Name:        name,
		Description: description,
		Permissions: rolePermissions,
	}

	err := db.Create(&role).Error
	if err != nil {
		return models.GroupRoleDefinition{}, err
	}

	return role, nil
}

func (repository RoleRepository) GetAppRoleById(id uint) (models.AppRole, error) {
	db := repository.GetDB()

	var role models.AppRole
	err := db.Preload("Permissions").First(&role, id).Error
	if err != nil {
		return models.AppRole{}, err
	}

	return role, nil
}

func (repository RoleRepository) GetGroupRoleById(id uint) (models.GroupRoleDefinition, error) {
	db := repository.GetDB()

	var role models.GroupRoleDefinition
	err := db.Preload("Permissions").First(&role, id).Error
	if err != nil {
		return models.GroupRoleDefinition{}, err
	}

	return role, nil
}

func (repository RoleRepository) UpdateAppRole(id uint, name string, description string, perms []string) (models.AppRole, error) {
	db := repository.GetDB()

	err := db.Where("app_role_id = ?", id).Delete(&models.AppRolePermission{}).Error
	if err != nil {
		return models.AppRole{}, err
	}

	err = db.Model(&models.AppRole{}).Where("id = ?", id).Updates(map[string]interface{}{
		"name":        name,
		"description": description,
	}).Error
	if err != nil {
		return models.AppRole{}, err
	}

	if len(perms) > 0 {
		rolePermissions := make([]models.AppRolePermission, 0, len(perms))
		for _, permission := range perms {
			rolePermissions = append(rolePermissions, models.AppRolePermission{AppRoleID: id, Permission: permission})
		}

		err = db.Create(&rolePermissions).Error
		if err != nil {
			return models.AppRole{}, err
		}
	}

	return repository.GetAppRoleById(id)
}

func (repository RoleRepository) UpdateGroupRole(id uint, name string, description string, perms []string) (models.GroupRoleDefinition, error) {
	db := repository.GetDB()

	err := db.Where("group_role_id = ?", id).Delete(&models.GroupRolePermission{}).Error
	if err != nil {
		return models.GroupRoleDefinition{}, err
	}

	err = db.Model(&models.GroupRoleDefinition{}).Where("id = ?", id).Updates(map[string]interface{}{
		"name":        name,
		"description": description,
	}).Error
	if err != nil {
		return models.GroupRoleDefinition{}, err
	}

	if len(perms) > 0 {
		rolePermissions := make([]models.GroupRolePermission, 0, len(perms))
		for _, permission := range perms {
			rolePermissions = append(rolePermissions, models.GroupRolePermission{GroupRoleID: id, Permission: permission})
		}

		err = db.Create(&rolePermissions).Error
		if err != nil {
			return models.GroupRoleDefinition{}, err
		}
	}

	return repository.GetGroupRoleById(id)
}

func (repository RoleRepository) GetAllRoles() ([]structs.RoleView, error) {
	db := repository.GetDB()

	var appRoles []models.AppRole
	err := db.Preload("Permissions").Find(&appRoles).Error
	if err != nil {
		return nil, err
	}

	var groupRoles []models.GroupRoleDefinition
	err = db.Preload("Permissions").Find(&groupRoles).Error
	if err != nil {
		return nil, err
	}

	appRoleCounts, err := repository.countAppRoleAssignments(nil)
	if err != nil {
		return nil, err
	}

	groupRoleCounts, err := repository.countGroupRoleAssignments(nil)
	if err != nil {
		return nil, err
	}

	roles := make([]structs.RoleView, 0, len(appRoles)+len(groupRoles))

	for _, role := range appRoles {
		perms := make([]string, 0, len(role.Permissions))
		for _, permission := range role.Permissions {
			perms = append(perms, permission.Permission)
		}

		roles = append(roles, structs.RoleView{
			Id:            role.ID,
			Name:          role.Name,
			Description:   role.Description,
			Scope:         permissions.ScopeApp,
			IsSystem:      role.IsSystem,
			Permissions:   perms,
			AssignedCount: appRoleCounts[role.ID],
		})
	}

	for _, role := range groupRoles {
		perms := make([]string, 0, len(role.Permissions))
		for _, permission := range role.Permissions {
			perms = append(perms, permission.Permission)
		}

		roles = append(roles, structs.RoleView{
			Id:            role.ID,
			Name:          role.Name,
			Description:   role.Description,
			Scope:         permissions.ScopeGroup,
			IsSystem:      role.IsSystem,
			Permissions:   perms,
			AssignedCount: groupRoleCounts[role.ID],
		})
	}

	return roles, nil
}

// roleUnionRow is the row shape scanned from the app/group role union derived
// table. Scope is the literal "APP"/"GROUP" tag added per side of the union.
type roleUnionRow struct {
	ID          uint
	Name        string
	Description string
	IsSystem    bool
	Scope       string
}

// roleUnionSQL unions the two role tables into a single result set so a page can
// span both app- and group-scoped roles. The id sequences are independent, so a
// row is identified by the (id, scope) pair.
const roleUnionSQL = "SELECT id, name, description, is_system, 'APP' AS scope FROM app_roles " +
	"UNION ALL " +
	"SELECT id, name, description, is_system, 'GROUP' AS scope FROM group_role_definitions"

// GetPagedRoles returns a page of roles across both scopes, ordered and counted
// at the database level. The page rows are then enriched with their permissions
// and assignment counts in a second pass.
func (repository RoleRepository) GetPagedRoles(command commands.PagedRoleRequestCommand) ([]structs.RoleView, int64, error) {
	db := repository.GetDB()

	orderBy := command.OrderBy
	if !repository.isValidColumn(orderBy) {
		orderBy = "name"
	}

	buildUnionQuery := func() *gorm.DB {
		query := db.Table("(?) as roles", gorm.Expr(roleUnionSQL))
		if command.Filter.Scope != "" {
			query = query.Where("scope = ?", command.Filter.Scope)
		}
		return query
	}

	var count int64
	err := buildUnionQuery().Count(&count).Error
	if err != nil {
		return nil, 0, err
	}

	// Append a deterministic tie-breaker: ordering by the sort column alone leaves
	// tied rows (e.g. an app role and a group role with the same name) in an
	// unspecified order, which can skip or duplicate rows across LIMIT/OFFSET
	// pages. The (scope, id) pair uniquely identifies a row in the union.
	query := repository.Sort(buildUnionQuery(), orderBy, command.SortDirection)
	query = query.Order("scope").Order("id")
	query = query.Scopes(repository.Paginate(command.Page, command.PageSize))

	var rows []roleUnionRow
	err = query.Scan(&rows).Error
	if err != nil {
		return nil, 0, err
	}

	roleViews, err := repository.buildRoleViews(rows)
	if err != nil {
		return nil, 0, err
	}

	return roleViews, count, nil
}

// isValidColumn whitelists the columns of the role union derived table that may
// be used for ordering (guards against unknown-column errors / injection).
func (repository RoleRepository) isValidColumn(orderBy string) bool {
	return orderBy == "name" ||
		orderBy == "description" ||
		orderBy == "is_system" ||
		orderBy == "scope"
}

// buildRoleViews enriches the page's union rows with their permissions and
// assignment counts, preserving the order of the rows.
func (repository RoleRepository) buildRoleViews(rows []roleUnionRow) ([]structs.RoleView, error) {
	appRoleIds := make([]uint, 0, len(rows))
	groupRoleIds := make([]uint, 0, len(rows))
	for _, row := range rows {
		if row.Scope == string(permissions.ScopeGroup) {
			groupRoleIds = append(groupRoleIds, row.ID)
		} else {
			appRoleIds = append(appRoleIds, row.ID)
		}
	}

	appPermissions, err := repository.getAppRolePermissions(appRoleIds)
	if err != nil {
		return nil, err
	}

	groupPermissions, err := repository.getGroupRolePermissions(groupRoleIds)
	if err != nil {
		return nil, err
	}

	appRoleCounts, err := repository.countAppRoleAssignments(appRoleIds)
	if err != nil {
		return nil, err
	}

	groupRoleCounts, err := repository.countGroupRoleAssignments(groupRoleIds)
	if err != nil {
		return nil, err
	}

	roleViews := make([]structs.RoleView, 0, len(rows))
	for _, row := range rows {
		scope := permissions.ScopeApp
		perms := appPermissions[row.ID]
		assignedCount := appRoleCounts[row.ID]
		if row.Scope == string(permissions.ScopeGroup) {
			scope = permissions.ScopeGroup
			perms = groupPermissions[row.ID]
			assignedCount = groupRoleCounts[row.ID]
		}

		if perms == nil {
			perms = []string{}
		}

		roleViews = append(roleViews, structs.RoleView{
			Id:            row.ID,
			Name:          row.Name,
			Description:   row.Description,
			Scope:         scope,
			IsSystem:      row.IsSystem,
			Permissions:   perms,
			AssignedCount: assignedCount,
		})
	}

	return roleViews, nil
}

// getAppRolePermissions returns a map of app role id -> permission strings for
// the given role ids, in a single query.
func (repository RoleRepository) getAppRolePermissions(roleIds []uint) (map[uint][]string, error) {
	result := make(map[uint][]string)
	if len(roleIds) == 0 {
		return result, nil
	}

	var permissionRows []models.AppRolePermission
	err := repository.GetDB().Where("app_role_id IN ?", roleIds).Find(&permissionRows).Error
	if err != nil {
		return nil, err
	}

	for _, row := range permissionRows {
		result[row.AppRoleID] = append(result[row.AppRoleID], row.Permission)
	}

	return result, nil
}

// getGroupRolePermissions returns a map of group role id -> permission strings
// for the given role ids, in a single query.
func (repository RoleRepository) getGroupRolePermissions(roleIds []uint) (map[uint][]string, error) {
	result := make(map[uint][]string)
	if len(roleIds) == 0 {
		return result, nil
	}

	var permissionRows []models.GroupRolePermission
	err := repository.GetDB().Where("group_role_id IN ?", roleIds).Find(&permissionRows).Error
	if err != nil {
		return nil, err
	}

	for _, row := range permissionRows {
		result[row.GroupRoleID] = append(result[row.GroupRoleID], row.Permission)
	}

	return result, nil
}

// roleAssignmentCount is the row shape returned by the grouped assignment-count
// queries: a role id and the number of users/members assigned to it.
type roleAssignmentCount struct {
	ID    uint
	Count int
}

// countAppRoleAssignments returns a map of app role id -> number of users
// currently assigned that role, in a single grouped query (avoids N+1). A nil
// roleIds counts every role; a non-nil slice scopes the count to those ids (so a
// page request doesn't scan the whole users table).
func (repository RoleRepository) countAppRoleAssignments(roleIds []uint) (map[uint]int, error) {
	if roleIds != nil && len(roleIds) == 0 {
		return map[uint]int{}, nil
	}

	query := repository.GetDB().Model(&models.User{}).
		Select("app_role_id as id, count(*) as count").
		Where("app_role_id IS NOT NULL").
		Group("app_role_id")
	if roleIds != nil {
		query = query.Where("app_role_id IN ?", roleIds)
	}

	var rows []roleAssignmentCount
	err := query.Scan(&rows).Error
	if err != nil {
		return nil, err
	}

	counts := make(map[uint]int, len(rows))
	for _, row := range rows {
		counts[row.ID] = row.Count
	}

	return counts, nil
}

// countGroupRoleAssignments returns a map of group role id -> number of group
// members currently assigned that role, in a single grouped query. A nil
// roleIds counts every role; a non-nil slice scopes the count to those ids.
func (repository RoleRepository) countGroupRoleAssignments(roleIds []uint) (map[uint]int, error) {
	if roleIds != nil && len(roleIds) == 0 {
		return map[uint]int{}, nil
	}

	query := repository.GetDB().Model(&models.GroupMember{}).
		Select("group_role_id as id, count(*) as count").
		Where("group_role_id IS NOT NULL").
		Group("group_role_id")
	if roleIds != nil {
		query = query.Where("group_role_id IN ?", roleIds)
	}

	var rows []roleAssignmentCount
	err := query.Scan(&rows).Error
	if err != nil {
		return nil, err
	}

	counts := make(map[uint]int, len(rows))
	for _, row := range rows {
		counts[row.ID] = row.Count
	}

	return counts, nil
}

func (repository RoleRepository) CountUsersWithAppRole(id uint) (int64, error) {
	db := repository.GetDB()

	var count int64
	err := db.Model(&models.User{}).Where("app_role_id = ?", id).Count(&count).Error
	if err != nil {
		return 0, err
	}

	return count, nil
}

func (repository RoleRepository) CountGroupMembersWithGroupRole(id uint) (int64, error) {
	db := repository.GetDB()

	var count int64
	err := db.Model(&models.GroupMember{}).Where("group_role_id = ?", id).Count(&count).Error
	if err != nil {
		return 0, err
	}

	return count, nil
}

// DeleteAppRole deletes the app role and its permissions (children cascade via
// the AppRolePermission OnDelete:CASCADE constraint).
func (repository RoleRepository) DeleteAppRole(id uint) error {
	db := repository.GetDB()
	return db.Delete(&models.AppRole{}, id).Error
}

// DeleteGroupRole deletes the group role and its permissions and resource
// grants (children cascade via their OnDelete:CASCADE constraints).
func (repository RoleRepository) DeleteGroupRole(id uint) error {
	db := repository.GetDB()
	return db.Delete(&models.GroupRoleDefinition{}, id).Error
}
