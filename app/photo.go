package app

import (
	"github.com/asaskevich/govalidator"
)

// Photo struct
type Photo struct {
	ID       uint   `gorm:"primaryKey" json:"id"`
	Title    string `json:"title"`
	Caption  string `json:"caption"`
	PhotoUrl string `json:"photo_url" valid:"required,url"`
	UserID   uint   `json:"user_id"`
	User     User
}
type PhotoInput struct {
	Title    string `json:"title" valid:"required"`
	Caption  string `json:"caption" valid:"required"`
	PhotoUrl string `json:"photo_url" valid:"required,url"`
	UserID   uint   `json:"user_id" valid:"required"`
}

// Validate validates the photo fields
func (p *Photo) Validate() error {
	_, err := govalidator.ValidateStruct(p)
	return err
}
