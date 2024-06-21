package campaign

import (
	"crowdfunding-minpro-alterra/modules/user"
	"time"
)

type Campaign struct {
	ID               int          `gorm:"column:id;primaryKey"`
	UserID           int          `gorm:"column:user_id"`
	Name             string       `gorm:"column:name"`
	ShortDescription string       `gorm:"column:short_description"`
	Description      string       `gorm:"column:description;type:TEXT"`
	// Perks            string       `gorm:"column:perks;type:TEXT"`
	BackerCount      int          `gorm:"column:backer_count"`
	GoalAmount       int          `gorm:"column:goal_amount"`
	CurrentAmount    int          `gorm:"column:current_amount"`
	Slug             string       `gorm:"column:slug"`
	CreatedAt        time.Time    `gorm:"column:created_at"`
	UpdatedAt        time.Time    `gorm:"column:updated_at"`
	CampaignImages   []CampaignImage `gorm:"foreignKey:CampaignID"`
	User             user.User    `gorm:"foreignKey:UserID"`
}

type CampaignImage struct {
	ID         int       `gorm:"column:id;primaryKey"`
	CampaignID int       `gorm:"column:campaign_id"`
	FileName   string    `gorm:"column:file_name"`
	IsPrimary  int       `gorm:"column:is_primary"`
	CreatedAt  time.Time `gorm:"column:created_at"`
	UpdatedAt  time.Time `gorm:"column:updated_at"`
	DeletedAt  *time.Time `gorm:"column:deleted_at"`
}
