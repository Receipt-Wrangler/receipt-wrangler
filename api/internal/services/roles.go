package services

import (
	"receipt-wrangler/api/internal/commands"
	"receipt-wrangler/api/internal/permissions"
	"receipt-wrangler/api/internal/repositories"
	"receipt-wrangler/api/internal/structs"

	"gorm.io/gorm"
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

	err := service.GetDB().Transaction(func(tx *gorm.DB) error {
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

func (service RoleService) GetRoles() ([]structs.RoleView, error) {
	roleRepository := repositories.NewRoleRepository(service.GetDB())
	return roleRepository.GetAllRoles()
}
