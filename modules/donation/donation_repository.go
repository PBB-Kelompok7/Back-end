package donation

import "gorm.io/gorm"

type repository struct {
	db *gorm.DB
}

type Repository interface {
	GetByCampaignID(CampaignID int) ([]Donation, error)
	GetByUserID(UserID int) ([]Donation, error)
	GetByID(ID int) (Donation, error)
	Save(donation Donation) (Donation, error)
	Update(donation Donation) (Donation, error)
}

func NewRepository(db *gorm.DB) *repository {
	return &repository{db}
}

func (r *repository) GetByCampaignID(CampaignID int) ([]Donation, error) {
	var donations []Donation

	err := r.db.Preload("User").Where("campaign_id = ?", CampaignID).Order("created_at desc").Find(&donations).Error

	if err != nil {
		return donations, err
	}

	return donations, nil
}

func (r *repository) GetByUserID(UserID int) ([]Donation, error) {
	var donations []Donation

	err := r.db.Preload("Campaign.CampaignImages", "campaign_images.is_primary = 1").Where("user_id = ?", UserID).Order("created_at desc").Find(&donations).Error

	if err != nil {
		return donations, err
	}

	return donations, nil
}

func (r *repository) GetByID(ID int) (Donation, error) {
	var donation Donation

	err := r.db.Where("id = ?", ID).Find(&donation).Error

	if err != nil {
		return donation, err
	}

	return donation, nil
}

func (r *repository) Save(donation Donation) (Donation, error) {
	err := r.db.Create(&donation).Error

	if err != nil {
		return donation, err
	}

	return donation, nil
}

func (r *repository) Update(donation Donation) (Donation, error) {
	err := r.db.Save(&donation).Error

	if err != nil {
		return donation, err
	}

	return donation, nil
}