package middleware

import (
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
	"github.com/guptaharsh13/balkanid-task/initializers"
	"github.com/guptaharsh13/balkanid-task/models"
	"github.com/guptaharsh13/balkanid-task/utils"
)

func RequireAuth(c *gin.Context) {

	tokenString := c.Request.Header.Get("Authorization")
	prefix := "Bearer"
	if strings.HasPrefix(tokenString, prefix) {
		tokenString = strings.TrimSpace(tokenString[len(prefix):])
	}
	if len(tokenString) == 0 {
		c.AbortWithStatusJSON(http.StatusUnauthorized, utils.UnauthorizedResponse("Authorization token not found"))
		return
	}

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(os.Getenv("JWT_SECRET")), nil
	})
	if err != nil {
		c.AbortWithStatusJSON(http.StatusUnauthorized, utils.UnauthorizedResponse("Invalid token"))
		return
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		if exp, ok := claims["exp"]; ok {
			if exp.(float64) < float64(time.Now().Unix()) {
				c.AbortWithStatusJSON(http.StatusUnauthorized, utils.UnauthorizedResponse("Token expired"))
				return
			}
		}
		username, ok := claims["username"]
		if !ok {
			c.AbortWithStatusJSON(http.StatusInternalServerError, utils.InternalServerErrorResponse())
			return
		}
		var user models.User
		if result := initializers.DB.Take(&user, "username = ?", username); result.Error != nil {
			fmt.Printf("User with username %s not found: %s", username, result.Error.Error())
			c.AbortWithStatusJSON(http.StatusInternalServerError, utils.InternalServerErrorResponse())
			return
		}
		fmt.Println(user)
		c.Set("is_admin", user.IsAdmin)
		c.Set("username", user.Username)
		c.Next()
	} else {
		c.AbortWithStatusJSON(http.StatusUnauthorized, utils.UnauthorizedResponse("Invalid token"))
	}
}
