package controllers

import (
	"encoding/csv"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/guptaharsh13/balkanid-task/initializers"
	"github.com/guptaharsh13/balkanid-task/models"
	"github.com/guptaharsh13/balkanid-task/utils"
)

func CreateTask(c *gin.Context) {

	username, ok := c.Get("username")
	if !ok {
		c.JSON(http.StatusUnauthorized, utils.UnauthorizedResponse("Unauthorized"))
		return
	}

	var body struct {
		Name        string   `json:"name" validate:"required"`
		Description string   `json:"description"`
		Asignees    []string `json:"asignees"`
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

	var creator models.User
	if result := initializers.DB.Take(&creator, "username = ?", username); result.Error != nil {
		c.JSON(http.StatusInternalServerError, utils.InternalServerErrorResponse())
		return
	}

	var asignees []models.User
	result := initializers.DB.Find(&asignees, "username IN ?", body.Asignees)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, utils.InternalServerErrorResponse())
		return
	}
	if result.RowsAffected != int64(len(asignees)) {
		c.JSON(http.StatusBadRequest, utils.BadRequestResponse("Couldn't find all users"))
		return
	}

	task := models.Task{
		Name:        body.Name,
		Description: body.Description,
		Creator:     creator.Username,
		Asignees:    asignees,
	}
	if result := initializers.DB.Create(&task); result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Internal Server Error",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"task": task,
	})
}

func GetTasks(c *gin.Context) {

	var tasks []models.Task
	if result := initializers.DB.Find(&tasks); result.Error != nil {
		c.JSON(http.StatusInternalServerError, utils.InternalServerErrorResponse())
		return
	}

	data := struct {
		Tasks []models.Task `json:"tasks"`
	}{
		Tasks: tasks,
	}
	c.JSON(http.StatusOK, utils.SuccessResponse(data))
}

func GetTaskByID(c *gin.Context) {

	id := c.Param("id")
	if len(strings.TrimSpace(id)) == 0 {
		c.JSON(http.StatusBadRequest, utils.BadRequestResponse("ID is required"))
		return
	}
	var task models.Task
	if result := initializers.DB.Preload("Asignees").Take(&task, "id = ?", id); result.Error != nil {
		c.JSON(http.StatusNotFound, utils.NotFoundResponse(fmt.Sprintf("Couldn't find task with id %s", id)))
		return
	}
	data := struct {
		Task models.Task `json:"task"`
	}{
		Task: task,
	}
	c.JSON(http.StatusOK, utils.SuccessResponse(data))
}

func DeleteTask(c *gin.Context) {

	id := c.Param("id")
	if len(strings.TrimSpace(id)) == 0 {
		c.JSON(http.StatusBadRequest, utils.BadRequestResponse("ID is required"))
		return
	}
	var task models.Task
	if result := initializers.DB.Take(&task, "id = ?", id); result.Error != nil {
		c.JSON(http.StatusNotFound, utils.NotFoundResponse(fmt.Sprintf("Couldn't find task with id %s", id)))
		return
	}
	username, ok := c.Get("username")
	if !ok {
		c.JSON(http.StatusUnauthorized, utils.UnauthorizedResponse("Unauthorized"))
		return
	}
	isAdmin, ok := c.Get("is_admin")
	if !ok {
		c.JSON(http.StatusUnauthorized, utils.UnauthorizedResponse("Unauthorized"))
	}
	if username != task.Creator && !isAdmin.(bool) {
		c.JSON(http.StatusForbidden, utils.ForbiddenResponse("Forbidden"))
		return
	}

	if result := initializers.DB.Delete(&task); result.Error != nil {
		c.JSON(http.StatusInternalServerError, utils.InternalServerErrorResponse())
		return
	}
	c.JSON(http.StatusOK, utils.SuccessResponse(nil))
}

func BulkUploadTasks(c *gin.Context) {

	username, ok := c.Get("username")
	if !ok {
		c.JSON(http.StatusUnauthorized, utils.UnauthorizedResponse("Unauthorized"))
		return
	}
	var user models.User
	if result := initializers.DB.Take(&user, "username = ?", username); result.Error != nil {
		c.JSON(http.StatusInternalServerError, utils.InternalServerErrorResponse())
		return
	}
	file, _, err := c.Request.FormFile("tasks")
	if err != nil {
		c.JSON(http.StatusBadRequest, utils.BadRequestResponse("File not found (tasks)"))
		return
	}
	defer file.Close()

	reader := csv.NewReader(file)
	if _, err := reader.Read(); err != nil && err != io.EOF {
		c.JSON(http.StatusBadRequest, utils.BadRequestResponse("Couldn't read CSV file (should have Name,Description as headers)"))
		return
	}

	var tasks []models.Task
	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			c.JSON(http.StatusBadRequest, utils.BadRequestResponse("Invalid CSV format"))
			return
		}
		if len(record) != 2 {
			c.JSON(http.StatusBadRequest, utils.BadRequestResponse("Invalid CSV format"))
			return
		}

		name := record[0]
		description := record[1]

		tasks = append(tasks, models.Task{
			Name:        name,
			Description: description,
			Creator:     user.Username,
		})
	}
	if result := initializers.DB.Create(&tasks); result.Error != nil {
		c.JSON(http.StatusInternalServerError, utils.InternalServerErrorResponse())
		fmt.Printf("Couldn't create tasks: %s", err.Error())
		return
	}
	data := struct {
		Tasks []models.Task `json:"tasks"`
	}{
		Tasks: tasks,
	}
	c.JSON(http.StatusOK, utils.SuccessResponse(data))
}

func AssignTaskToUsers(c *gin.Context) {

	id := c.Param("id")
	if len(strings.TrimSpace(id)) == 0 {
		c.JSON(http.StatusBadRequest, utils.BadRequestResponse("ID is required"))
		return
	}
	var task models.Task
	if result := initializers.DB.Take(&task, "id = ?", id); result.Error != nil {
		c.JSON(http.StatusNotFound, utils.NotFoundResponse(fmt.Sprintf("Couldn't find task with id %s", id)))
		return
	}
	username, ok := c.Get("username")
	if !ok {
		c.JSON(http.StatusUnauthorized, utils.UnauthorizedResponse("Unauthorized"))
		return
	}
	isAdmin, ok := c.Get("is_admin")
	if !ok {
		c.JSON(http.StatusUnauthorized, utils.UnauthorizedResponse("Unauthorized"))
	}
	if username != task.Creator && !isAdmin.(bool) {
		c.JSON(http.StatusForbidden, utils.ForbiddenResponse("Forbidden"))
		return
	}

	var body struct {
		Asignees []string `json:"asignees"`
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
	var asignees []models.User
	result := initializers.DB.Find(&asignees, "username IN ?", body.Asignees)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, utils.InternalServerErrorResponse())
		return
	}
	if result.RowsAffected != int64(len(asignees)) {
		c.JSON(http.StatusBadRequest, utils.BadRequestResponse("Couldn't find all users"))
		return
	}

	task.Asignees = asignees
	if result := initializers.DB.Save(&task); result.Error != nil {
		c.JSON(http.StatusInternalServerError, utils.InternalServerErrorResponse())
		return
	}
	data := struct {
		Task models.Task `json:"task"`
	}{
		Task: task,
	}
	c.JSON(http.StatusOK, utils.SuccessResponse(data))
}
