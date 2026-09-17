package redirect

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/gocql/gocql"
)

type RedirectService interface {
	ResolveAlias(c *gin.Context, input RedirectDto) (string, error)
}

type RedirectServiceImpl struct {
	session *gocql.Session
}

func NewRedirectService(session *gocql.Session) RedirectService {
	return &RedirectServiceImpl{
		session: session,
	}
}

func (s *RedirectServiceImpl) ResolveAlias(c *gin.Context, input RedirectDto) (string, error) {
	var longUrl string

	err := s.session.Query(
		`SELECT long_url FROM links WHERE alias = ?`,
		input.Alias,
	).WithContext(c.Request.Context()).Scan(&longUrl)

	if err != nil {
		return "", fmt.Errorf("resolve alias: %w", err)
	}

	return longUrl, nil
}

var _ RedirectService = (*RedirectServiceImpl)(nil)
