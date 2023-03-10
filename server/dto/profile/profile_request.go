package profiledto

import "time"

type CreateProfileRequest struct {
	Phone          string    `json:"phone" form:"phone" validate:"required"`
	Address        string    `json:"address" form:"address" validate:"required"`
	UserID         int       `json:"user_id" form:"user_id" validate:"required"`
	ProfilePicture string    `json:"profile_picture" form:"profile_picture" validate:"required"`
	CreatedAt      time.Time `json:"-"`
	UpdatedAt      time.Time `json:"-"`
}

type UpdateProfileRequest struct {
	Phone   string `json:"phone" form:"phone" `
	Address string `json:"address" form:"address"`
	UserID  int    `json:"user_id" form:"user_id"`
	Photo   string `json:"photo" form:"photo" validate:"required"`
}
