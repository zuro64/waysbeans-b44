package repositories

import (
	"nis-waybeans/models"

	"gorm.io/gorm"
)

type ProfileRepository interface {
	FindProfiles() ([]models.Profile, error)
	GetProfile(ID int) (models.Profile, error)
	GetProfileByUserID(userID int) (models.Profile, error)
	CreateProfile(Profile models.Profile) (models.Profile, error)
	UpdateProfile(Profile models.Profile) (models.Profile, error)
	DeleteProfile(Profile models.Profile, ID int) (models.Profile, error)
}

func RepositoryProfile(db *gorm.DB) *repository {
	return &repository{db}
}

func (r *repository) FindProfiles() ([]models.Profile, error) {
	var profiles []models.Profile
	err := r.db.Preload("User").Find(&profiles).Error
	return profiles, err
}

func (r *repository) GetProfile(ID int) (models.Profile, error) {
	var profile models.Profile
	err := r.db.First(&profile, ID).Preload("User").Error
	return profile, err
}

func (r *repository) GetProfileByUserID(userID int) (models.Profile, error) {
	var profile models.Profile
	err := r.db.Where("user_id = ?", userID).First(&profile).Error
	return profile, err
}

func (r *repository) CreateProfile(profile models.Profile) (models.Profile, error) {
	err := r.db.Create(&profile).Error
	return profile, err
}

func (r *repository) UpdateProfile(profile models.Profile) (models.Profile, error) {
	err := r.db.Save(&profile).Error
	return profile, err
}

func (r *repository) DeleteProfile(profile models.Profile, ID int) (models.Profile, error) {
	err := r.db.Delete(&profile, ID).Error
	return profile, err
}
