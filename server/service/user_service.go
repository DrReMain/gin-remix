package service

import (
	"context"
	"go-remix/model"
	"go-remix/model/apperrors"
	"log"
	"strings"
)

type userService struct {
	UserRepository  model.UserRepository
	RedisRepository model.RedisRepository
	MailRepository  model.MailRepository
}

type USConfig struct {
	UserRepository  model.UserRepository
	RedisRepository model.RedisRepository
	MailRepository  model.MailRepository
}

func NewUserService(c *USConfig) model.UserService {
	return &userService{
		UserRepository:  c.UserRepository,
		RedisRepository: c.RedisRepository,
		MailRepository:  c.MailRepository,
	}
}

func (s *userService) Get(uid string) (*model.User, error) {
	return s.UserRepository.FindById(uid)
}

func (s *userService) GetByEmail(email string) (*model.User, error) {
	email = strings.ToLower(email)
	email = strings.TrimSpace(email)

	return s.UserRepository.FindByEmail(email)
}

func (s *userService) Register(user *model.User) (*model.User, error) {
	hashedPassword, err := hashPassword(user.Password)
	if err != nil {
		log.Printf("Unable to signup user for email: %v\n", user.Email)
		return nil, apperrors.NewInternal()
	}

	user.ID = GenerateId()
	user.Password = hashedPassword

	return s.UserRepository.Create(user)
}

func (s *userService) Login(email, password string) (*model.User, error) {
	user, err := s.UserRepository.FindByEmail(email)
	if err != nil {
		return nil, apperrors.NewAuthorization(apperrors.InvalidCredentials)
	}

	match, err := comparePasswords(user.Password, password)
	if err != nil {
		return nil, apperrors.NewInternal()
	}

	if !match {
		return nil, apperrors.NewAuthorization(apperrors.InvalidCredentials)
	}

	return user, nil
}

func (s *userService) UpdateAccount(u *model.User) error {
	return s.UserRepository.Update(u)
}

func (s *userService) IsEmailAlreadyInUse(email string) bool {
	user, err := s.UserRepository.FindByEmail(email)
	if err != nil {
		return true
	}

	return user.ID != ""
}

func (s *userService) ChangePassword(currentPassword, newPassword string, user *model.User) error {
	match, err := comparePasswords(user.Password, currentPassword)
	if err != nil {
		return apperrors.NewInternal()
	}

	if !match {
		return apperrors.NewAuthorization(apperrors.InvalidCredentials)
	}

	hashedPassword, err := hashPassword(newPassword)
	if err != nil {
		log.Printf("Unable to change password for email: %v\n", user.Email)
		return apperrors.NewInternal()
	}

	user.Password = hashedPassword

	return s.UserRepository.Update(user)
}

func (s *userService) ForgotPassword(ctx context.Context, user *model.User) error {
	token, err := s.RedisRepository.SetResetToken(ctx, user.ID)

	if err != nil {
		return err
	}

	return s.MailRepository.SendResetMail(user.Email, token)
}

func (s *userService) ResetPassword(ctx context.Context, password string, token string) (*model.User, error) {
	id, err := s.RedisRepository.GetIdFromToken(ctx, token)
	if err != nil {
		return nil, err
	}

	user, err := s.UserRepository.FindById(id)
	if err != nil {
		return nil, err
	}

	hashedPassword, err := hashPassword(password)
	if err != nil {
		log.Printf("Unable to reset password")
		return nil, apperrors.NewInternal()
	}

	user.Password = hashedPassword

	if err = s.UserRepository.Update(user); err != nil {
		return nil, err
	}
	return user, nil
}
