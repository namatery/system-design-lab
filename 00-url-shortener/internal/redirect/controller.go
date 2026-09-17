package redirect

import "github.com/gin-gonic/gin"

type RedirectController struct {
	service RedirectService
}

func NewRedirectController(service RedirectService) *RedirectController {
	return &RedirectController{
		service: service,
	}
}

func (rc *RedirectController) Redirect(c *gin.Context) {
	var input RedirectDto
	if err := c.ShouldBindUri(&input); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	longUrl, err := rc.service.ResolveAlias(c, input)
	if err != nil { // Handle error appropriately
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	c.Redirect(302, longUrl)
}
