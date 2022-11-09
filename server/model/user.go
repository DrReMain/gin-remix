package model

type User struct {
	BaseModel
	Username string `gorm:"not null" json:"username"`
	Email    string `gorm:"not null;uniqueIndex" json:"email"`
	Password string `gorm:"not null" json:"-"`
}

type UserService interface {
	Register(user *User) (*User, error)
	Login(email, password string) (*User, error)
}

type UserRepository interface {
	Create(user *User) (*User, error)
	FindByEmail(email string) (*User, error)
}
