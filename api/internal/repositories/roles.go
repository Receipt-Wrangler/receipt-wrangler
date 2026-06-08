package repositories

import (
	"errors"
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

	appRoleCounts, err := repository.countAppRoleAssignments()
	if err != nil {
		return nil, err
	}

	groupRoleCounts, err := repository.countGroupRoleAssignments()
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
			IsDefault:     role.IsDefault,
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
			IsDefault:     role.IsDefault,
			IsSystem:      role.IsSystem,
			Permissions:   perms,
			AssignedCount: groupRoleCounts[role.ID],
		})
	}

	return roles, nil
}

// roleAssignmentCount is the row shape returned by the grouped assignment-count
// queries: a role id and the number of users/members assigned to it.
type roleAssignmentCount struct {
	ID    uint
	Count int
}

// countAppRoleAssignments returns a map of app role id -> number of users
// currently assigned that role, in a single grouped query (avoids N+1).
func (repository RoleRepository) countAppRoleAssignments() (map[uint]int, error) {
	db := repository.GetDB()

	var rows []roleAssignmentCount
	err := db.Model(&models.User{}).
		Select("app_role_id as id, count(*) as count").
		Where("app_role_id IS NOT NULL").
		Group("app_role_id").
		Scan(&rows).Error
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
// members currently assigned that role, in a single grouped query.
func (repository RoleRepository) countGroupRoleAssignments() (map[uint]int, error) {
	db := repository.GetDB()

	var rows []roleAssignmentCount
	err := db.Model(&models.GroupMember{}).
		Select("group_role_id as id, count(*) as count").
		Where("group_role_id IS NOT NULL").
		Group("group_role_id").
		Scan(&rows).Error
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

// GetAppRolePermissions returns the permission strings granted by an app role.
func (repository RoleRepository) GetAppRolePermissions(appRoleId uint) ([]string, error) {
	db := repository.GetDB()

	perms := make([]string, 0)
	err := db.Model(&models.AppRolePermission{}).
		Where("app_role_id = ?", appRoleId).
		Pluck("permission", &perms).Error
	if err != nil {
		return nil, err
	}

	return perms, nil
}

// GetGroupRolePermissions returns the permission strings granted by a group role.
func (repository RoleRepository) GetGroupRolePermissions(groupRoleId uint) ([]string, error) {
	db := repository.GetDB()

	perms := make([]string, 0)
	err := db.Model(&models.GroupRolePermission{}).
		Where("group_role_id = ?", groupRoleId).
		Pluck("permission", &perms).Error
	if err != nil {
		return nil, err
	}

	return perms, nil
}

// GetUserAppRoleId returns the app role id assigned to a user, or nil when the
// user has no app role. Returns gorm.ErrRecordNotFound if the user is missing.
func (repository RoleRepository) GetUserAppRoleId(userId uint) (*uint, error) {
	db := repository.GetDB()

	var user models.User
	err := db.Select("id", "app_role_id").Where("id = ?", userId).First(&user).Error
	if err != nil {
		return nil, err
	}

	return user.AppRoleID, nil
}

// GetGroupMemberRoleId returns the group role id assigned to a user within a
// group, or nil when the membership has no group role. Returns
// gorm.ErrRecordNotFound if the user is not a member of the group.
func (repository RoleRepository) GetGroupMemberRoleId(userId uint, groupId uint) (*uint, error) {
	db := repository.GetDB()

	var member models.GroupMember
	err := db.Select("user_id", "group_id", "group_role_id").
		Where("user_id = ? AND group_id = ?", userId, groupId).
		First(&member).Error
	if err != nil {
		return nil, err
	}

	return member.GroupRoleID, nil
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

// SetDefaultAppRole makes the given app role the single default: it clears the
// flag on every other app role, then sets it on id. Callers run this inside a
// transaction so the "exactly one default" invariant holds atomically.
func (repository RoleRepository) SetDefaultAppRole(id uint) error {
	db := repository.GetDB()

	err := db.Model(&models.AppRole{}).Where("id <> ?", id).Update("is_default", false).Error
	if err != nil {
		return err
	}

	return db.Model(&models.AppRole{}).Where("id = ?", id).Update("is_default", true).Error
}

// SetDefaultGroupRole makes the given group role the single default: it clears
// the flag on every other group role, then sets it on id.
func (repository RoleRepository) SetDefaultGroupRole(id uint) error {
	db := repository.GetDB()

	err := db.Model(&models.GroupRoleDefinition{}).Where("id <> ?", id).Update("is_default", false).Error
	if err != nil {
		return err
	}

	return db.Model(&models.GroupRoleDefinition{}).Where("id = ?", id).Update("is_default", true).Error
}

// GetAppRoleIdByName returns the id of the app role with the given name, or nil
// when no such role exists (e.g. an unseeded test database).
func (repository RoleRepository) GetAppRoleIdByName(name string) (*uint, error) {
	db := repository.GetDB()

	var role models.AppRole
	err := db.Select("id").Where("name = ?", name).First(&role).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	return &role.ID, nil
}

// GetDefaultAppRoleId returns the id of the default app role, or nil when none
// is flagged default (e.g. an unseeded test database).
func (repository RoleRepository) GetDefaultAppRoleId() (*uint, error) {
	db := repository.GetDB()

	var role models.AppRole
	err := db.Select("id").Where("is_default = ?", true).First(&role).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	return &role.ID, nil
}

// GetDefaultGroupRoleId returns the id of the default group role, or nil when
// none is flagged default (e.g. an unseeded test database).
func (repository RoleRepository) GetDefaultGroupRoleId() (*uint, error) {
	db := repository.GetDB()

	var role models.GroupRoleDefinition
	err := db.Select("id").Where("is_default = ?", true).First(&role).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	return &role.ID, nil
}
