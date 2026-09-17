package link

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type LinkController struct {
	service LinkService
}

func NewLinkController(service LinkService) *LinkController {
	return &LinkController{
		service: service,
	}
}

func (lc *LinkController) CreateShortenLink(c *gin.Context) {
	var input CreateLinkDto
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	shortUrl, err := lc.service.Create(c, input)
	if err != nil { // Handle error appropriately
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"short_url": shortUrl})
}
