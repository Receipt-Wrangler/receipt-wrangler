package repositories

import (
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

	roles := make([]structs.RoleView, 0, len(appRoles)+len(groupRoles))

	for _, role := range appRoles {
		perms := make([]string, 0, len(role.Permissions))
		for _, permission := range role.Permissions {
			perms = append(perms, permission.Permission)
		}

		roles = append(roles, structs.RoleView{
			Id:          role.ID,
			Name:        role.Name,
			Description: role.Description,
			Scope:       permissions.ScopeApp,
			IsSystem:    role.IsSystem,
			Permissions: perms,
		})
	}

	for _, role := range groupRoles {
		perms := make([]string, 0, len(role.Permissions))
		for _, permission := range role.Permissions {
			perms = append(perms, permission.Permission)
		}

		roles = append(roles, structs.RoleView{
			Id:          role.ID,
			Name:        role.Name,
			Description: role.Description,
			Scope:       permissions.ScopeGroup,
			IsSystem:    role.IsSystem,
			Permissions: perms,
		})
	}

	return roles, nil
}
