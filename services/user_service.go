package services

import (
	"encoding/json"
	"fmt"
	"reflect"

	domainErrors "github.com/thoulee21/go-learn/errors"
	"github.com/thoulee21/go-learn/models"
	"gorm.io/gorm"
)

type UserService struct {
	DB *gorm.DB
}

func NewUserService(db *gorm.DB) (*UserService, error) {
	return &UserService{DB: db}, nil
}

func (r *UserService) GetAll() (*[]models.User, error) {
	var users []models.User
	if err := r.DB.Find(&users).Error; err != nil {
		return nil, domainErrors.NewAppErrorWithType(domainErrors.UnknownError)
	}
	return &users, nil
}

func (r *UserService) Create(userDomain *models.User) (*models.User, error) {
	userRepository := userDomain
	txDb := r.DB.Create(userRepository)
	err := txDb.Error
	if err != nil {
		byteErr, _ := json.Marshal(err)
		var newError domainErrors.GormErr
		errUnmarshal := json.Unmarshal(byteErr, &newError)
		if errUnmarshal != nil {
			return &models.User{}, errUnmarshal
		}
		switch newError.Number {
		case 1062:
			err = domainErrors.NewAppErrorWithType(domainErrors.ResourceAlreadyExists)
			return &models.User{}, err
		default:
			err = domainErrors.NewAppErrorWithType(domainErrors.UnknownError)
		}
	}
	return userRepository, err
}

func IsZeroValue(value any) bool {
	return reflect.DeepEqual(value, reflect.Zero(reflect.TypeOf(value)).Interface())
}

func (r *UserService) GetOneByMap(userMap map[string]interface{}) (*models.User, error) {
	var userRepository models.User
	tx := r.DB.Limit(1)

	for key, value := range userMap {
		if !IsZeroValue(value) {
			tx = tx.Where(fmt.Sprintf("%s = ?", key), value)
		}
	}
	if err := tx.Find(&userRepository).Error; err != nil {
		return &models.User{}, domainErrors.NewAppErrorWithType(domainErrors.UnknownError)
	}
	return &userRepository, nil
}

func (r *UserService) GetByID(id uint) (*models.User, error) {
	var user models.User

	dbRes := r.DB.Where("id = ?", id)
	err := dbRes.Error
	if err != nil {
		return &models.User{}, err
	}

	err = dbRes.First(&user).Error
	if err != nil {
		return &models.User{}, err
	}
	return &user, nil
}

func (r *UserService) Update(id uint, userMap *models.User) (*models.User, error) {
	var userObj models.User
	userObj.ID = id
	err := r.DB.Model(&userObj).Select("user_name", "email").Updates(userMap).Error
	if err != nil {
		byteErr, _ := json.Marshal(err)
		var newError domainErrors.GormErr
		errUnmarshal := json.Unmarshal(byteErr, &newError)
		if errUnmarshal != nil {
			return &models.User{}, errUnmarshal
		}
		switch newError.Number {
		case 1062:
			return &models.User{}, domainErrors.NewAppErrorWithType(domainErrors.ResourceAlreadyExists)
		default:
			return &models.User{}, domainErrors.NewAppErrorWithType(domainErrors.UnknownError)
		}
	}
	if err := r.DB.Where("id = ?", id).First(&userObj).Error; err != nil {
		return &models.User{}, err
	}
	return &userObj, nil
}

func (r *UserService) Delete(id uint) error {
	tx := r.DB.Delete(&models.User{}, id)
	if tx.Error != nil {
		return domainErrors.NewAppErrorWithType(domainErrors.UnknownError)
	}
	if tx.RowsAffected == 0 {
		return domainErrors.NewAppErrorWithType(domainErrors.NotFound)
	}
	return nil
}
