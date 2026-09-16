package analytics_transport_http

import (
	"context"
	"log/slog"

	"github.com/go-chi/chi/v5"
	"github.com/loundxr/expense-tracker/internal/core/domain"
)

type AnalyticsHTTPHandler struct {
	analyticsService AnalyticsService
	logger           *slog.Logger
}

type AnalyticsService interface {
	Summary(ctx context.Context, filter domain.SummaryFilter) (domain.Summary, error)
	CategoriesBreakdown(ctx context.Context, filter domain.SummaryFilter) ([]domain.CategoryBreakdown, error)
	Trends(ctx context.Context, filter domain.TrendsFilter) ([]domain.ExpenseTrendPoint, error)
	MembersContribution(ctx context.Context, filter domain.SummaryFilter) ([]domain.MemberContribution, error)
}

func NewStatisticsHTTPHandler(as AnalyticsService, l *slog.Logger) *AnalyticsHTTPHandler {
	return &AnalyticsHTTPHandler{
		analyticsService: as,
		logger:           l,
	}
}

func (h *AnalyticsHTTPHandler) RegisterRoutes(r chi.Router) {
	r.Route("/accounts/{account_id}/analytics", func(r chi.Router) {
		r.Get("/summary", h.Summary)
		r.Get("/categories", h.CategoriesBreakdown)
		r.Get("/trends", h.Trends)
		r.Get("/members", h.MembersContribution)
	})
}
