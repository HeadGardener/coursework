package handlers

import (
	"net/http"
	"strconv"

	"github.com/HeadGardener/coursework/internal/dto"
	"github.com/HeadGardener/coursework/internal/models"
	"github.com/gin-gonic/gin"
)

func (h *Handler) viewDrinks(c *gin.Context) {
	adult, err := getIsAdult(c)
	if err != nil {
		newErrResponse(c, http.StatusForbidden, "failed while identifying age", err)
		return
	}

	drinks, err := h.drinkService.GetAll(c, adult)
	if err != nil {
		newErrResponse(c, http.StatusInternalServerError, "failed while getting drinks", err)
		return
	}

	c.JSON(http.StatusOK, drinks)
}

func (h *Handler) viewByID(c *gin.Context) {
	drinkID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		newErrResponse(c, http.StatusBadRequest, "failed while checking id", err)
		return
	}

	adult, err := getIsAdult(c)
	if err != nil {
		newErrResponse(c, http.StatusForbidden, "failed while identifying age", err)
		return
	}

	drinks, err := h.drinkService.GetByID(c, drinkID, adult)
	if err != nil {
		newErrResponse(c, http.StatusInternalServerError, "failed while getting drinks", err)
		return
	}

	c.JSON(http.StatusOK, drinks)
}

func (h *Handler) addDrink(c *gin.Context) {
	var req dto.DrinkRequest
	if err := c.BindJSON(&req); err != nil {
		newErrResponse(c, http.StatusBadRequest, "failed while decoding drink request", err)
		return
	}

	if err := req.Validate(); err != nil {
		newErrResponse(c, http.StatusBadRequest, "failed while validating drink request", err)
		return
	}

	drink := &models.Drink{
		Name:   req.Name,
		Type:   req.Type,
		Bottle: req.Bottle,
		Cost:   req.Cost,
		Soft:   req.Soft,
	}

	id, err := h.drinkService.Add(c, drink)
	if err != nil {
		newErrResponse(c, http.StatusInternalServerError, "failed while adding drink", err)
		return
	}

	c.JSON(http.StatusCreated, map[string]any{
		"id": id,
	})
}

func (h *Handler) updateDrink(c *gin.Context) {
	drinkID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		newErrResponse(c, http.StatusBadRequest, "failed while checking id", err)
		return
	}

	var req dto.DrinkRequest
	if err = c.BindJSON(&req); err != nil {
		newErrResponse(c, http.StatusBadRequest, "failed while decoding update drink request", err)
		return
	}

	if err = req.Validate(); err != nil {
		newErrResponse(c, http.StatusBadRequest, "failed while validating update drink request", err)
		return
	}

	drink := &models.Drink{
		Name:   req.Name,
		Type:   req.Type,
		Bottle: req.Bottle,
		Cost:   req.Cost,
	}

	if err := h.drinkService.Update(c, drinkID, drink); err != nil {
		newErrResponse(c, http.StatusInternalServerError, "failed while updating drink", err)
		return
	}

	c.JSON(http.StatusOK, map[string]any{
		"status": "updated",
	})
}

func (h *Handler) deleteDrink(c *gin.Context) {
	drinkID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		newErrResponse(c, http.StatusBadRequest, "failed while checking id", err)
		return
	}

	if err := h.drinkService.Delete(c, drinkID); err != nil {
		newErrResponse(c, http.StatusInternalServerError, "failed while deleting drink", err)
		return
	}

	c.JSON(http.StatusOK, map[string]any{
		"status": "deleted",
	})
}
