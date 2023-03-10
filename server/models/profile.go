package models

import "time"

type Profile struct {
	ID             int          `json:"id" gorm:"primary_key:auto_increment"`
	Phone          string       `json:"phone" gorm:"type: varchar(255)"`
	Address        string       `json:"address" gorm:"type: varchar(255)"`
	ProfilePicture string       `json:"profile_picture" gorm:"type: varchar(255)"`
	UserID         int          `json:"user_id" gorm:"type: int" form:"user_id"`
	User           UserResponse `json:"user"`
	CreatedAt      time.Time    `json:"-"`
	UpdatedAt      time.Time    `json:"-"`
}

type ProfileResponse struct {
	Phone          string `json:"phone" gorm:"type: varchar(255)"`
	Address        string `json:"address" gorm:"type: varchar(255)"`
	ProfilePicture string `json:"profile_picture" gorm:"type: varchar(255)"`
	UserID         int    `json:"-"`
}

func (ProfileResponse) TableName() string {
	return "profiles"
}
