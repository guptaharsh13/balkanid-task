package controllers

import (
	"encoding/csv"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
	"github.com/google/uuid"
	"github.com/guptaharsh13/balkanid-task/initializers"
	"github.com/guptaharsh13/balkanid-task/models"
	"github.com/guptaharsh13/balkanid-task/utils"
	"golang.org/x/crypto/bcrypt"
)

func Signup(c *gin.Context) {

	var body struct {
		Username string `json:"username" validate:"username,required"`
		Email    string `json:"email" validate:"email,required"`
		Password string `json:"password" validate:"password,required"`
	}
	if c.Bind(&body) != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Bad Request",
		})
		return
	}
	err := utils.ValidateStruct(body)
	if err != nil {
		c.JSON(http.StatusBadRequest, utils.ValidationErrorResponse(err))
		return
	}

	var user models.User
	if result := initializers.DB.Take(&user, "username = ?", body.Username); result.RowsAffected > 0 {
		c.JSON(http.StatusConflict, utils.ConflictResponse("Username already taken"))
		return
	}
	if result := initializers.DB.Take(&user, "email = ?", body.Email); result.RowsAffected > 0 {
		c.JSON(http.StatusConflict, utils.ConflictResponse("Email already taken"))
		return
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(body.Password), 10)
	if err != nil {
		c.JSON(http.StatusInternalServerError, utils.InternalServerErrorResponse())
		fmt.Printf("Couldn't hash password: %s", err.Error())
		return
	}

	user = models.User{
		Username: body.Username,
		Email:    body.Email,
		Password: string(hash),
	}
	if result := initializers.DB.Create(&user); result.Error != nil {
		c.JSON(http.StatusInternalServerError, utils.InternalServerErrorResponse())
		fmt.Printf("Couldn't create User: %s", err.Error())
		return
	}
	data := struct {
		User models.User `json:"user"`
	}{
		User: user,
	}
	c.JSON(http.StatusOK, utils.SuccessResponse(data))
}

func Login(c *gin.Context) {

	var body struct {
		Username string `json:"username" validate:"omitempty"`
		Email    string `json:"email" validate:"omitempty,email"`
		Password string `json:"password" validate:"required"`
	}

	if c.Bind(&body) != nil {
		c.JSON(http.StatusBadRequest, utils.BadRequestResponse("Bad Request"))
		return
	}
	err := utils.ValidateStruct(body)
	if err != nil {
		c.JSON(http.StatusBadRequest, utils.ValidationErrorResponse(err))
		return
	}

	var user models.User
	if len(body.Username) > 0 {
		initializers.DB.Take(&user, "username = ?", body.Username)
	} else if len(body.Email) > 0 {
		initializers.DB.Take(&user, "email = ?", body.Email)
	} else {
		c.JSON(http.StatusBadRequest, utils.BadRequestResponse("Username or Email required"))
		return
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(body.Password))
	if err != nil {
		c.JSON(http.StatusUnauthorized, utils.UnauthorizedResponse("Invalid Credentials"))
		return
	}
	if !user.IsActive {
		c.JSON(http.StatusUnauthorized, utils.UnauthorizedResponse("Email not verified"))
		return
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"username": user.Username,
		"exp":      time.Now().Add(time.Hour * 24 * 15).Unix(),
	})

	tokenString, err := token.SignedString([]byte(os.Getenv("JWT_SECRET")))
	if err != nil {
		c.JSON(http.StatusInternalServerError, utils.InternalServerErrorResponse())
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

func VerifyEmail(c *gin.Context) {

	username := c.Param("username")
	if len(strings.TrimSpace(username)) == 0 {
		c.JSON(http.StatusBadRequest, utils.BadRequestResponse("Username is required"))
		return
	}

	var user models.User
	if result := initializers.DB.Take(&user, "username = ?", username); result.Error != nil {
		c.JSON(http.StatusNotFound, utils.NotFoundResponse("User not found"))
		return
	}
	if user.IsActive {
		c.JSON(http.StatusConflict, utils.ConflictResponse("Email already verified"))
		return
	}
	var verifyEmail models.VerifyEmail
	if result := initializers.DB.Take(&verifyEmail, "user_id = ?", user.Username); result.RowsAffected > 0 {

		verifyEmail.Code = uuid.New().String()
		verifyEmail.Expiration = time.Now().Add(time.Hour * 15)
		verifyEmail.IsUsed = false

		if result := initializers.DB.Save(&verifyEmail); result.Error != nil {
			c.JSON(http.StatusInternalServerError, utils.InternalServerErrorResponse())
			fmt.Printf("Couldn't update verifyEmail: %s", result.Error.Error())
			return
		}
	} else {
		verifyEmail = models.VerifyEmail{
			User:       user,
			Code:       uuid.New().String(),
			Expiration: time.Now().Add(time.Hour * 15),
		}
		if result := initializers.DB.Create(&verifyEmail); result.Error != nil {
			c.JSON(http.StatusInternalServerError, utils.InternalServerErrorResponse())
			fmt.Printf("Couldn't create verifyEmail: %s", result.Error.Error())
			return
		}
	}

	data := struct {
		Code string `json:"code"`
	}{
		Code: verifyEmail.Code,
	}
	c.JSON(http.StatusOK, utils.SuccessResponse(data))
}

func ActivateUser(c *gin.Context) {

	username := c.Param("username")
	if len(strings.TrimSpace(username)) == 0 {
		c.JSON(http.StatusBadRequest, utils.BadRequestResponse("Username is required"))
		return
	}
	code := c.Param("code")
	if len(strings.TrimSpace(code)) == 0 {
		c.JSON(http.StatusBadRequest, utils.BadRequestResponse("Verification code is required"))
		return
	}
	var user models.User
	if result := initializers.DB.Take(&user, "username = ?", username); result.Error != nil {
		c.JSON(http.StatusNotFound, utils.NotFoundResponse("User not found"))
		return
	}
	if user.IsActive {
		c.JSON(http.StatusConflict, utils.ConflictResponse("User already active"))
		return
	}

	var verifyEmail models.VerifyEmail
	if result := initializers.DB.Take(&verifyEmail, "user_id = ?", user.Username); result.Error != nil {
		c.JSON(http.StatusBadRequest, utils.BadRequestResponse("Send verification code first"))
		return
	}
	if verifyEmail.IsUsed {
		c.JSON(http.StatusUnauthorized, utils.UnauthorizedResponse("Verification code already used"))
		return
	}
	if verifyEmail.Expiration.Unix() < time.Now().Unix() {
		c.JSON(http.StatusUnauthorized, utils.UnauthorizedResponse("Verification code expired"))
		return
	}
	if verifyEmail.Code != code {
		c.JSON(http.StatusUnauthorized, utils.UnauthorizedResponse("Invalid verification code"))
		return
	}
	if result := initializers.DB.Model(&user).Update("is_active", true); result.Error != nil {
		c.JSON(http.StatusInternalServerError, utils.InternalServerErrorResponse())
		fmt.Printf("Couldn't update user: %s", result.Error.Error())
		return
	}
	verifyEmail.IsUsed = true
	if result := initializers.DB.Save(&verifyEmail); result.Error != nil {
		c.JSON(http.StatusInternalServerError, utils.InternalServerErrorResponse())
		fmt.Printf("Couldn't update verifyEmail: %s", result.Error.Error())
		return
	}
	data := struct {
		User models.User `json:"user"`
	}{
		User: user,
	}
	c.JSON(http.StatusOK, utils.SuccessResponse(data))
}

func DeactivateUser(c *gin.Context) {
	username := c.Param("username")
	if len(strings.TrimSpace(username)) == 0 {
		c.JSON(http.StatusBadRequest, utils.BadRequestResponse("Username is required"))
		return
	}
	var user models.User
	if result := initializers.DB.Take(&user, "username = ?", username); result.Error != nil {
		c.JSON(http.StatusNotFound, utils.NotFoundResponse("User not found"))
		return
	}
	if !user.IsActive {
		c.JSON(http.StatusConflict, utils.ConflictResponse("User already inactive"))
		return
	}
	if result := initializers.DB.Model(&user).Update("is_active", false); result.Error != nil {
		c.JSON(http.StatusInternalServerError, utils.InternalServerErrorResponse())
		fmt.Printf("Couldn't update user: %s", result.Error.Error())
		return
	}
	data := struct {
		User models.User `json:"user"`
	}{
		User: user,
	}
	c.JSON(http.StatusOK, utils.SuccessResponse(data))
}

func GetCurrentUser(c *gin.Context) {

	username, ok := c.Get("username")
	if !ok {
		c.JSON(http.StatusUnauthorized, utils.UnauthorizedResponse("Unauthorized"))
		return
	}
	var user models.User
	if result := initializers.DB.Take(&user, "username = ?", username); result.Error != nil {
		c.JSON(http.StatusNotFound, utils.NotFoundResponse("User not found"))
		return
	}
	data := struct {
		User models.User `json:"user"`
	}{
		User: user,
	}
	c.JSON(http.StatusOK, utils.SuccessResponse(data))
}

func GetUsers(c *gin.Context) {

	var users []models.User
	if result := initializers.DB.Find(&users); result.Error != nil {
		c.JSON(http.StatusInternalServerError, utils.InternalServerErrorResponse())
		fmt.Printf("Couldn't fetch users: %s", result.Error.Error())
		return
	}
	data := struct {
		Users []models.User `json:"users"`
	}{
		Users: users,
	}
	c.JSON(http.StatusOK, utils.SuccessResponse(data))
}

func GetUserByUsername(c *gin.Context) {

	username := c.Param("username")
	if len(strings.TrimSpace(username)) == 0 {
		c.JSON(http.StatusBadRequest, utils.BadRequestResponse("Username is required"))
		return
	}
	var user models.User
	if result := initializers.DB.Take(&user, "username = ?", username); result.Error != nil {
		c.JSON(http.StatusNotFound, utils.NotFoundResponse("User not found"))
		return
	}
	data := struct {
		User models.User `json:"user"`
	}{
		User: user,
	}
	c.JSON(http.StatusOK, utils.SuccessResponse(data))
}

func DeleteUser(c *gin.Context) {

	username := c.Param("username")
	if len(strings.TrimSpace(username)) == 0 {
		c.JSON(http.StatusBadRequest, utils.BadRequestResponse("Username is required"))
		return
	}
	var user models.User
	if result := initializers.DB.Take(&user, "username = ?", username); result.Error != nil {
		c.JSON(http.StatusNotFound, utils.NotFoundResponse("User not found"))
		return
	}
	if result := initializers.DB.Unscoped().Delete(&user); result.Error != nil {
		c.JSON(http.StatusInternalServerError, utils.InternalServerErrorResponse())
		fmt.Printf("Couldn't delete user: %s", result.Error.Error())
		return
	}
	c.JSON(http.StatusOK, utils.SuccessResponse(nil))
}

func BulkUploadUsers(c *gin.Context) {

	file, _, err := c.Request.FormFile("users")
	if err != nil {
		c.JSON(http.StatusBadRequest, utils.BadRequestResponse("File not found (users)"))
		return
	}
	defer file.Close()

	reader := csv.NewReader(file)
	if _, err := reader.Read(); err != nil && err != io.EOF {
		c.JSON(http.StatusBadRequest, utils.BadRequestResponse("Couldn't read CSV file (should have Username,Email,Password as headers)"))
		return
	}

	var users []models.User
	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			c.JSON(http.StatusBadRequest, utils.BadRequestResponse("Invalid CSV format"))
			return
		}
		if len(record) != 3 {
			c.JSON(http.StatusBadRequest, utils.BadRequestResponse("Invalid CSV format"))
			return
		}

		username := record[0]
		email := record[1]
		password := record[2]
		hash, err := bcrypt.GenerateFromPassword([]byte(password), 10)
		if err != nil {
			c.JSON(http.StatusInternalServerError, utils.InternalServerErrorResponse())
			fmt.Printf("Couldn't hash password: %s", err.Error())
			return
		}
		users = append(users, models.User{
			Username: username,
			Email:    email,
			Password: string(hash),
		})
	}
	if result := initializers.DB.Create(&users); result.Error != nil {
		c.JSON(http.StatusInternalServerError, utils.InternalServerErrorResponse())
		fmt.Printf("Couldn't create users: %s", err.Error())
		return
	}
	data := struct {
		Users []models.User `json:"users"`
	}{
		Users: users,
	}
	c.JSON(http.StatusOK, utils.SuccessResponse(data))
}
