package controllers

import (
	"encoding/json"
	"io"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/guptaharsh13/balkanid-task/initializers"
	"github.com/guptaharsh13/balkanid-task/models"
	"github.com/guptaharsh13/balkanid-task/utils"
)

func CreateRole(c *gin.Context) {

	var body struct {
		Name        string   `json:"name" validate:"required"`
		Description string   `json:"description"`
		Users       []string `json:"users"`
		Permissions []string `json:"permissions"`
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

	if result := initializers.DB.Take(&models.Role{}, "name = ?", body.Name); result.RowsAffected > 0 {
		c.JSON(http.StatusConflict, utils.ConflictResponse("Role already exists"))
		return
	}

	var users []models.User
	result := initializers.DB.Where("username IN ?", body.Users).Find(&users)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, utils.InternalServerErrorResponse())
		return
	}
	if result.RowsAffected != int64(len(body.Users)) {
		c.JSON(http.StatusBadRequest, utils.BadRequestResponse("Couldn't find all users"))
		return
	}

	var permissions []models.Permission
	result = initializers.DB.Where("name IN ?", body.Permissions).Find(&permissions)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, utils.InternalServerErrorResponse())
		return
	}
	if result.RowsAffected != int64(len(body.Permissions)) {
		c.JSON(http.StatusBadRequest, utils.BadRequestResponse("Permissions not found"))
		return
	}

	role := models.Role{
		Name:        body.Name,
		Description: body.Description,
		Users:       users,
		Permissions: permissions,
	}
	if result := initializers.DB.Create(&role); result.Error != nil {
		c.JSON(http.StatusInternalServerError, utils.InternalServerErrorResponse())
		return
	}
	data := struct {
		Role models.Role `json:"role"`
	}{
		Role: role,
	}
	c.JSON(http.StatusOK, utils.SuccessResponse(data))
}

func GetRoles(c *gin.Context) {

	var roles []models.Role
	if result := initializers.DB.Preload("Users").Preload("Permissions").Find(&roles); result.Error != nil {
		c.JSON(http.StatusInternalServerError, utils.InternalServerErrorResponse())
		return
	}
	data := struct {
		Roles []models.Role `json:"roles"`
	}{
		Roles: roles,
	}
	c.JSON(http.StatusOK, utils.SuccessResponse(data))
}

func GetRoleByName(c *gin.Context) {

	name := c.Param("name")
	if len(strings.TrimSpace(name)) == 0 {
		c.JSON(http.StatusBadRequest, utils.BadRequestResponse("Name is required"))
		return
	}
	var role models.Role
	if result := initializers.DB.Preload("Users").Preload("Permissions").Take(&role, "name = ?", name); result.Error != nil {
		c.JSON(http.StatusNotFound, utils.NotFoundResponse("Couldn't find role"))
		return
	}
	data := struct {
		Role models.Role `json:"role"`
	}{
		Role: role,
	}
	c.JSON(http.StatusOK, utils.SuccessResponse(data))
}

func UpdateRolePut(c *gin.Context) {

	name := c.Param("name")
	if len(strings.TrimSpace(name)) == 0 {
		c.JSON(http.StatusBadRequest, utils.BadRequestResponse("Name is required"))
		return
	}
	var role models.Role
	if result := initializers.DB.Preload("Users").Preload("Permissions").Take(&role, "name = ?", name); result.Error != nil {
		c.JSON(http.StatusNotFound, utils.NotFoundResponse("Couldn't find role"))
		return
	}

	var body struct {
		Name        string   `json:"name"`
		Description string   `json:"description"`
		Users       []string `json:"users"`
		Permissions []string `json:"permissions"`
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

	if body.Name != role.Name {
		if result := initializers.DB.Take(&models.Role{}, "name = ?", body.Name); result.RowsAffected > 0 {
			c.JSON(http.StatusConflict, utils.ConflictResponse("Role already exists"))
			return
		}
	}

	var users []models.User
	result := initializers.DB.Where("username IN ?", body.Users).Find(&users)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, utils.InternalServerErrorResponse())
		return
	}
	if result.RowsAffected != int64(len(body.Users)) {
		c.JSON(http.StatusBadRequest, utils.BadRequestResponse("Couldn't find all users"))
		return
	}

	var permissions []models.Permission
	result = initializers.DB.Where("name IN ?", body.Permissions).Find(&permissions)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, utils.InternalServerErrorResponse())
		return
	}
	if result.RowsAffected != int64(len(body.Permissions)) {
		c.JSON(http.StatusBadRequest, utils.BadRequestResponse("Permissio"))
		return
	}

	if len(body.Name) > 0 {
		role.Name = body.Name
	}
	role.Description = body.Description
	role.Users = users
	role.Permissions = permissions
	if result := initializers.DB.Save(&role); result.Error != nil {
		c.JSON(http.StatusInternalServerError, utils.InternalServerErrorResponse())
		return
	}
	data := struct {
		Role models.Role `json:"role"`
	}{
		Role: role,
	}
	c.JSON(http.StatusOK, utils.SuccessResponse(data))
}

func UpdateRolePatch(c *gin.Context) {

	name := c.Param("name")
	if len(strings.TrimSpace(name)) == 0 {
		c.JSON(http.StatusBadRequest, utils.BadRequestResponse("Name is required"))
		return
	}
	var role models.Role
	if result := initializers.DB.Preload("Users").Preload("Permissions").Take(&role, "name = ?", name); result.Error != nil {
		c.JSON(http.StatusNotFound, utils.NotFoundResponse("Couldn't find role"))
		return
	}

	requestBytes, err := io.ReadAll(c.Request.Body)
	defer c.Request.Body.Close()
	if err != nil {
		c.JSON(http.StatusInternalServerError, utils.InternalServerErrorResponse())
		return
	}
	body := make(map[string]interface{})
	if err = json.Unmarshal(requestBytes, &body); err != nil {
		c.JSON(http.StatusInternalServerError, utils.InternalServerErrorResponse())
		return
	}

	if name, ok := body["name"]; ok {
		if name != role.Name {
			if result := initializers.DB.Take(&models.Role{}, "name = ?", name); result.RowsAffected > 0 {
				c.JSON(http.StatusConflict, utils.ConflictResponse("Role already exists"))
				return
			}
		}
		role.Name = name.(string)
	}

	if description, ok := body["description"]; ok {
		role.Description = description.(string)
	}

	if value, ok := body["users"]; ok {
		usernames := value.([]string)
		var users []models.User
		result := initializers.DB.Where("username IN ?", usernames).Find(&users)
		if result.Error != nil {
			c.JSON(http.StatusInternalServerError, utils.InternalServerErrorResponse())
			return
		}
		if result.RowsAffected != int64(len(usernames)) {
			c.JSON(http.StatusBadRequest, utils.BadRequestResponse("Couldn't find all users"))
			return
		}
		role.Users = users
	}

	if value, ok := body["permissions"]; ok {
		permissionNames := value.([]string)
		var permissions []models.Permission
		result := initializers.DB.Where("name IN ?", permissionNames).Find(&permissions)
		if result.Error != nil {
			c.JSON(http.StatusInternalServerError, utils.InternalServerErrorResponse())
			return
		}
		if result.RowsAffected != int64(len(permissionNames)) {
			c.JSON(http.StatusBadRequest, utils.BadRequestResponse("Permissio"))
			return
		}
		role.Permissions = permissions
	}

	if result := initializers.DB.Save(&role); result.Error != nil {
		c.JSON(http.StatusInternalServerError, utils.InternalServerErrorResponse())
		return
	}
	data := struct {
		Role models.Role `json:"role"`
	}{
		Role: role,
	}
	c.JSON(http.StatusOK, utils.SuccessResponse(data))
}

func DeleteRole(c *gin.Context) {

	name := c.Param("name")
	if len(strings.TrimSpace(name)) == 0 {
		c.JSON(http.StatusBadRequest, utils.BadRequestResponse("Name is required"))
		return
	}
	var role models.Role
	if result := initializers.DB.Take(&role, "name = ?", name); result.Error != nil {
		c.JSON(http.StatusNotFound, utils.NotFoundResponse("Couldn't find role"))
		return
	}
	if result := initializers.DB.Unscoped().Delete(&role); result.Error != nil {
		c.JSON(http.StatusInternalServerError, utils.InternalServerErrorResponse())
		return
	}
	c.JSON(http.StatusOK, utils.SuccessResponse(nil))
}

func AddPermissionsToRole(c *gin.Context) {

	name := c.Param("name")
	if len(strings.TrimSpace(name)) == 0 {
		c.JSON(http.StatusBadRequest, utils.BadRequestResponse("Role name is required"))
		return
	}
	var role models.Role
	if result := initializers.DB.Preload("Permissions").Take(&role, "name = ?", name); result.Error != nil {
		c.JSON(http.StatusNotFound, utils.NotFoundResponse("Couldn't find role"))
		return
	}

	var body struct {
		Permissions []string `json:"permissions"`
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
	var permissions []models.Permission
	result := initializers.DB.Where("name IN ?", body.Permissions).Find(&permissions)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, utils.InternalServerErrorResponse())
		return
	}
	if result.RowsAffected != int64(len(body.Permissions)) {
		c.JSON(http.StatusBadRequest, utils.BadRequestResponse("Permissions not found"))
		return
	}

	role.Permissions = append(role.Permissions, permissions...)
	if result := initializers.DB.Save(&role); result.Error != nil {
		c.JSON(http.StatusInternalServerError, utils.InternalServerErrorResponse())
		return
	}
	data := struct {
		Role models.Role `json:"role"`
	}{
		Role: role,
	}
	c.JSON(http.StatusOK, utils.SuccessResponse(data))
}

func AssignRoleToUser(c *gin.Context) {

	name := c.Param("name")
	if len(strings.TrimSpace(name)) == 0 {
		c.JSON(http.StatusBadRequest, utils.BadRequestResponse("Role name is required"))
		return
	}
	var role models.Role
	if result := initializers.DB.Preload("Users").Take(&role, "name = ?", name); result.Error != nil {
		c.JSON(http.StatusNotFound, utils.NotFoundResponse("Couldn't find role"))
		return
	}

	var body struct {
		Users []string `json:"users"`
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
	var users []models.User
	result := initializers.DB.Where("username IN ?", body.Users).Find(&users)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, utils.InternalServerErrorResponse())
		return
	}
	if result.RowsAffected != int64(len(body.Users)) {
		c.JSON(http.StatusBadRequest, utils.BadRequestResponse("Couldn't find all users"))
		return
	}

	role.Users = append(role.Users, users...)
	if result := initializers.DB.Save(&role); result.Error != nil {
		c.JSON(http.StatusInternalServerError, utils.InternalServerErrorResponse())
		return
	}
	data := struct {
		Role models.Role `json:"role"`
	}{
		Role: role,
	}
	c.JSON(http.StatusOK, utils.SuccessResponse(data))
}

func GetUsersByRole(c *gin.Context) {
	name := c.Param("name")
	if len(strings.TrimSpace(name)) == 0 {
		c.JSON(http.StatusBadRequest, utils.BadRequestResponse("Role name is required"))
		return
	}
	var role models.Role
	if result := initializers.DB.Preload("Users").Take(&role, "name = ?", name); result.Error != nil {
		c.JSON(http.StatusNotFound, utils.NotFoundResponse("Couldn't find role"))
		return
	}

	var users []models.User
	if err := initializers.DB.Model(&role).Association("Users").Find(&users); err != nil {
		c.JSON(http.StatusInternalServerError, utils.InternalServerErrorResponse())
		return
	}
	data := struct {
		Users []models.User `json:"users"`
	}{
		Users: users,
	}
	c.JSON(http.StatusOK, utils.SuccessResponse(data))
}
