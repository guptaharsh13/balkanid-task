package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/guptaharsh13/balkanid-task/utils"
)

func Health(c *gin.Context) {
	c.JSON(http.StatusOK, utils.SuccessResponse(nil))
}
