package redirect

import (
	"github.com/gin-gonic/gin"
	"github.com/namatery/system-design-lab/00-url-shortener/internal/analytics"
	"time"
)

type ClickPublisher interface {
	TryPublish(analytics.ClickEventDto) bool
}

type RedirectController struct {
	service   RedirectService
	collector ClickPublisher
}

func NewRedirectController(service RedirectService) *RedirectController {
	return &RedirectController{
		service: service,
	}
}

func NewRedirectControllerWithAnalytics(service RedirectService, collector ClickPublisher) *RedirectController {
	return &RedirectController{service: service, collector: collector}
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

	if rc.collector != nil {
		rc.collector.TryPublish(analytics.ClickEventDto{
			Alias: input.Alias, Timestamp: time.Now().UTC(),
			Referrer: c.Request.Referer(), UserAgent: c.Request.UserAgent(),
		})
	}
	c.Redirect(302, longUrl)
}
