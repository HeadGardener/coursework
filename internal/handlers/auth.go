package handlers

import (
	"net/http"

	"github.com/HeadGardener/coursework/internal/dto"
	"github.com/gin-gonic/gin"
)

func (h *Handler) signUp(c *gin.Context) {
	var req dto.SignUpReq

	if err := c.BindJSON(&req); err != nil {
		newErrResponse(c, http.StatusBadRequest, "failed while decoding sign up request", err)
		return
	}

	if err := req.Validate(); err != nil {
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
	var req dto.SignInReq

	if err := c.BindJSON(&req); err != nil {
		newErrResponse(c, http.StatusBadRequest, "failed while decoding sign in request", err)
		return
	}

	if err := req.Validate(); err != nil {
		newErrResponse(c, http.StatusBadRequest, "failed while validating sign in request", err)
		return
	}

	tokens, err := h.authService.SignIn(c, req.Username, req.Password)
	if err != nil {
		newErrResponse(c, http.StatusInternalServerError, "failed while signing in", err)
		return
	}

	c.JSON(http.StatusCreated, tokens)
}

func (h *Handler) refresh(c *gin.Context) {
	var req dto.RefreshRequest
	if err := c.BindJSON(&req); err != nil {
		newErrResponse(c, http.StatusBadRequest, "failed while decoding RefreshRequest", err)
		return
	}

	tokens, err := h.authService.Refresh(c.Request.Context(), req.AccessToken, req.RefreshToken)
	if err != nil {
		newErrResponse(c, http.StatusBadRequest, "failed while refreshing", err)
		return
	}

	c.JSON(http.StatusOK, map[string]any{
		"access_token":  tokens.AccessToken,
		"refresh_token": tokens.RefreshToken,
	})
}

func (h *Handler) logout(c *gin.Context) {
	userID, err := getUserID(c)
	if err != nil {
		newErrResponse(c, http.StatusForbidden, "failed while getting user id", err)
		return
	}

	if err = h.authService.LogOut(c, userID); err != nil {
		newErrResponse(c, http.StatusInternalServerError, "failed while logging out", err)
		return
	}

	c.JSON(http.StatusOK, map[string]any{
		"status": "logged out",
	})
}
