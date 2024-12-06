package storage

// import (
// 	"log"
// 	"tender_management/models"

// 	"gorm.io/gorm"
// )


// type UserStorage struct {
// 	db *gorm.DB
// }

// func (s *UserStorage) Create(user *models.Users) error {
// 	if err := s.db.Create(user).Error; err != nil {
// 		log.Printf("Create: Error creating user: %v", err)
// 		return err
// 	}
// 	log.Printf("Create: User created with ID %d", user.ID)
// 	return nil
// }

// func (s *UserStorage) FindByEmail(email string) (*models.Users, error) {
// 	var user models.Users

// 	err := s.db.Where("email = ?", email).First(&user).Error
// 	if err != nil {
// 		log.Printf("FindByEmail error getting user by email %s: %v", email, err)
// 		return nil, err
// 	}

// 	log.Printf("FindByEmail User with email %s retrieved", email)
// 	return &user, nil
// }

// func (s *UserStorage) ActivateUser(phoneNumber string) error {
// 	result := s.db.Model(&models.Users{}).Where("phone_number = ?", phoneNumber).Update("is_active", true)
// 	return result.Error
// }
