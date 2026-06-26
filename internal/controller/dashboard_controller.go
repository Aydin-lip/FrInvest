package controller

import (
	"net/http"
	"recruitment-api/internal/dto"
	"recruitment-api/internal/service"

	"github.com/gin-gonic/gin"
)

type DashboardController struct {
	userService service.UserService
}

func NewDashboardController(userService service.UserService) *DashboardController {
	return &DashboardController{userService: userService}
}

func (dc *DashboardController) GetStatusPercentages(c *gin.Context) {
	stats, err := dc.userService.GetStatusPercentages()
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{Message: err.Error()})
		return
	}

	c.JSON(http.StatusOK, stats)
}
