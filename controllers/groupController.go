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

func CreateGroup(c *gin.Context) {

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

	if result := initializers.DB.Take(&models.Group{}, "name = ?", body.Name); result.RowsAffected > 0 {
		c.JSON(http.StatusConflict, utils.ConflictResponse("Group already exists"))
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

	group := models.Group{
		Name:        body.Name,
		Description: body.Description,
		Users:       users,
		Permissions: permissions,
	}
	if result := initializers.DB.Create(&group); result.Error != nil {
		c.JSON(http.StatusInternalServerError, utils.InternalServerErrorResponse())
		return
	}
	data := struct {
		Group models.Group `json:"group"`
	}{
		Group: group,
	}
	c.JSON(http.StatusOK, utils.SuccessResponse(data))
}

func GetGroups(c *gin.Context) {

	var groups []models.Group
	if result := initializers.DB.Preload("Users").Preload("Permissions").Find(&groups); result.Error != nil {
		c.JSON(http.StatusInternalServerError, utils.InternalServerErrorResponse())
		return
	}
	data := struct {
		Groups []models.Group `json:"groups"`
	}{
		Groups: groups,
	}
	c.JSON(http.StatusOK, utils.SuccessResponse(data))
}

func GetGroupByName(c *gin.Context) {

	name := c.Param("name")
	if len(strings.TrimSpace(name)) == 0 {
		c.JSON(http.StatusBadRequest, utils.BadRequestResponse("Name is required"))
		return
	}
	var group models.Group
	if result := initializers.DB.Preload("Users").Preload("Permissions").Take(&group, "name = ?", name); result.Error != nil {
		c.JSON(http.StatusNotFound, utils.NotFoundResponse("Couldn't find group"))
		return
	}
	data := struct {
		Group models.Group `json:"group"`
	}{
		Group: group,
	}
	c.JSON(http.StatusOK, utils.SuccessResponse(data))
}

func UpdateGroupPut(c *gin.Context) {

	name := c.Param("name")
	if len(strings.TrimSpace(name)) == 0 {
		c.JSON(http.StatusBadRequest, utils.BadRequestResponse("Name is required"))
		return
	}
	var group models.Group
	if result := initializers.DB.Preload("Users").Preload("Permissions").Take(&group, "name = ?", name); result.Error != nil {
		c.JSON(http.StatusNotFound, utils.NotFoundResponse("Couldn't find group"))
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

	if body.Name != group.Name {
		if result := initializers.DB.Take(&models.Group{}, "name = ?", body.Name); result.RowsAffected > 0 {
			c.JSON(http.StatusConflict, utils.ConflictResponse("Group already exists"))
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
		group.Name = body.Name
	}
	group.Description = body.Description
	group.Users = users
	group.Permissions = permissions
	if result := initializers.DB.Save(&group); result.Error != nil {
		c.JSON(http.StatusInternalServerError, utils.InternalServerErrorResponse())
		return
	}
	data := struct {
		Group models.Group `json:"group"`
	}{
		Group: group,
	}
	c.JSON(http.StatusOK, utils.SuccessResponse(data))
}

func UpdateGroupPatch(c *gin.Context) {

	name := c.Param("name")
	if len(strings.TrimSpace(name)) == 0 {
		c.JSON(http.StatusBadRequest, utils.BadRequestResponse("Name is required"))
		return
	}
	var group models.Group
	if result := initializers.DB.Preload("Users").Preload("Permissions").Take(&group, "name = ?", name); result.Error != nil {
		c.JSON(http.StatusNotFound, utils.NotFoundResponse("Couldn't find group"))
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
		if name != group.Name {
			if result := initializers.DB.Take(&models.Group{}, "name = ?", name); result.RowsAffected > 0 {
				c.JSON(http.StatusConflict, utils.ConflictResponse("Group already exists"))
				return
			}
		}
		group.Name = name.(string)
	}

	if description, ok := body["description"]; ok {
		group.Description = description.(string)
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
		group.Users = users
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
		group.Permissions = permissions
	}

	if result := initializers.DB.Save(&group); result.Error != nil {
		c.JSON(http.StatusInternalServerError, utils.InternalServerErrorResponse())
		return
	}
	data := struct {
		Group models.Group `json:"group"`
	}{
		Group: group,
	}
	c.JSON(http.StatusOK, utils.SuccessResponse(data))
}

func DeleteGroup(c *gin.Context) {

	name := c.Param("name")
	if len(strings.TrimSpace(name)) == 0 {
		c.JSON(http.StatusBadRequest, utils.BadRequestResponse("Name is required"))
		return
	}
	var group models.Group
	if result := initializers.DB.Take(&group, "name = ?", name); result.Error != nil {
		c.JSON(http.StatusNotFound, utils.NotFoundResponse("Couldn't find group"))
		return
	}
	if result := initializers.DB.Unscoped().Delete(&group); result.Error != nil {
		c.JSON(http.StatusInternalServerError, utils.InternalServerErrorResponse())
		return
	}
	c.JSON(http.StatusOK, utils.SuccessResponse(nil))
}

func AddPermissionsToGroup(c *gin.Context) {

	name := c.Param("name")
	if len(strings.TrimSpace(name)) == 0 {
		c.JSON(http.StatusBadRequest, utils.BadRequestResponse("Group name is required"))
		return
	}
	var group models.Group
	if result := initializers.DB.Preload("Permissions").Take(&group, "name = ?", name); result.Error != nil {
		c.JSON(http.StatusNotFound, utils.NotFoundResponse("Couldn't find group"))
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

	group.Permissions = append(group.Permissions, permissions...)
	if result := initializers.DB.Save(&group); result.Error != nil {
		c.JSON(http.StatusInternalServerError, utils.InternalServerErrorResponse())
		return
	}
	data := struct {
		Group models.Group `json:"group"`
	}{
		Group: group,
	}
	c.JSON(http.StatusOK, utils.SuccessResponse(data))
}

func AddUsersToGroup(c *gin.Context) {

	name := c.Param("name")
	if len(strings.TrimSpace(name)) == 0 {
		c.JSON(http.StatusBadRequest, utils.BadRequestResponse("Group name is required"))
		return
	}
	var group models.Group
	if result := initializers.DB.Preload("Users").Take(&group, "name = ?", name); result.Error != nil {
		c.JSON(http.StatusNotFound, utils.NotFoundResponse("Couldn't find group"))
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

	group.Users = append(group.Users, users...)
	if result := initializers.DB.Save(&group); result.Error != nil {
		c.JSON(http.StatusInternalServerError, utils.InternalServerErrorResponse())
		return
	}
	data := struct {
		Group models.Group `json:"group"`
	}{
		Group: group,
	}
	c.JSON(http.StatusOK, utils.SuccessResponse(data))
}

func GetUsersByGroup(c *gin.Context) {
	name := c.Param("name")
	if len(strings.TrimSpace(name)) == 0 {
		c.JSON(http.StatusBadRequest, utils.BadRequestResponse("Group name is required"))
		return
	}
	var group models.Group
	if result := initializers.DB.Preload("Users").Take(&group, "name = ?", name); result.Error != nil {
		c.JSON(http.StatusNotFound, utils.NotFoundResponse("Couldn't find group"))
		return
	}

	var users []models.User
	if err := initializers.DB.Model(&group).Association("Users").Find(&users); err != nil {
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
