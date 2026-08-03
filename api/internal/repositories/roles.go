package repositories

import (
	"errors"
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

func (repository RoleRepository) CreateAppRole(name string, description string, perms []string, skipDefaultGroupCreation bool) (models.AppRole, error) {
	db := repository.GetDB()

	rolePermissions := make([]models.AppRolePermission, 0, len(perms))
	for _, permission := range perms {
		rolePermissions = append(rolePermissions, models.AppRolePermission{Permission: permission})
	}

	role := models.AppRole{
		Name:                     name,
		Description:              description,
		SkipDefaultGroupCreation: skipDefaultGroupCreation,
		Permissions:              rolePermissions,
	}

	err := db.Create(&role).Error
	if err != nil {
		return models.AppRole{}, err
	}

	return role, nil
}

func (repository RoleRepository) CreateGroupRole(name string, description string, perms []string, categoryGrantIds []uint, tagGrantIds []uint, paidByUserGrantIds []uint, includeOwnPaidReceipts bool, seesAllMembers bool) (models.GroupRoleDefinition, error) {
	db := repository.GetDB()

	rolePermissions := make([]models.GroupRolePermission, 0, len(perms))
	for _, permission := range perms {
		rolePermissions = append(rolePermissions, models.GroupRolePermission{Permission: permission})
	}

	role := models.GroupRoleDefinition{
		Name:                       name,
		Description:                description,
		IncludeOwnPaidReceipts:     includeOwnPaidReceipts,
		PaidByVisibilityRestricted: includeOwnPaidReceipts || len(paidByUserGrantIds) > 0,
		SeesAllMembers:             seesAllMembers,
		Permissions:                rolePermissions,
	}

	err := db.Create(&role).Error
	if err != nil {
		return models.GroupRoleDefinition{}, err
	}

	err = repository.replaceGroupRoleGrants(role.ID, categoryGrantIds, tagGrantIds, paidByUserGrantIds)
	if err != nil {
		return models.GroupRoleDefinition{}, err
	}

	return repository.GetGroupRoleById(role.ID)
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
	err := db.Preload("Permissions").
		Preload("CategoryGrants").
		Preload("TagGrants").
		Preload("PaidByUserGrants").
		Preload("ReportTemplateGrants").
		First(&role, id).Error
	if err != nil {
		return models.GroupRoleDefinition{}, err
	}

	return role, nil
}

func (repository RoleRepository) UpdateAppRole(id uint, name string, description string, perms []string, skipDefaultGroupCreation bool) (models.AppRole, error) {
	db := repository.GetDB()

	err := db.Where("app_role_id = ?", id).Delete(&models.AppRolePermission{}).Error
	if err != nil {
		return models.AppRole{}, err
	}

	// Use the map form so false bools persist (GORM's struct Updates skips
	// zero-value bools, which would leave a toggled-off flag set).
	err = db.Model(&models.AppRole{}).Where("id = ?", id).Updates(map[string]interface{}{
		"name":                        name,
		"description":                 description,
		"skip_default_group_creation": skipDefaultGroupCreation,
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

func (repository RoleRepository) UpdateGroupRole(id uint, name string, description string, perms []string, categoryGrantIds []uint, tagGrantIds []uint, paidByUserGrantIds []uint, includeOwnPaidReceipts bool, seesAllMembers bool) (models.GroupRoleDefinition, error) {
	db := repository.GetDB()

	err := db.Where("group_role_id = ?", id).Delete(&models.GroupRolePermission{}).Error
	if err != nil {
		return models.GroupRoleDefinition{}, err
	}

	// Use the map form so false bools persist (GORM's struct Updates skips
	// zero-value bools, which would leave a toggled-off flag set).
	err = db.Model(&models.GroupRoleDefinition{}).Where("id = ?", id).Updates(map[string]interface{}{
		"name":                          name,
		"description":                   description,
		"include_own_paid_receipts":     includeOwnPaidReceipts,
		"paid_by_visibility_restricted": includeOwnPaidReceipts || len(paidByUserGrantIds) > 0,
		"sees_all_members":              seesAllMembers,
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

	err = repository.replaceGroupRoleGrants(id, categoryGrantIds, tagGrantIds, paidByUserGrantIds)
	if err != nil {
		return models.GroupRoleDefinition{}, err
	}

	return repository.GetGroupRoleById(id)
}

// replaceGroupRoleGrants resets a group role's category, tag, and paid-by user
// grants to exactly the given id sets (delete-all-then-insert, mirroring the
// permission sync). The nested Category/Tag/User belongs-to associations are
// Omit-ted so GORM never tries to upsert a zero-valued row from the grant rows —
// only the join rows are written.
func (repository RoleRepository) replaceGroupRoleGrants(groupRoleId uint, categoryGrantIds []uint, tagGrantIds []uint, paidByUserGrantIds []uint) error {
	db := repository.GetDB()

	err := db.Where("group_role_id = ?", groupRoleId).Delete(&models.GroupRoleCategoryGrant{}).Error
	if err != nil {
		return err
	}

	err = db.Where("group_role_id = ?", groupRoleId).Delete(&models.GroupRoleTagGrant{}).Error
	if err != nil {
		return err
	}

	err = db.Where("group_role_id = ?", groupRoleId).Delete(&models.GroupRolePaidByUserGrant{}).Error
	if err != nil {
		return err
	}

	if len(categoryGrantIds) > 0 {
		categoryGrants := make([]models.GroupRoleCategoryGrant, 0, len(categoryGrantIds))
		for _, categoryId := range categoryGrantIds {
			categoryGrants = append(categoryGrants, models.GroupRoleCategoryGrant{GroupRoleID: groupRoleId, CategoryID: categoryId})
		}

		err = db.Omit("Category").Create(&categoryGrants).Error
		if err != nil {
			return err
		}
	}

	if len(tagGrantIds) > 0 {
		tagGrants := make([]models.GroupRoleTagGrant, 0, len(tagGrantIds))
		for _, tagId := range tagGrantIds {
			tagGrants = append(tagGrants, models.GroupRoleTagGrant{GroupRoleID: groupRoleId, TagID: tagId})
		}

		err = db.Omit("Tag").Create(&tagGrants).Error
		if err != nil {
			return err
		}
	}

	if len(paidByUserGrantIds) > 0 {
		paidByUserGrants := make([]models.GroupRolePaidByUserGrant, 0, len(paidByUserGrantIds))
		for _, userId := range paidByUserGrantIds {
			paidByUserGrants = append(paidByUserGrants, models.GroupRolePaidByUserGrant{GroupRoleID: groupRoleId, UserID: userId})
		}

		err = db.Omit("User").Create(&paidByUserGrants).Error
		if err != nil {
			return err
		}
	}

	return nil
}

// ReplaceGroupRoleReportTemplateGrants resets a group role's report-template grants
// to exactly the given matrix (delete-all-then-insert, one row per template+action),
// and records whether the role is now restricted so it fails closed once its last
// granted template is deleted. Kept a separate method from replaceGroupRoleGrants so
// the widely-called CreateGroupRole/UpdateGroupRole signatures stay unchanged. The
// ReportTemplate belongs-to association is Omit-ted so only the join rows are written.
func (repository RoleRepository) ReplaceGroupRoleReportTemplateGrants(groupRoleId uint, grants []commands.ReportTemplateGrantCommand) error {
	db := repository.GetDB()

	err := db.Where("group_role_id = ?", groupRoleId).Delete(&models.GroupRoleReportTemplateGrant{}).Error
	if err != nil {
		return err
	}

	// Map-style Update so the false value persists when a role is cleared of grants.
	err = db.Model(&models.GroupRoleDefinition{}).Where("id = ?", groupRoleId).
		Update("report_template_grants_restricted", len(grants) > 0).Error
	if err != nil {
		return err
	}

	if len(grants) == 0 {
		return nil
	}

	rows := make([]models.GroupRoleReportTemplateGrant, 0, len(grants))
	for _, grant := range grants {
		for _, permission := range grant.Permissions {
			rows = append(rows, models.GroupRoleReportTemplateGrant{
				GroupRoleID:      groupRoleId,
				ReportTemplateID: grant.ReportTemplateId,
				Permission:       permission,
			})
		}
	}

	return db.Omit("ReportTemplate").Create(&rows).Error
}

func (repository RoleRepository) GetAllRoles() ([]structs.RoleView, error) {
	db := repository.GetDB()

	var appRoles []models.AppRole
	err := db.Preload("Permissions").Find(&appRoles).Error
	if err != nil {
		return nil, err
	}

	var groupRoles []models.GroupRoleDefinition
	err = db.Preload("Permissions").
		Preload("CategoryGrants").
		Preload("TagGrants").
		Preload("PaidByUserGrants").
		Preload("ReportTemplateGrants").
		Find(&groupRoles).Error
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
			Id:                       role.ID,
			Name:                     role.Name,
			Description:              role.Description,
			Scope:                    permissions.ScopeApp,
			IsDefault:                role.IsDefault,
			IsSystem:                 role.IsSystem,
			Permissions:              perms,
			AssignedCount:            appRoleCounts[role.ID],
			SkipDefaultGroupCreation: role.SkipDefaultGroupCreation,
			CategoryGrants:           []uint{},
			TagGrants:                []uint{},
			PaidByUserGrants:         []uint{},
			ReportTemplateGrants:     []structs.ReportTemplateGrantView{},
		})
	}

	for _, role := range groupRoles {
		perms := make([]string, 0, len(role.Permissions))
		for _, permission := range role.Permissions {
			perms = append(perms, permission.Permission)
		}

		roles = append(roles, structs.RoleView{
			Id:                     role.ID,
			Name:                   role.Name,
			Description:            role.Description,
			Scope:                  permissions.ScopeGroup,
			IsDefault:              role.IsDefault,
			IsSystem:               role.IsSystem,
			Permissions:            perms,
			AssignedCount:          groupRoleCounts[role.ID],
			CategoryGrants:         categoryGrantIdsFromRole(role),
			TagGrants:              tagGrantIdsFromRole(role),
			PaidByUserGrants:       paidByUserGrantIdsFromRole(role),
			IncludeOwnPaidReceipts: role.IncludeOwnPaidReceipts,
			SeesAllMembers:         role.SeesAllMembers,
			ReportTemplateGrants:   ReportTemplateGrantsFromRole(role),
		})
	}

	return roles, nil
}

// categoryGrantIdsFromRole extracts the category-grant ids from a loaded group
// role's preloaded CategoryGrants, normalized to a non-nil slice.
func categoryGrantIdsFromRole(role models.GroupRoleDefinition) []uint {
	ids := make([]uint, 0, len(role.CategoryGrants))
	for _, grant := range role.CategoryGrants {
		ids = append(ids, grant.CategoryID)
	}
	return ids
}

// tagGrantIdsFromRole extracts the tag-grant ids from a loaded group role's
// preloaded TagGrants, normalized to a non-nil slice.
func tagGrantIdsFromRole(role models.GroupRoleDefinition) []uint {
	ids := make([]uint, 0, len(role.TagGrants))
	for _, grant := range role.TagGrants {
		ids = append(ids, grant.TagID)
	}
	return ids
}

// paidByUserGrantIdsFromRole extracts the paid-by user-grant ids from a loaded
// group role's preloaded PaidByUserGrants, normalized to a non-nil slice.
func paidByUserGrantIdsFromRole(role models.GroupRoleDefinition) []uint {
	ids := make([]uint, 0, len(role.PaidByUserGrants))
	for _, grant := range role.PaidByUserGrants {
		ids = append(ids, grant.UserID)
	}
	return ids
}

// ReportTemplateGrantsFromRole groups a loaded group role's flat report-template
// grant rows (one per template+action) into the per-template matrix rows the API
// exposes, preserving first-seen template order for a stable response. Exported so
// the role service's view builders reuse the grouping rather than duplicating it.
func ReportTemplateGrantsFromRole(role models.GroupRoleDefinition) []structs.ReportTemplateGrantView {
	byTemplate := make(map[uint][]string)
	order := make([]uint, 0)
	for _, grant := range role.ReportTemplateGrants {
		if _, seen := byTemplate[grant.ReportTemplateID]; !seen {
			order = append(order, grant.ReportTemplateID)
		}
		byTemplate[grant.ReportTemplateID] = append(byTemplate[grant.ReportTemplateID], grant.Permission)
	}

	views := make([]structs.ReportTemplateGrantView, 0, len(order))
	for _, templateId := range order {
		views = append(views, structs.ReportTemplateGrantView{
			ReportTemplateId: templateId,
			Permissions:      byTemplate[templateId],
		})
	}
	return views
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

// appRoleIdsWithPermission returns the ids of every app role whose granted
// permissions satisfy perm, honoring wildcard grants (e.g. "*", "app.*") via the
// permission matcher — something a raw SQL equality check on the permission
// strings could not do.
//
// TODO: this scans every app role in memory to run the wildcard-aware match. App
// roles are a small, admin-defined set and the callers (the last-admin delete
// guard and the first-admin-login check) are infrequent, so this is fine today.
// If app-role counts ever grow large, cache the reverse permission->roleIds
// lookup (invalidated on role changes, mirroring services/permission_cache.go)
// instead of scanning on every call.
func (repository RoleRepository) appRoleIdsWithPermission(perm string) ([]uint, error) {
	db := repository.GetDB()

	var appRoles []models.AppRole
	err := db.Preload("Permissions").Find(&appRoles).Error
	if err != nil {
		return nil, err
	}

	roleIds := make([]uint, 0)
	for _, role := range appRoles {
		granted := make([]string, 0, len(role.Permissions))
		for _, permission := range role.Permissions {
			granted = append(granted, permission.Permission)
		}
		if permissions.HasAll(granted, perm) {
			roleIds = append(roleIds, role.ID)
		}
	}

	return roleIds, nil
}

// CountUsersWithAppPermission returns the number of users whose assigned app role
// grants perm (wildcard-aware via appRoleIdsWithPermission). It defines the
// "admin" population in place of the removed legacy UserRole enum — e.g. counting
// holders of app.users.read to guard against deleting the last administrator.
func (repository RoleRepository) CountUsersWithAppPermission(perm string) (int64, error) {
	db := repository.GetDB()

	roleIds, err := repository.appRoleIdsWithPermission(perm)
	if err != nil {
		return 0, err
	}
	if len(roleIds) == 0 {
		return 0, nil
	}

	var count int64
	err = db.Model(&models.User{}).Where("app_role_id IN ?", roleIds).Count(&count).Error
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

// GetGroupRoleCategoryIds returns the category ids a group role grants its
// members. An empty result means the role is unrestricted (see all categories).
func (repository RoleRepository) GetGroupRoleCategoryIds(groupRoleId uint) ([]uint, error) {
	db := repository.GetDB()

	ids := make([]uint, 0)
	err := db.Model(&models.GroupRoleCategoryGrant{}).
		Where("group_role_id = ?", groupRoleId).
		Pluck("category_id", &ids).Error
	if err != nil {
		return nil, err
	}

	return ids, nil
}

// GetGroupRoleTagIds returns the tag ids a group role grants its members. An
// empty result means the role is unrestricted (see all tags).
func (repository RoleRepository) GetGroupRoleTagIds(groupRoleId uint) ([]uint, error) {
	db := repository.GetDB()

	ids := make([]uint, 0)
	err := db.Model(&models.GroupRoleTagGrant{}).
		Where("group_role_id = ?", groupRoleId).
		Pluck("tag_id", &ids).Error
	if err != nil {
		return nil, err
	}

	return ids, nil
}

// GetGroupRolePaidByUserIds returns the user ids whose receipts a group role
// lets its members see (the absolute paid-by grants only — the relative "their
// own" token is on GetGroupRolePaidByConfig). An empty result with include-own
// false means the role is unrestricted (members see every payer).
func (repository RoleRepository) GetGroupRolePaidByUserIds(groupRoleId uint) ([]uint, error) {
	db := repository.GetDB()

	ids := make([]uint, 0)
	err := db.Model(&models.GroupRolePaidByUserGrant{}).
		Where("group_role_id = ?", groupRoleId).
		Pluck("user_id", &ids).Error
	if err != nil {
		return nil, err
	}

	return ids, nil
}

// GetGroupRolePaidByConfig returns a group role's two scalar paid-by flags in a
// single row read: includeOwn (the relative "their own receipts" token) and
// restricted (whether the role opted into paid-by filtering at all — it stays
// true after a granted user is deleted and the grant rows cascade away, so a
// configured role keeps failing closed instead of widening to see-all).
func (repository RoleRepository) GetGroupRolePaidByConfig(groupRoleId uint) (includeOwn bool, restricted bool, err error) {
	db := repository.GetDB()

	var role models.GroupRoleDefinition
	err = db.Select("include_own_paid_receipts", "paid_by_visibility_restricted").
		Where("id = ?", groupRoleId).
		First(&role).Error
	if err != nil {
		return false, false, err
	}

	return role.IncludeOwnPaidReceipts, role.PaidByVisibilityRestricted, nil
}

// GetGroupRoleIndividualGrantConfig returns whether a group role requires
// per-member category / tag assignment. When set, a member holding the role with
// no membership grants of their own sees nothing at all rather than falling back
// to the role's set — so forgetting to assign a new member fails closed.
func (repository RoleRepository) GetGroupRoleIndividualGrantConfig(groupRoleId uint) (categories bool, tags bool, err error) {
	db := repository.GetDB()

	var role models.GroupRoleDefinition
	err = db.Select("requires_individual_category_grants", "requires_individual_tag_grants").
		Where("id = ?", groupRoleId).
		First(&role).Error
	if err != nil {
		return false, false, err
	}

	return role.RequiresIndividualCategoryGrants, role.RequiresIndividualTagGrants, nil
}

// SetGroupRoleIndividualGrantConfig records whether a group role requires
// per-member category / tag assignment.
//
// Kept a separate method from CreateGroupRole / UpdateGroupRole — whose positional
// signatures already end in two bools — for the same reason
// ReplaceGroupRoleReportTemplateGrants is separate: appending two more would make
// four adjacent bools at every call site, which is trivially transposable and
// silently wrong when it is. Uses the map form so a toggled-off false persists.
func (repository RoleRepository) SetGroupRoleIndividualGrantConfig(groupRoleId uint, requiresCategories bool, requiresTags bool) error {
	return repository.GetDB().Model(&models.GroupRoleDefinition{}).
		Where("id = ?", groupRoleId).
		Updates(map[string]interface{}{
			"requires_individual_category_grants": requiresCategories,
			"requires_individual_tag_grants":      requiresTags,
		}).Error
}

// GetGroupRoleReportTemplateGrants returns a group role's report-template grant
// rows (one per template+action). An empty result means the role is unrestricted
// (act on every template its group access reaches) unless
// GetGroupRoleReportTemplateGrantsRestricted reports it opted into restriction.
func (repository RoleRepository) GetGroupRoleReportTemplateGrants(groupRoleId uint) ([]models.GroupRoleReportTemplateGrant, error) {
	db := repository.GetDB()

	grants := make([]models.GroupRoleReportTemplateGrant, 0)
	err := db.Where("group_role_id = ?", groupRoleId).Find(&grants).Error
	if err != nil {
		return nil, err
	}

	return grants, nil
}

// GetGroupRoleReportTemplateGrantsRestricted reports whether a group role opted
// into report-template restriction. Like the paid-by flag, it stays true after the
// grant rows cascade away (last granted template deleted), so a configured role
// keeps failing closed instead of widening back to see-all.
func (repository RoleRepository) GetGroupRoleReportTemplateGrantsRestricted(groupRoleId uint) (bool, error) {
	db := repository.GetDB()

	var role models.GroupRoleDefinition
	err := db.Select("report_template_grants_restricted").
		Where("id = ?", groupRoleId).
		First(&role).Error
	if err != nil {
		return false, err
	}

	return role.ReportTemplateGrantsRestricted, nil
}

// GetUserAppRoleId returns the app role id assigned to a user, or nil when the
// user has no app role. Returns gorm.ErrRecordNotFound if the user is missing.
// AppRoleSkipsDefaultGroup reports whether users created with this app role skip
// the personal "My Receipts" group. Reads the single column so user creation does
// not preload the role's whole permission set.
func (repository RoleRepository) AppRoleSkipsDefaultGroup(id uint) (bool, error) {
	db := repository.GetDB()

	var role models.AppRole
	err := db.Select("skip_default_group_creation").
		Where("id = ?", id).
		First(&role).Error
	if err != nil {
		return false, err
	}

	return role.SkipDefaultGroupCreation, nil
}

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
