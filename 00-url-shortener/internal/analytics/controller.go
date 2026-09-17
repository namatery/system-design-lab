package analytics

import (
	"context"
	"sync"

	"go.uber.org/zap"
)

type AnalyticsCollector struct {
	queue   chan ClickEventDto
	service AnalyticsService
	once    sync.Once
	workers sync.WaitGroup
	logger  *zap.Logger
}

func NewAnalyticsCollector(queueSize int, service AnalyticsService, logger *zap.Logger) *AnalyticsCollector {
	return &AnalyticsCollector{
		queue:   make(chan ClickEventDto, queueSize),
		service: service,
		logger:  logger,
	}
}

func (c *AnalyticsCollector) TryPublish(event ClickEventDto) bool {
	select {
	case c.queue <- event:
		return true
	default:
		return false
	}
}

func (c *AnalyticsCollector) Start(ctx context.Context, workers int) {
	if workers <= 0 {
		panic("analytics: workers must be positive")
	}
	c.once.Do(func() {
		c.workers.Add(workers)
		for i := 0; i < workers; i++ {
			go func() {
				defer c.workers.Done()
				for {
					if ctx.Err() != nil {
						return
					}
					select {
					case <-ctx.Done():
						return
					case event := <-c.queue:
						if err := c.service.RecordClick(ctx, event); err != nil {
							c.logger.Warn("Failed to record click", zap.String("alias", event.Alias), zap.Error(err))
						}
					}
				}
			}()
		}
	})
}

func (c *AnalyticsCollector) Wait() { c.workers.Wait() }
