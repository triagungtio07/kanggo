package utils

import (
	"github.com/gin-gonic/gin"
)

type (
	Respond struct {
		Status  int         `json:"status"`
		Message string      `json:"message"`
		Data    interface{} `json:"data"`
	}
)

func Response(c *gin.Context, status int, msg string, data interface{}) {
	c.JSON(status, &Respond{
		Status:  status,
		Message: msg,
		Data:    data,
	})
}
