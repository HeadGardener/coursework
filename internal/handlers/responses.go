package handlers

import (
	"log"

	"github.com/gin-gonic/gin"
)

type response struct {
	Msg   string `json:"Msg"`
	Error string `json:"Error"`
}

func newErrResponse(c *gin.Context, code int, msg string, err error) {
	log.Printf("[ERROR] %s: %s", msg, err.Error())

	c.AbortWithStatusJSON(code, response{
		Msg:   msg,
		Error: err.Error(),
	})
}
