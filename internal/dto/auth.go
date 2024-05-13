package dto

import (
	"errors"
	"regexp"
)

var (
	checkUsername = regexp.MustCompile(`^[0-9A-Za-z]+$`)
	checkName     = regexp.MustCompile(`^[A-Za-z]+$`)
	checkPassword = regexp.MustCompile(`[0-9A-z]{8,16}$`)
)

type SignUpReq struct {
	Username string `json:"username"`
	Name     string `json:"name"`
	Age      int    `json:"age"`
	Password string `json:"password"`
}

type SignInReq struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type RefreshRequest struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

func (r *SignUpReq) Validate() error {
	if !checkUsername.MatchString(r.Username) {
		return errors.New("invalid username: must contain only letters, numbers and symbols(_-) ")
	}

	if !checkName.MatchString(r.Name) {
		return errors.New("invalid name: must contain only letters")
	}

	if r.Age <= 0 || r.Age >= 112 {
		return errors.New("invalid age: can't be less than 0 or greater than 111")
	}

	if !checkPassword.MatchString(r.Password) {
		return errors.New("invalid password: must contain only letters and numbers")
	}

	return nil
}

func (r *SignInReq) Validate() error {
	if !checkUsername.MatchString(r.Username) {
		return errors.New("invalid username: must contain only letters, numbers and symbols(_-) ")
	}

	if !checkPassword.MatchString(r.Password) {
		return errors.New("invalid password: must contain only letters and numbers")
	}

	return nil
}
