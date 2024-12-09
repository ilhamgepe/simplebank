package val

import (
	"fmt"
	"net/mail"
	"regexp"
)

var (
	isValidUsername = regexp.MustCompile(`^[a-z0-9_]+$`).MatchString
	isValidFullname = regexp.MustCompile(`^[a-zA-Z\s]+$`).MatchString
)

func ValidateString(value string, min,max int) error{
	n := len(value)
	if n < min || n > max {
		return fmt.Errorf("length must be between %d and %d", min, max)
	}
	return nil
}

func ValidateUsername(value string) error{
	if err := ValidateString(value, 3, 20); err != nil {
		return err
	}
	if !isValidUsername(value){
		return fmt.Errorf("username must contain only lowercase letters, numbers, and underscores")
	}
	return nil
}
func ValidateFullname(value string) error{
	if err := ValidateString(value, 3, 20); err != nil {
		return err
	}
	if !isValidFullname(value){
		return fmt.Errorf("fullname must contain only letters or spaces")
	}
	return nil
}

func ValidatePassword(value string) error{
	return ValidateString(value, 6, 100)
}


func ValidateEmail(value string) error{
	if err := ValidateString(value, 3, 100); err != nil {
		return err
	}
	if _,err  := mail.ParseAddress(value); err != nil {
		return fmt.Errorf("invalid email address")
	}
	
	return nil
}