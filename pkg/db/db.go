package db

import (
	"fmt"
	"tender_management/config"
	"tender_management/models"
	"log"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func ConnectDB(cfg config.Config) (*gorm.DB, error) {
	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		cfg.DBHost, cfg.DBPort, cfg.DBUser, cfg.DBPassword, cfg.DBName,
	)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	if err := db.AutoMigrate(&models.Users{}, &models.Tenders{}, &models.Offers{}); err != nil {
		log.Fatal("Error Migratilon")
	}


	log.Println("Connected database... ")
	return db, nil
}
