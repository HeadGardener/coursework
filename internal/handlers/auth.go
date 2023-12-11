package handlers

import (
	"errors"
	"net/http"
	"regexp"

	"github.com/gin-gonic/gin"
)

var (
	checkUsername = regexp.MustCompile(`[0-9A-z]$`)
	checkName     = regexp.MustCompile(`[A-z]$`)
	checkPassword = regexp.MustCompile(`[0-9A-z]{8,16}$`)
)

type signUpReq struct {
	Username string `json:"username"`
	Name     string `json:"name"`
	Age      int    `json:"age"`
	Password string `json:"password"`
}

type signInReq struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func (h *Handler) signUp(c *gin.Context) {
	var req signUpReq

	if err := c.BindJSON(&req); err != nil {
		newErrResponse(c, http.StatusBadRequest, "failed while decoding sign up request", err)
		return
	}

	if err := req.validate(); err != nil {
		newErrResponse(c, http.StatusBadRequest, "failed while validating sign up request", err)
		return
	}

	id, err := h.authService.SignUp(c, req.Username, req.Name, req.Age, req.Password)
	if err != nil {
		newErrResponse(c, http.StatusInternalServerError, "failed while signing up", err)
		return
	}

	c.JSON(http.StatusCreated, map[string]any{
		"id": id,
	})
}

func (h *Handler) signIn(c *gin.Context) {
	var req signInReq

	if err := c.BindJSON(&req); err != nil {
		newErrResponse(c, http.StatusBadRequest, "failed while decoding sign in request", err)
		return
	}

	if err := req.validate(); err != nil {
		newErrResponse(c, http.StatusBadRequest, "failed while validating sign in request", err)
		return
	}

	token, err := h.authService.SignIn(c, req.Username, req.Password)
	if err != nil {
		newErrResponse(c, http.StatusInternalServerError, "failed while signing in", err)
		return
	}

	c.JSON(http.StatusCreated, map[string]any{
		"token": token,
	})
}

func (h *Handler) logout(c *gin.Context) {
	userID, err := getUserID(c)
	if err != nil {
		newErrResponse(c, http.StatusForbidden, "failed while getting user id", err)
	}

	if err = h.authService.LogOut(c, userID); err != nil {
		newErrResponse(c, http.StatusInternalServerError, "failed while logging out", err)
	}

	c.JSON(http.StatusOK, map[string]any{
		"status": "logged out",
	})
}

func (r *signUpReq) validate() error {
	if !checkUsername.MatchString(r.Username) {
		return errors.New("invalid user name: must contain only letters, numbers and symbols(_-) ")
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

func (r *signInReq) validate() error {
	if !checkUsername.MatchString(r.Username) {
		return errors.New("invalid username: must contain only letters, numbers and symbols(_-) ")
	}

	if !checkPassword.MatchString(r.Password) {
		return errors.New("invalid password: must contain only letters and numbers")
	}

	return nil
}
