package database

import (
	"crowdfunding-minpro-alterra/modules/campaign"
	"crowdfunding-minpro-alterra/modules/donation"
	"crowdfunding-minpro-alterra/modules/user"
	"fmt"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type Config struct {
	DBUser string
	DBPass string
	DBHost string
	DBPort string
	DBName string
}

func ConnectDB() (*gorm.DB, error) {
	config := Config{
		DBName: "pbb_kelompok_7",
		DBUser: "root",
		DBPass: "",
		DBHost: "localhost",
		DBPort: "3306",
	}

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local", config.DBUser, config.DBPass, config.DBHost, config.DBPort, config.DBName)
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	MigrateAllEntities(db)
	return db, nil
}

func MigrateAllEntities(db *gorm.DB) {
	db.AutoMigrate(&user.User{}, &campaign.Campaign{}, &campaign.CampaignImage{}, &donation.Donation{})
}
