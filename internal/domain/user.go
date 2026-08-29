package domain

import (
	"net/mail"
	"regexp"
	"strings"
	"time"
	"unicode"
	"unicode/utf8"
)

type User struct {
	ID           int64
	Email        string
	Username     string
	PasswordHash string
	CreatedAt    time.Time
}

func NewUser(username string, email string) (User, error) {
	user := User{
		Username: strings.ToLower(strings.TrimSpace(username)),
		Email:    strings.ToLower(strings.TrimSpace(email)),
	}

	if err := validateUsername(username); err != nil {
		return User{}, err
	}

	if err := user.validateEmail(); err != nil {
		return User{}, err
	}
	return user, nil
}

func (u *User) validateEmail() error {
	if u.Email == "" {
		return nil
	}

	addr, err := mail.ParseAddress(u.Email)
	if err != nil || addr.Address != u.Email {
		return ErrInvalidEmail
	}

	return nil
}

var usernameRegex = regexp.MustCompile(`^[a-z0-9_]+$`)
var onlyNumbersRegex = regexp.MustCompile(`^[0-9]+$`)

func validateUsername(username string) error {
	username = strings.ToLower(strings.TrimSpace(username))

	if len(username) < 3 || len(username) > 15 {
		return ErrInvalidUsernameLength
	}

	if !usernameRegex.MatchString(username) {
		return ErrInvalidUsernameChars
	}

	if onlyNumbersRegex.MatchString(username) {
		return ErrNumericUsername
	}

	return nil
}

func ValidatePassword(password string) error {
	if utf8.RuneCountInString(password) < 10 {
		return ErrPasswordTooShort
	}

	var hasNumber bool
	var hasUpper bool

	for _, char := range password {
		if unicode.IsSpace(char) {
			return ErrPasswordContainsSpace
		}

		if unicode.IsDigit(char) {
			hasNumber = true
		}

		if unicode.IsUpper(char) {
			hasUpper = true
		}
	}

	if !hasNumber {
		return ErrPasswordRequiresDigit
	}

	if !hasUpper {
		return ErrPasswordRequiresUpper
	}

	return nil
}
