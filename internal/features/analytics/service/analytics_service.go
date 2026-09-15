package analytics_service

import (
	"context"
	"log/slog"

	"github.com/loundxr/expense-tracker/internal/core/domain"
)

type AnalyticsService struct {
	analyticsRepository AnalyticsRepository
	accChecker          AccountChecker
	logger              *slog.Logger
}

type AnalyticsRepository interface {
	Summary(ctx context.Context, filter domain.SummaryFilter) (domain.Summary, error)
}

type AccountChecker interface {
	HasAccess(ctx context.Context, uid, accountID int64) (bool, error)
}

func NewAnalyticsService(ar AnalyticsRepository, ac AccountChecker, l *slog.Logger) *AnalyticsService {
	return &AnalyticsService{
		analyticsRepository: ar,
		accChecker:          ac,
		logger:              l,
	}
}
