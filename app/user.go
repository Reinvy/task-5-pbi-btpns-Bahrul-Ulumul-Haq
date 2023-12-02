package app

import (
	"github.com/asaskevich/govalidator"
	"golang.org/x/crypto/bcrypt"
)

// User struct
type User struct {
	ID        uint   `gorm:"primaryKey" json:"id"`
	Username  string `gorm:"unique;not null" json:"username" valid:"required"`
	Email     string `gorm:"unique;not null" json:"email" valid:"required,email"`
	Password  string `gorm:"not null" json:"password" valid:"required,minstringlength(6)"`
	Photos    []Photo
	CreatedAt int64 `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt int64 `gorm:"autoUpdateTime" json:"updated_at"`
}

type UserInput struct {
	Email    string `json:"email" valid:"required,email"`
	Password string `json:"password" valid:"required,minstringlength(6)"`
}

// HashPassword hashes the user password
func (u *User) HashPassword() error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	u.Password = string(hashedPassword)
	return nil
}

// CheckPassword checks the user password
func (u *User) CheckPassword(password string, rawPassword string) error {
	return bcrypt.CompareHashAndPassword([]byte(password), []byte(rawPassword))
}

// Validate validates the user fields
func (u *User) Validate() error {
	_, err := govalidator.ValidateStruct(u)
	return err
}
