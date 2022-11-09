package repository

import (
	"errors"
	"go-remix/appo"
	"gorm.io/gorm"
	"log"
	"regexp"

	"go-remix/model"
)

type userRepository struct {
	DB *gorm.DB
}

func NewUserRepository(db *gorm.DB) model.UserRepository {
	return &userRepository{
		DB: db,
	}
}

func (r *userRepository) Create(user *model.User) (*model.User, error) {
	if result := r.DB.Create(&user); result.Error != nil {

		if isDuplicateKeyError(result.Error) {
			return nil, appo.NewBadRequest(appo.MsgDuplicateEmail)
		}

		log.Printf("Could not create a user with email: %v. Reason: \n", user.Email, result.Error)

		return nil, appo.NewInternal()
	}

	return user, nil
}

func (r *userRepository) FindByEmail(email string) (*model.User, error) {
	user := &model.User{}

	if err := r.DB.Where("email = ?", email).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return user, appo.NewNotFound("email", email)
		}
		return user, appo.NewInternal()
	}

	return user, nil
}

func isDuplicateKeyError(err error) bool {
	duplicate := regexp.MustCompile(`\(SQLSTATE 23505\)$`)
	return duplicate.MatchString(err.Error())
}
