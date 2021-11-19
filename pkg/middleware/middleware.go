package middleware

import (
	"fmt"
	"kanggo/config"
	"kanggo/utils"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
)

func RoleAdmin() gin.HandlerFunc {
	return func(c *gin.Context) {
		bearer := c.Request.Header.Get("Authorization")

		token, err := jwt.Parse(bearer, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}
			return []byte(config.EnvFile.ApiSecret), nil
		})
		if err != nil {
			utils.Response(c, 401, err.Error(), nil)
			c.Abort()
			return
		}

		value := token.Claims.(jwt.MapClaims)

		if token.Claims == nil && err != nil {
			utils.Response(c, 401, "not authorized", nil)
			c.Abort()
			return
		}

		if value["admin"] == false {
			utils.Response(c, 401, "not authorized", nil)
			c.Abort()
			return
		}

		uid, _ := strconv.ParseUint(fmt.Sprintf("%.0f", value["user_id"]), 10, 32)

		c.Set("user_id", uid)
		c.Next()
	}
}

func RoleUser() gin.HandlerFunc {
	return func(c *gin.Context) {
		bearer := c.Request.Header.Get("Authorization")

		token, err := jwt.Parse(bearer, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}
			return []byte(config.EnvFile.ApiSecret), nil
		})
		if err != nil {
			utils.Response(c, 401, err.Error(), nil)
			c.Abort()
			return
		}

		if token.Claims == nil && err != nil {
			utils.Response(c, 401, "not authorized", nil)
			c.Abort()
			return
		}
		value := token.Claims.(jwt.MapClaims)

		uid, _ := strconv.ParseUint(fmt.Sprintf("%.0f", value["user_id"]), 10, 32)

		c.Set("user_id", uid)
		c.Next()
	}
}
