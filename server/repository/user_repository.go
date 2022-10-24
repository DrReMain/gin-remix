package repository

import (
	"errors"
	"go-remix/model"
	"go-remix/model/apperrors"
	"gorm.io/gorm"
	"log"
	"regexp"
)

type userRepository struct {
	DB *gorm.DB
}

func NewUserRepository(db *gorm.DB) model.UserRepository {
	return &userRepository{
		DB: db,
	}
}

func (r *userRepository) FindById(id string) (*model.User, error) {
	user := &model.User{}

	if err := r.DB.Where("id = ?", id).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return user, apperrors.NewNotFound("uid", id)
		}
		return user, apperrors.NewInternal()
	}

	return user, nil
}

func (r *userRepository) Create(user *model.User) (*model.User, error) {
	if result := r.DB.Create(&user); result.Error != nil {
		if isDuplicateKeyError(result.Error) {
			return nil, apperrors.NewBadRequest(apperrors.DuplicateEmail)
		}

		log.Printf("Cound not create a user with email: %v. Reason: %v\n", user.Email, result.Error)
		return nil, apperrors.NewInternal()
	}

	return user, nil
}

func (r *userRepository) FindByEmail(email string) (*model.User, error) {
	user := &model.User{}

	if err := r.DB.Where("email = ?", email).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return user, apperrors.NewNotFound("email", email)
		}
		return user, apperrors.NewInternal()
	}

	return user, nil
}

func (r *userRepository) Update(user *model.User) error {
	return r.DB.Save(&user).Error
}

func isDuplicateKeyError(err error) bool {
	duplicate := regexp.MustCompile(`\(SQLSTATE 23505\)$`)
	return duplicate.MatchString(err.Error())
}
