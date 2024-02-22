package validator

import (
	"fmt"
	"net/mail"
	"regexp"
)

var (
	isValidFormatUsername = regexp.MustCompile(`^[a-z0-9_]+$`).MatchString
	isValidFormatFullName = regexp.MustCompile(`^[a-zA-Z\s]+$`).MatchString
)

func ValidateString(s string, minLen int, maxLen int) error {
	n := len(s)
	if n < minLen || n > maxLen {
		return fmt.Errorf("must contain from %d-%d characters", minLen, maxLen)
	}
	return nil
}

func ValidateUsername(u string) error {
	if err := ValidateString(u, 3, 100); err != nil {
		return err
	}

	if !isValidFormatUsername(u) {
		return fmt.Errorf("must contain only lowercase letters, digits, or underscores")
	}

	return nil
}

func ValidateFullName(fn string) error {
	if err := ValidateString(fn, 3, 100); err != nil {
		return err
	}

	if !isValidFormatFullName(fn) {
		return fmt.Errorf("must contain only letters or spaces")
	}

	return nil
}

func ValidatePassword(p string) error {
	return ValidateString(p, 6, 100)
}

func ValidateEmail(e string) error {
	if err := ValidateString(e, 3, 200); err != nil {
		return err
	}

	if _, err := mail.ParseAddress(e); err != nil {
		return fmt.Errorf("email is not valid")
	}

	return nil
}
