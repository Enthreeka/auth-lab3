package entity

import (
	"regexp"
	"unicode/utf8"
)

type User struct {
	ID       string `json:"id"`
	Password string `json:"password"`
	Login    string `json:"login"`

	Role Role `json:"role"`
}

func IsLoginValid(value string) bool {
	re := regexp.MustCompile("^[0-9A-Za-za-яА-Я]")

	if !re.MatchString(value) {
		return false
	}

	if utf8.RuneCountInString(value) < 6 {
		return false
	}

	return true
}

func IsPasswordValid(password string) bool {
	re := regexp.MustCompile(`^(?=.*[A-Z])(?=.*[a-z])(?=.*\d).{8,}$`)
	return re.MatchString(password)
}
