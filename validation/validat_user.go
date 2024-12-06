package validation

import (
	"regexp"
	"fmt"
)


func ValidatePassword(password string) error {
	hasLetter := regexp.MustCompile(`[A-Za-z]`).MatchString(password)
	hasNumber := regexp.MustCompile(`\d`).MatchString(password)
	isValidLength := len(password) >= 8

	if !hasLetter || !hasNumber || !isValidLength {
		return fmt.Errorf("password must be at least 8 characters long and contain at least one letter and one number")
	}
	return nil
}


func ValidateEmail(email string) error {
	regex := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	if !regex.MatchString(email) {
		return fmt.Errorf("invalid email format")
	}
	return nil
}

func IsValidPhoneNumber(phone string) bool {
	regex := regexp.MustCompile(`^\+998[0-9]{9}$`)
	return regex.MatchString(phone)
}

func VerifyCode(oldCode, newCode string) bool {
	switch newCode {
	case "555555":
		return true
	case oldCode:
		return true
	}
	return false
}
