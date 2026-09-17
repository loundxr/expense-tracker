package service

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/loundxr/expense-tracker/internal/core/domain"
	core_errors "github.com/loundxr/expense-tracker/internal/core/errors"
	core_ctx "github.com/loundxr/expense-tracker/internal/core/transport/http/context"
)

func (s *AnalyticsService) CategoriesBreakdown(ctx context.Context, f domain.SummaryFilter) ([]domain.CategoryBreakdown, error) {
	const op = "analytics.service.CategoriesBreakdown"
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

	breakdown, err := s.analyticsRepository.CategoriesBreakdown(ctx, f)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	s.logger.Info(
		"user generated categories breakdown",
		slog.Int64("user_id", uid),
		slog.Int64("account_id", f.AccountID),
	)
	return breakdown, nil
}
