package repository

import (
	"smartgas-payment/internal/models"
	"smartgas-payment/internal/schemas"
	"smartgas-payment/internal/utils"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

//go:generate mockery --name UserRepository --filename=mock_user.go --inpackage=true
type UserRepository interface {
	GetUserByEmail(string) (*models.User, error)
	CreateUser(*models.User) error
	GetUserByID(uuid.UUID) (*models.User, error)
	List(*schemas.Pagination, any) ([]*models.User, error)
	UpdateByID(uuid.UUID, *models.User) error
	GetUserDetailByID(uuid.UUID) (*models.User, error)
}

type userRepository struct {
	DB *gorm.DB
}

func ProvideUserRepository(db *gorm.DB) *userRepository {
	return &userRepository{
		DB: db,
	}
}

func (ur *userRepository) GetUserByEmail(email string) (*models.User, error) {

	var user models.User

	result := ur.DB.Where(models.User{Email: email, Active: utils.BoolAddr(true)}).First(&user)

	if result.Error != nil {
		return nil, result.Error
	}

	return &user, nil
}

func (ur *userRepository) CreateUser(user *models.User) error {
	result := ur.DB.
		Omit("Permissions.*", "Groups.*", "GasStations.*").
		Create(&user)

	return result.Error
}

func (ur *userRepository) GetUserByID(id uuid.UUID) (*models.User, error) {
	var user models.User

	result := ur.DB.
		Preload("Groups").
		Preload("Permissions").
		Preload("Groups.Permissions").
		Preload("GasStations").
		Where(models.User{Active: utils.BoolAddr(true)}).First(&user, id)

	if result.Error != nil {
		return nil, result.Error
	}

	return &user, nil
}

func (ur *userRepository) GetUserDetailByID(id uuid.UUID) (*models.User, error) {
	var user models.User

	result := ur.DB.
		Preload("Groups").
		Preload("Permissions").
		Preload("Groups.Permissions").
		Preload("GasStations").
		First(&user, id)

	if result.Error != nil {
		return nil, result.Error
	}

	return &user, nil
}

func (ur *userRepository) List(pagination *schemas.Pagination, filters any) ([]*models.User, error) {
	var users []*models.User

	sql := `
    CONCAT(users.first_name, ' ', users.last_name) LIKE @search OR
    users.email LIKE @search
  `

	result := ur.DB.
		Preload("Permissions").
		Preload("Groups").
		Preload("Groups.Permissions").
		Preload("GasStations").
		Scopes(utils.Paginate(pagination, users, ur.DB, sql, filters, "")).
		Where(sql, filters).
		Order("created_at desc").
		Find(&users)

	if result.Error != nil {
		return nil, result.Error
	}

	return users, nil
}

func (ur *userRepository) UpdateByID(id uuid.UUID, user *models.User) error {
	user.ID = id
	result := ur.DB.
		Omit("Groups", "Permissions", "GasStations").
		Updates(user)

	if err := ur.DB.Model(user).Omit("Permissions.*").Association("Permissions").Replace(user.Permissions); err != nil {
		return err
	}

	if err := ur.DB.Model(user).Omit("Groups.*").Association("Groups").Replace(user.Groups); err != nil {
		return err
	}

	if err := ur.DB.Model(user).Omit("GasStations.*").Association("GasStations").Replace(user.GasStations); err != nil {
		return err
	}

	return result.Error
}
