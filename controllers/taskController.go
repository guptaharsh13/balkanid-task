package controllers

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/guptaharsh13/balkanid-task/initializers"
	"github.com/guptaharsh13/balkanid-task/models"
	"github.com/guptaharsh13/balkanid-task/utils"
)

func CreateTask(c *gin.Context) {

	var body struct {
		Name        string   `json:"name" validate:"required"`
		Description string   `json:"description"`
		Creator     string   `json:"creator" validate:"required"`
		Asignees    []string `json:"asignees"`
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

	var creator models.User
	result := initializers.DB.First(&creator, "username = ?", body.Creator)
	if result.Error != nil {
		c.JSON(http.StatusBadRequest, utils.BadRequestResponse(fmt.Sprintf("Could not find user %s", body.Creator)))
		return
	}
	fmt.Println(creator)

	var asignees []models.User
	for _, username := range body.Asignees {
		var asignee models.User
		result := initializers.DB.First(&asignee, "username = ?", username)
		if result.Error != nil {
			c.JSON(http.StatusBadRequest, utils.BadRequestResponse(fmt.Sprintf("Could not find user %s", username)))
			return
		}
		asignees = append(asignees, asignee)
	}

	task := models.Task{
		Name:        body.Name,
		Description: body.Description,
		Creator:     creator.Username,
		Asignees:    asignees,
	}
	result = initializers.DB.Create(&task)
	if result.Error != nil {
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
	result := initializers.DB.Find(&tasks)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Internal Server Error",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"tasks": tasks,
	})
}
