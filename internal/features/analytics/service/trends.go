package service

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/loundxr/expense-tracker/internal/core/domain"
	core_errors "github.com/loundxr/expense-tracker/internal/core/errors"
	core_ctx "github.com/loundxr/expense-tracker/internal/core/transport/http/context"
)

func (s *AnalyticsService) Trends(ctx context.Context, f domain.TrendsFilter) ([]domain.ExpenseTrendPoint, error) {
	const op = "analytics.service.Trends"

	uid := core_ctx.GetUserID(ctx)
	role := core_ctx.GetUserRole(ctx)

	if role != domain.RoleAdmin {
		hasAccess, err := s.accChecker.HasAccess(ctx, uid, f.AccountID)
		if err != nil {
			return nil, fmt.Errorf("%s: %w", op, err)
		}
		if !hasAccess {
			return nil, fmt.Errorf("%s: %w", op, core_errors.ErrForbidden)
		}
	}

	points, err := s.analyticsRepository.Trends(ctx, f)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	s.logger.Info(
		"user generated trends analytics",
		slog.Int64("user_id", uid),
	)

	return points, nil
}
