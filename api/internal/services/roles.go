package services

import (
	"errors"
	"receipt-wrangler/api/internal/commands"
	"receipt-wrangler/api/internal/models"
	"receipt-wrangler/api/internal/permissions"
	"receipt-wrangler/api/internal/repositories"
	"receipt-wrangler/api/internal/structs"

	"gorm.io/gorm"
)

var (
	ErrRoleNotFound          = errors.New("role not found")
	ErrRoleTypeMismatch      = errors.New("role type cannot be changed")
	ErrSystemRoleImmutable   = errors.New("system roles cannot be modified")
	ErrSystemRoleUndeletable = errors.New("system roles cannot be deleted")
	ErrRoleAssigned          = errors.New("role is assigned and cannot be deleted")
	ErrRoleIsDefault         = errors.New("the default role cannot be deleted")
	ErrInvalidGrant          = errors.New("a category, tag, paid-by, or report-template grant references a non-existent category, tag, user, or template")
)

type RoleService struct {
	BaseService
}

func NewRoleService(tx *gorm.DB) RoleService {
	service := RoleService{BaseService: BaseService{
		DB: repositories.GetDB(),
		TX: tx,
	}}
	return service
}

func (service RoleService) CreateRole(command commands.UpsertRoleCommand) (structs.RoleView, error) {
	var roleView structs.RoleView

	perms := command.Permissions
	if perms == nil {
		perms = []string{}
	}
	categoryGrants := normalizeUintSlice(command.CategoryGrants)
	tagGrants := normalizeUintSlice(command.TagGrants)
	paidByUserGrants := normalizeUintSlice(command.PaidByUserGrants)
	reportTemplateGrants := normalizeReportTemplateGrants(command.ReportTemplateGrants)

	err := repositories.GetDB().Transaction(func(tx *gorm.DB) error {
		roleRepository := repositories.NewRoleRepository(tx)

		if command.Scope == permissions.ScopeApp {
			role, txErr := roleRepository.CreateAppRole(command.Name, command.Description, perms)
			if txErr != nil {
				return txErr
			}

			roleView = structs.RoleView{
				Id:                   role.ID,
				Name:                 role.Name,
				Description:          role.Description,
				Scope:                permissions.ScopeApp,
				IsSystem:             role.IsSystem,
				Permissions:          perms,
				CategoryGrants:       []uint{},
				TagGrants:            []uint{},
				PaidByUserGrants:     []uint{},
				ReportTemplateGrants: []structs.ReportTemplateGrantView{},
			}

			return nil
		}

		if txErr := validateGrantsExist(tx, categoryGrants, tagGrants, paidByUserGrants, distinctReportTemplateIds(reportTemplateGrants)); txErr != nil {
			return txErr
		}

		role, txErr := roleRepository.CreateGroupRole(command.Name, command.Description, perms, categoryGrants, tagGrants, paidByUserGrants, command.IncludeOwnPaidReceipts, command.SeesAllMembers)
		if txErr != nil {
			return txErr
		}

		if txErr := roleRepository.ReplaceGroupRoleReportTemplateGrants(role.ID, reportTemplateGrants); txErr != nil {
			return txErr
		}

		if txErr := roleRepository.SetGroupRoleIndividualGrantConfig(
			role.ID,
			command.RequiresIndividualCategoryGrants,
			command.RequiresIndividualTagGrants,
		); txErr != nil {
			return txErr
		}

		roleView = structs.RoleView{
			Id:                               role.ID,
			Name:                             role.Name,
			Description:                      role.Description,
			Scope:                            permissions.ScopeGroup,
			IsSystem:                         role.IsSystem,
			Permissions:                      perms,
			CategoryGrants:                   categoryGrants,
			TagGrants:                        tagGrants,
			PaidByUserGrants:                 paidByUserGrants,
			IncludeOwnPaidReceipts:           command.IncludeOwnPaidReceipts,
			SeesAllMembers:                   command.SeesAllMembers,
			RequiresIndividualCategoryGrants: command.RequiresIndividualCategoryGrants,
			RequiresIndividualTagGrants:      command.RequiresIndividualTagGrants,
			ReportTemplateGrants:             reportTemplateGrantsToView(reportTemplateGrants),
		}

		return nil
	})
	if err != nil {
		return structs.RoleView{}, err
	}

	return roleView, nil
}

func (service RoleService) UpdateRole(id uint, command commands.UpsertRoleCommand) (structs.RoleView, error) {
	var roleView structs.RoleView

	perms := command.Permissions
	if perms == nil {
		perms = []string{}
	}
	categoryGrants := normalizeUintSlice(command.CategoryGrants)
	tagGrants := normalizeUintSlice(command.TagGrants)
	paidByUserGrants := normalizeUintSlice(command.PaidByUserGrants)
	reportTemplateGrants := normalizeReportTemplateGrants(command.ReportTemplateGrants)

	err := repositories.GetDB().Transaction(func(tx *gorm.DB) error {
		roleRepository := repositories.NewRoleRepository(tx)

		if command.Scope == permissions.ScopeApp {
			existing, txErr := roleRepository.GetAppRoleById(id)
			if errors.Is(txErr, gorm.ErrRecordNotFound) {
				return resolveMissingRoleError(roleRepository, permissions.ScopeGroup, id)
			}
			if txErr != nil {
				return txErr
			}
			if existing.IsSystem {
				return ErrSystemRoleImmutable
			}

			role, txErr := roleRepository.UpdateAppRole(id, command.Name, command.Description, perms)
			if txErr != nil {
				return txErr
			}

			roleView = structs.RoleView{
				Id:                   role.ID,
				Name:                 role.Name,
				Description:          role.Description,
				Scope:                permissions.ScopeApp,
				IsDefault:            role.IsDefault,
				IsSystem:             role.IsSystem,
				Permissions:          perms,
				CategoryGrants:       []uint{},
				TagGrants:            []uint{},
				PaidByUserGrants:     []uint{},
				ReportTemplateGrants: []structs.ReportTemplateGrantView{},
			}

			return nil
		}

		existing, txErr := roleRepository.GetGroupRoleById(id)
		if errors.Is(txErr, gorm.ErrRecordNotFound) {
			return resolveMissingRoleError(roleRepository, permissions.ScopeApp, id)
		}
		if txErr != nil {
			return txErr
		}
		if existing.IsSystem {
			return ErrSystemRoleImmutable
		}

		if txErr := validateGrantsExist(tx, categoryGrants, tagGrants, paidByUserGrants, distinctReportTemplateIds(reportTemplateGrants)); txErr != nil {
			return txErr
		}

		role, txErr := roleRepository.UpdateGroupRole(id, command.Name, command.Description, perms, categoryGrants, tagGrants, paidByUserGrants, command.IncludeOwnPaidReceipts, command.SeesAllMembers)
		if txErr != nil {
			return txErr
		}

		if txErr := roleRepository.ReplaceGroupRoleReportTemplateGrants(id, reportTemplateGrants); txErr != nil {
			return txErr
		}

		if txErr := roleRepository.SetGroupRoleIndividualGrantConfig(
			id,
			command.RequiresIndividualCategoryGrants,
			command.RequiresIndividualTagGrants,
		); txErr != nil {
			return txErr
		}

		roleView = structs.RoleView{
			Id:                               role.ID,
			Name:                             role.Name,
			Description:                      role.Description,
			Scope:                            permissions.ScopeGroup,
			IsDefault:                        role.IsDefault,
			IsSystem:                         role.IsSystem,
			Permissions:                      perms,
			CategoryGrants:                   categoryGrants,
			TagGrants:                        tagGrants,
			PaidByUserGrants:                 paidByUserGrants,
			IncludeOwnPaidReceipts:           command.IncludeOwnPaidReceipts,
			SeesAllMembers:                   command.SeesAllMembers,
			RequiresIndividualCategoryGrants: command.RequiresIndividualCategoryGrants,
			RequiresIndividualTagGrants:      command.RequiresIndividualTagGrants,
			ReportTemplateGrants:             reportTemplateGrantsToView(reportTemplateGrants),
		}

		return nil
	})
	if err != nil {
		return structs.RoleView{}, err
	}

	// A role's permission list just changed; drop its cached permissions so the
	// next permission check resolves the new set.
	clearRolePermissionCache(command.Scope, id)

	// Group roles also carry category/tag grants, which the update just
	// replaced — evict the cached grant sets so enforcement uses the new grants.
	if command.Scope == permissions.ScopeGroup {
		clearGroupRoleGrantCache(id)
	}

	return roleView, nil
}

// normalizeUintSlice returns a non-nil slice so empty grant sets serialize as
// [] rather than null and never alias the caller's nil.
func normalizeUintSlice(ids []uint) []uint {
	if ids == nil {
		return []uint{}
	}
	return ids
}

// validateGrantsExist confirms every category/tag/paid-by-user/report-template
// grant id references a real row. Because UpsertRoleCommand.Validate already
// rejects duplicate grant ids, a matching count means all ids exist. Returns
// ErrInvalidGrant otherwise.
func validateGrantsExist(tx *gorm.DB, categoryGrants []uint, tagGrants []uint, paidByUserGrants []uint, reportTemplateIds []uint) error {
	if len(categoryGrants) > 0 {
		categoryRepository := repositories.NewCategoryRepository(tx)
		count, err := categoryRepository.CountByIds(categoryGrants)
		if err != nil {
			return err
		}
		if int(count) != len(categoryGrants) {
			return ErrInvalidGrant
		}
	}

	if len(tagGrants) > 0 {
		tagsRepository := repositories.NewTagsRepository(tx)
		count, err := tagsRepository.CountByIds(tagGrants)
		if err != nil {
			return err
		}
		if int(count) != len(tagGrants) {
			return ErrInvalidGrant
		}
	}

	if len(paidByUserGrants) > 0 {
		userRepository := repositories.NewUserRepository(tx)
		count, err := userRepository.CountByIds(paidByUserGrants)
		if err != nil {
			return err
		}
		if int(count) != len(paidByUserGrants) {
			return ErrInvalidGrant
		}
	}

	if len(reportTemplateIds) > 0 {
		reportTemplateRepository := repositories.NewReportTemplateRepository(tx)
		count, err := reportTemplateRepository.CountByIds(reportTemplateIds)
		if err != nil {
			return err
		}
		if int(count) != len(reportTemplateIds) {
			return ErrInvalidGrant
		}
	}

	return nil
}

// normalizeReportTemplateGrants returns a non-nil slice so an unset matrix
// serializes as [] and never aliases the caller's nil.
func normalizeReportTemplateGrants(grants []commands.ReportTemplateGrantCommand) []commands.ReportTemplateGrantCommand {
	if grants == nil {
		return []commands.ReportTemplateGrantCommand{}
	}
	return grants
}

// distinctReportTemplateIds returns the unique template ids referenced by a
// matrix, for the existence check. Validate already rejects duplicate template
// ids, so this is defensive.
func distinctReportTemplateIds(grants []commands.ReportTemplateGrantCommand) []uint {
	seen := make(map[uint]bool, len(grants))
	ids := make([]uint, 0, len(grants))
	for _, grant := range grants {
		if seen[grant.ReportTemplateId] {
			continue
		}
		seen[grant.ReportTemplateId] = true
		ids = append(ids, grant.ReportTemplateId)
	}
	return ids
}

// reportTemplateGrantsToView converts validated command matrix rows into the read
// model returned after a create/update (built from the command, mirroring how the
// other grants are echoed back rather than re-read).
func reportTemplateGrantsToView(grants []commands.ReportTemplateGrantCommand) []structs.ReportTemplateGrantView {
	views := make([]structs.ReportTemplateGrantView, 0, len(grants))
	for _, grant := range grants {
		permissions := make([]string, 0, len(grant.Permissions))
		permissions = append(permissions, grant.Permissions...)
		views = append(views, structs.ReportTemplateGrantView{
			ReportTemplateId: grant.ReportTemplateId,
			Permissions:      permissions,
		})
	}
	return views
}

// resolveMissingRoleError distinguishes a genuine not-found from a forbidden
// role-type switch: if the id exists under the other scope, the caller is
// trying to change a role's type, which is not allowed.
func resolveMissingRoleError(roleRepository repositories.RoleRepository, otherScope permissions.Scope, id uint) error {
	var otherErr error
	if otherScope == permissions.ScopeApp {
		_, otherErr = roleRepository.GetAppRoleById(id)
	} else {
		_, otherErr = roleRepository.GetGroupRoleById(id)
	}

	if otherErr == nil {
		return ErrRoleTypeMismatch
	}

	if errors.Is(otherErr, gorm.ErrRecordNotFound) {
		return ErrRoleNotFound
	}

	return otherErr
}

func (service RoleService) GetRoles() ([]structs.RoleView, error) {
	roleRepository := repositories.NewRoleRepository(nil)
	return roleRepository.GetAllRoles()
}

// DeleteRole deletes an app- or group-scoped role. The scope disambiguates the
// id (app and group role ids overlap). System roles cannot be deleted, and a
// role cannot be deleted while it is assigned to any user or group member.
func (service RoleService) DeleteRole(id uint, scope permissions.Scope) error {
	err := repositories.GetDB().Transaction(func(tx *gorm.DB) error {
		roleRepository := repositories.NewRoleRepository(tx)

		if scope == permissions.ScopeApp {
			existing, txErr := roleRepository.GetAppRoleById(id)
			if errors.Is(txErr, gorm.ErrRecordNotFound) {
				return resolveMissingRoleError(roleRepository, permissions.ScopeGroup, id)
			}
			if txErr != nil {
				return txErr
			}
			if existing.IsSystem {
				return ErrSystemRoleUndeletable
			}

			if existing.IsDefault {
				return ErrRoleIsDefault
			}

			count, txErr := roleRepository.CountUsersWithAppRole(id)
			if txErr != nil {
				return txErr
			}
			if count > 0 {
				return ErrRoleAssigned
			}

			return roleRepository.DeleteAppRole(id)
		}

		existing, txErr := roleRepository.GetGroupRoleById(id)
		if errors.Is(txErr, gorm.ErrRecordNotFound) {
			return resolveMissingRoleError(roleRepository, permissions.ScopeApp, id)
		}
		if txErr != nil {
			return txErr
		}
		if existing.IsSystem {
			return ErrSystemRoleUndeletable
		}

		if existing.IsDefault {
			return ErrRoleIsDefault
		}

		count, txErr := roleRepository.CountGroupMembersWithGroupRole(id)
		if txErr != nil {
			return txErr
		}
		if count > 0 {
			return ErrRoleAssigned
		}

		return roleRepository.DeleteGroupRole(id)
	})
	if err != nil {
		return err
	}

	// The role is gone; evict any cached permissions for it.
	clearRolePermissionCache(scope, id)

	// And its cached category/tag grants (group roles only).
	if scope == permissions.ScopeGroup {
		clearGroupRoleGrantCache(id)
	}

	return nil
}

// SetDefaultRole makes the given role the default for its scope (the role
// assigned to new accounts for APP, or to group creators for GROUP), clearing
// the previous default. The scope disambiguates the id (app and group role ids
// overlap). System roles are allowed — the seeded legacy roles are the initial
// defaults. The role's permission list is unchanged, so no cache eviction is
// needed.
func (service RoleService) SetDefaultRole(id uint, scope permissions.Scope) (structs.RoleView, error) {
	var roleView structs.RoleView

	err := repositories.GetDB().Transaction(func(tx *gorm.DB) error {
		roleRepository := repositories.NewRoleRepository(tx)

		if scope == permissions.ScopeApp {
			existing, txErr := roleRepository.GetAppRoleById(id)
			if errors.Is(txErr, gorm.ErrRecordNotFound) {
				return resolveMissingRoleError(roleRepository, permissions.ScopeGroup, id)
			}
			if txErr != nil {
				return txErr
			}

			if txErr := roleRepository.SetDefaultAppRole(id); txErr != nil {
				return txErr
			}

			roleView = appRoleToView(existing, true)
			return nil
		}

		existing, txErr := roleRepository.GetGroupRoleById(id)
		if errors.Is(txErr, gorm.ErrRecordNotFound) {
			return resolveMissingRoleError(roleRepository, permissions.ScopeApp, id)
		}
		if txErr != nil {
			return txErr
		}

		if txErr := roleRepository.SetDefaultGroupRole(id); txErr != nil {
			return txErr
		}

		roleView = groupRoleToView(existing, true)
		return nil
	})
	if err != nil {
		return structs.RoleView{}, err
	}

	return roleView, nil
}

// appRoleToView builds the read model for an app role, overriding IsDefault with
// the post-update value (the loaded row predates the flag change).
func appRoleToView(role models.AppRole, isDefault bool) structs.RoleView {
	perms := make([]string, 0, len(role.Permissions))
	for _, permission := range role.Permissions {
		perms = append(perms, permission.Permission)
	}

	return structs.RoleView{
		Id:                   role.ID,
		Name:                 role.Name,
		Description:          role.Description,
		Scope:                permissions.ScopeApp,
		IsDefault:            isDefault,
		IsSystem:             role.IsSystem,
		Permissions:          perms,
		CategoryGrants:       []uint{},
		TagGrants:            []uint{},
		PaidByUserGrants:     []uint{},
		ReportTemplateGrants: []structs.ReportTemplateGrantView{},
	}
}

// groupRoleToView builds the read model for a group role, overriding IsDefault
// with the post-update value. Grant ids are read from the role's preloaded
// CategoryGrants/TagGrants (GetGroupRoleById preloads them).
func groupRoleToView(role models.GroupRoleDefinition, isDefault bool) structs.RoleView {
	perms := make([]string, 0, len(role.Permissions))
	for _, permission := range role.Permissions {
		perms = append(perms, permission.Permission)
	}

	categoryGrants := make([]uint, 0, len(role.CategoryGrants))
	for _, grant := range role.CategoryGrants {
		categoryGrants = append(categoryGrants, grant.CategoryID)
	}

	tagGrants := make([]uint, 0, len(role.TagGrants))
	for _, grant := range role.TagGrants {
		tagGrants = append(tagGrants, grant.TagID)
	}

	paidByUserGrants := make([]uint, 0, len(role.PaidByUserGrants))
	for _, grant := range role.PaidByUserGrants {
		paidByUserGrants = append(paidByUserGrants, grant.UserID)
	}

	return structs.RoleView{
		Id:                               role.ID,
		Name:                             role.Name,
		Description:                      role.Description,
		Scope:                            permissions.ScopeGroup,
		IsDefault:                        isDefault,
		IsSystem:                         role.IsSystem,
		Permissions:                      perms,
		CategoryGrants:                   categoryGrants,
		TagGrants:                        tagGrants,
		PaidByUserGrants:                 paidByUserGrants,
		IncludeOwnPaidReceipts:           role.IncludeOwnPaidReceipts,
		SeesAllMembers:                   role.SeesAllMembers,
		RequiresIndividualCategoryGrants: role.RequiresIndividualCategoryGrants,
		RequiresIndividualTagGrants:      role.RequiresIndividualTagGrants,
		ReportTemplateGrants:             repositories.ReportTemplateGrantsFromRole(role),
	}
}
