package services

import (
	"errors"
	"receipt-wrangler/api/internal/commands"
	"receipt-wrangler/api/internal/permissions"
	"receipt-wrangler/api/internal/repositories"
	"receipt-wrangler/api/internal/structs"

	"gorm.io/gorm"
)

var (
	ErrRoleNotFound        = errors.New("role not found")
	ErrRoleTypeMismatch    = errors.New("role type cannot be changed")
	ErrSystemRoleImmutable = errors.New("system roles cannot be modified")
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

	err := repositories.GetDB().Transaction(func(tx *gorm.DB) error {
		roleRepository := repositories.NewRoleRepository(tx)

		if command.Scope == permissions.ScopeApp {
			role, txErr := roleRepository.CreateAppRole(command.Name, command.Description, perms)
			if txErr != nil {
				return txErr
			}

			roleView = structs.RoleView{
				Id:          role.ID,
				Name:        role.Name,
				Description: role.Description,
				Scope:       permissions.ScopeApp,
				IsSystem:    role.IsSystem,
				Permissions: perms,
			}

			return nil
		}

		role, txErr := roleRepository.CreateGroupRole(command.Name, command.Description, perms)
		if txErr != nil {
			return txErr
		}

		roleView = structs.RoleView{
			Id:          role.ID,
			Name:        role.Name,
			Description: role.Description,
			Scope:       permissions.ScopeGroup,
			IsSystem:    role.IsSystem,
			Permissions: perms,
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
				Id:          role.ID,
				Name:        role.Name,
				Description: role.Description,
				Scope:       permissions.ScopeApp,
				IsSystem:    role.IsSystem,
				Permissions: perms,
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

		role, txErr := roleRepository.UpdateGroupRole(id, command.Name, command.Description, perms)
		if txErr != nil {
			return txErr
		}

		roleView = structs.RoleView{
			Id:          role.ID,
			Name:        role.Name,
			Description: role.Description,
			Scope:       permissions.ScopeGroup,
			IsSystem:    role.IsSystem,
			Permissions: perms,
		}

		return nil
	})
	if err != nil {
		return structs.RoleView{}, err
	}

	return roleView, nil
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
