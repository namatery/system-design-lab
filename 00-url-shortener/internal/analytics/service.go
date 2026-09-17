package analytics

import (
	"context"
	"fmt"
	"time"

	"github.com/gocql/gocql"
)

type AnalyticsService interface {
	RecordClick(context.Context, ClickEventDto) error
}

type AnalyticsServiceImpl struct {
	session *gocql.Session
}

func NewAnalyticsService(session *gocql.Session) AnalyticsService {
	return &AnalyticsServiceImpl{session: session}
}

func (s *AnalyticsServiceImpl) RecordClick(ctx context.Context, event ClickEventDto) error {
	ctx, cancel := context.WithTimeout(ctx, 2*time.Second)
	defer cancel()

	timestamp := event.Timestamp.UTC()
	err := s.session.Query(
		`INSERT INTO click_events 
		(alias, day, event_id, clicked_at, referrer, user_agent)
  		VALUES (?, ?, ?, ?, ?, ?)`,
		event.Alias,
		timestamp.Format("2006-01-02"),
		gocql.UUIDFromTime(timestamp),
		timestamp,
		event.Referrer,
		event.UserAgent,
	).WithContext(ctx).Exec()

	if err != nil {
		return fmt.Errorf("record click: %w", err)
	}

	return nil
}

var _ AnalyticsService = (*AnalyticsServiceImpl)(nil)
