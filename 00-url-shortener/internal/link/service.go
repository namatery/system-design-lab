package link

import (
	"errors"
	"fmt"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gocql/gocql"
	"github.com/namatery/system-design-lab/00-url-shortener/internal/shared/config"
	"github.com/namatery/system-design-lab/00-url-shortener/internal/shared/model"
)

var (
	ErrAliasExists = errors.New("link alias already exists")
	ErrNotFound    = errors.New("link not found")
)

type LinkService interface {
	Create(c *gin.Context, input CreateLinkDto) (string, error)
}

type LinkServiceImpl struct {
	session *gocql.Session
}

func NewLinkService(session *gocql.Session) LinkService {
	return &LinkServiceImpl{
		session: session,
	}
}

func (s *LinkServiceImpl) Create(c *gin.Context, input CreateLinkDto) (string, error) {
	// Generate a short alias using the first 6 characters of a UUID
	alias := gocql.TimeUUID().String()[:6]

	expiresAt, err := time.Parse("2006-01-02", input.ExpirationDate)
	if err != nil {
		return "", fmt.Errorf("invalid expiration date: %w", err)
	}

	link := model.Link{
		Alias:     alias,
		LongURL:   input.URL,
		CreatedAt: time.Now().UTC(),
		ExpiresAt: &expiresAt,
	}

	existing := make(map[string]interface{})
	applied, err := s.session.Query(`
		INSERT INTO links (alias, long_url, created_at, expires_at)
		VALUES (?, ?, ?, ?)
		IF NOT EXISTS`,
		link.Alias,
		link.LongURL,
		link.CreatedAt,
		link.ExpiresAt,
	).WithContext(c.Request.Context()).MapScanCAS(existing)
	if err != nil {
		return "", fmt.Errorf("create link: %w", err)
	}
	if !applied {
		return "", ErrAliasExists
	}

	return config.GetConfig().PublicURL + "/" + link.Alias, nil
}

var _ LinkService = (*LinkServiceImpl)(nil)
