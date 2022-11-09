package service

import (
	"go-remix/appo"
	"go-remix/model"
	"log"
)

type UserService struct {
	UserRepository model.UserRepository
}

func NewUserService(p *UserService) model.UserService {
	return p
}

func (s *UserService) Register(user *model.User) (*model.User, error) {
	hashedPassword, err := hashPassword(user.Password)
	if err != nil {
		log.Printf("Unable to signup user for emai: %\n", user.Email)
		return nil, appo.NewInternal()
	}

	user.ID = GenerateId()
	user.Password = hashedPassword

	return s.UserRepository.Create(user)
}

func (s *UserService) Login(email, password string) (*model.User, error) {
	user, err := s.UserRepository.FindByEmail(email)
	if err != nil {
		return nil, appo.NewAuthorization(appo.MsgInvalidCredentials)
	}

	match, err := comparePasswords(user.Password, password)
	if err != nil {
		return nil, appo.NewInternal()
	}

	if !match {
		return nil, appo.NewAuthorization(appo.MsgInvalidCredentials)
	}

	return user, nil
}
