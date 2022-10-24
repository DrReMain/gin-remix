package model

import "context"

type User struct {
	BaseModel
	Username string `gorm:"not null" json:"username"`
	Email    string `gorm:"not null;uniqueIndex" json:"email"`
	Password string `gorm:"not null" json:"-"`
}

type UserService interface {
	Get(id string) (*User, error)
	GetByEmail(email string) (*User, error)
	Register(user *User) (*User, error)
	Login(email, password string) (*User, error)
	UpdateAccount(user *User) error
	IsEmailAlreadyInUse(email string) bool
	ChangePassword(currentPassword, newPassword string, user *User) error
	ForgotPassword(ctx context.Context, user *User) error
	ResetPassword(ctx context.Context, password string, token string) (*User, error)
}

type UserRepository interface {
	FindById(id string) (*User, error)
	Create(user *User) (*User, error)
	FindByEmail(email string) (*User, error)
	Update(user *User) error
}
