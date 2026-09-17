package analytics

import "time"

type ClickEventDto struct {
	Alias     string
	Timestamp time.Time
	Referrer  string
	UserAgent string
}
