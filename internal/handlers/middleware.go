package handlers

import (
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/HeadGardener/coursework/internal/lib/auth"
	"github.com/HeadGardener/coursework/internal/models"
	"github.com/gin-gonic/gin"
)

const (
	userCtx string = "userAtr"
	isAdult string = "is_adult"
)

const (
	headerPartsLen = 2
)

var (
	ErrUserCtxNotExist   = errors.New("userCtx not exists")
	ErrNotUserAttributes = errors.New("userCtx value is not of type UserAttributes")
	ErrInvalidRole       = errors.New("user not admin")
	ErrNotBool           = errors.New("value is not of bool type")
)

func (h *Handler) identifyUser(c *gin.Context) {
	header := c.GetHeader("Authorization")

	if header == "" {
		newErrResponse(c, http.StatusUnauthorized, "failed while identifying user",
			errors.New("empty auth header"))
	}

	headerParts := strings.Split(header, " ")
	if len(headerParts) != headerPartsLen {
		newErrResponse(c, http.StatusUnauthorized, "failed while identifying user",
			errors.New("invalid auth header, must be like `Bearer token`"))
	}

	if headerParts[0] != "Bearer" {
		newErrResponse(c, http.StatusUnauthorized, "failed while identifying user",
			fmt.Errorf("invalid auth header %s, must be Bearer", headerParts[0]))
	}

	token := headerParts[1]
	if token == "" {
		newErrResponse(c, http.StatusUnauthorized, "failed while identifying user",
			errors.New("jwt token is empty"))
	}

	userAttributes, err := auth.ParseToken(token)
	if err != nil {
		newErrResponse(c, http.StatusUnauthorized, "failed while parsing token", err)
	}

	if err = h.authService.Check(c.Request.Context(), userAttributes.ID, token); err != nil {
		newErrResponse(c, http.StatusUnauthorized, "failed while checking session", err)
	}

	c.Set(userCtx, userAttributes)
}

func (h *Handler) identifyRole(c *gin.Context) {
	userAttributes, err := getUserAttributes(c)
	if err != nil {
		newErrResponse(c, http.StatusForbidden, "invalid user ctx", err)
	}

	if userAttributes.Role != models.RoleAdmin {
		newErrResponse(c, http.StatusForbidden, "invalid user role", ErrInvalidRole)
	}
}

func (h *Handler) checkAge(c *gin.Context) {
	userAttributes, err := getUserAttributes(c)
	if err != nil {
		newErrResponse(c, http.StatusForbidden, "invalid user ctx", err)
	}

	if userAttributes.Age < models.AdultAge {
		c.Set(isAdult, false)
		c.Next()
	}

	c.Set(isAdult, true)
}

func getUserID(c *gin.Context) (string, error) {
	userAttributes, err := getUserAttributes(c)
	if err != nil {
		return "", err
	}

	return userAttributes.ID, nil
}

func getIsAdult(c *gin.Context) (bool, error) {
	v, ok := c.Get(isAdult)
	if !ok {
		return false, ErrUserCtxNotExist
	}

	adult, ok := v.(bool)
	if !ok {
		return false, ErrNotBool
	}

	return adult, nil
}

func getUserAttributes(c *gin.Context) (auth.UserAttributes, error) {
	v, ok := c.Get(userCtx)
	if !ok {
		return auth.UserAttributes{}, ErrUserCtxNotExist
	}

	userAttributes, ok := v.(auth.UserAttributes)
	if !ok {
		return auth.UserAttributes{}, ErrNotUserAttributes
	}

	return userAttributes, nil
}
