package controllers

import (
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
	"github.com/guptaharsh13/balkanid-task/initializers"
	"github.com/guptaharsh13/balkanid-task/models"
	"github.com/guptaharsh13/balkanid-task/utils"
	"golang.org/x/crypto/bcrypt"
)

func Signup(c *gin.Context) {

	var body struct {
		Username string `json:"username" validate:"required"`
		Email    string `json:"email" validate:"email,required"`
		Password string `json:"password" validate:"required"`
	}
	if c.Bind(&body) != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Bad Request",
		})
		return
	}
	err := utils.ValidateStruct(body)
	if err != nil {
		c.JSON(http.StatusBadRequest, utils.BadRequestResponse(err.Error()))
		return
	}

	var user models.User
	result := initializers.DB.First(&user, "username = ?", body.Username)
	if result.Error == nil {
		c.JSON(http.StatusConflict, utils.ConflictResponse("Username already taken"))
		return
	}
	result = initializers.DB.First(&user, "email = ?", body.Email)
	if result.Error == nil {
		c.JSON(http.StatusConflict, utils.ConflictResponse("Email already taken"))
		return
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(body.Password), 10)
	if err != nil {
		c.JSON(http.StatusInternalServerError, utils.InternalServerErrorResponse("Internal Server Error"))
		fmt.Printf("Couldn't hash password: %s", err.Error())
		return
	}

	user = models.User{
		Username: body.Username,
		Email:    body.Email,
		Password: string(hash),
	}
	result = initializers.DB.Create(&user)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, utils.InternalServerErrorResponse("Internal Server Error"))
		fmt.Printf("Couldn't create User: %s", err.Error())
		return
	}
	c.JSON(http.StatusOK, utils.SuccessResponse(user))
}

func Login(c *gin.Context) {

	var body struct {
		Username string `json:"username"`
		Email    string `json:"email" validate:"email"`
		Password string `json:"password" validate:"required"`
	}

	if c.Bind(&body) != nil {
		c.JSON(http.StatusBadRequest, utils.BadRequestResponse("Bad Request"))
		return
	}
	err := utils.ValidateStruct(body)
	if err != nil {
		c.JSON(http.StatusBadRequest, utils.BadRequestResponse(err.Error()))
		return
	}

	var user models.User
	if len(body.Username) > 0 {
		initializers.DB.First(&user, "username = ?", body.Username)
	} else if len(body.Email) > 0 {
		initializers.DB.First(&user, "email = ?", body.Email)
	} else {
		c.JSON(http.StatusBadRequest, utils.BadRequestResponse("Username or Email required"))
		return
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(body.Password))
	if err != nil {
		c.JSON(http.StatusUnauthorized, utils.UnauthorizedResponse("Invalid Credentials"))
		return
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"username": user.Username,
		"exp":      time.Now().Add(time.Hour * 24 * 15).Unix(),
	})

	tokenString, err := token.SignedString([]byte(os.Getenv("JWT_SECRET")))
	if err != nil {
		c.JSON(http.StatusInternalServerError, utils.InternalServerErrorResponse("Internal Server Error"))
		fmt.Printf("Couldn't create token: %s", err.Error())
		return
	}
	data := struct {
		Token string `json:"token"`
	}{
		Token: tokenString,
	}
	c.JSON(http.StatusOK, utils.SuccessResponse(data))
}
