package auth

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	jwt "github.com/golang-jwt/jwt"
)

func AccessToken(signature string) gin.HandlerFunc {

	return func(c *gin.Context) {

		token := jwt.NewWithClaims(jwt.SigningMethodHS256, &jwt.StandardClaims{
			ExpiresAt: time.Now().Add(5 * time.Minute).Unix(),
			Audience:  "Boom", //ของจริงอาจจะเป็น username หรืออย่างอื่น แต่ในตัวอย่าง hardcode ไปก่อน

		})

		// log.Println(token.Claims)

		ss, err := token.SignedString([]byte(signature))
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": err.Error(),
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"token": ss,
		})
	}
}
