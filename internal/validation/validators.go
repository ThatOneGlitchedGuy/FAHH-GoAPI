package validation

import (
	"regexp"
	"unicode"
)

type ValidationRule struct {
	Field string
	Error string
}

func ValidateEmail(email string) bool {
	const emailPattern = `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`
	match, _ := regexp.MatchString(emailPattern, email)
	return match && len(email) <= 254
}

func ValidatePassword(password string) []ValidationRule {
	var rules []ValidationRule

	if len(password) < 8 {
		rules = append(rules, ValidationRule{
			Field: "password",
			Error: "Password must be at least 8 characters long",
		})
	}

	if len(password) > 128 {
		rules = append(rules, ValidationRule{
			Field: "password",
			Error: "Password must not exceed 128 characters",
		})
	}

	hasUpper := false
	hasLower := false
	hasDigit := false
	hasSpecial := false

	for _, char := range password {
		switch {
		case unicode.IsUpper(char):
			hasUpper = true
		case unicode.IsLower(char):
			hasLower = true
		case unicode.IsDigit(char):
			hasDigit = true
		case unicode.IsPunct(char) || unicode.IsSymbol(char):
			hasSpecial = true
		}
	}

	if !hasUpper {
		rules = append(rules, ValidationRule{
			Field: "password",
			Error: "Password must contain at least one uppercase letter",
		})
	}

	if !hasLower {
		rules = append(rules, ValidationRule{
			Field: "password",
			Error: "Password must contain at least one lowercase letter",
		})
	}

	if !hasDigit {
		rules = append(rules, ValidationRule{
			Field: "password",
			Error: "Password must contain at least one digit",
		})
	}

	if !hasSpecial {
		rules = append(rules, ValidationRule{
			Field: "password",
			Error: "Password must contain at least one special character",
		})
	}

	return rules
}

func ValidateName(name string) bool {
	if len(name) < 2 || len(name) > 100 {
		return false
	}
	for _, char := range name {
		if !unicode.IsLetter(char) && !unicode.IsSpace(char) && char != '-' && char != '\'' {
			return false
		}
	}
	return true
}

func ValidatePhoneNumber(phone string) bool {
	const phonePattern = `^[0-9\-\+\(\)\s]{10,20}$`
	match, _ := regexp.MatchString(phonePattern, phone)
	return match
}

func ValidateURL(url string) bool {
	const urlPattern = `^https?://[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}(/[^\s]*)?$`
	match, _ := regexp.MatchString(urlPattern, url)
	return match
}

func ValidatePrice(price float64) bool {
	return price >= 0 && price <= 999999.99
}

func ValidateQuantity(quantity int) bool {
	return quantity > 0 && quantity <= 10000
}

func ValidateRating(rating int) bool {
	return rating >= 1 && rating <= 5
}
