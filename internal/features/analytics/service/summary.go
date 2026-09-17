package service

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/loundxr/expense-tracker/internal/core/domain"
	core_errors "github.com/loundxr/expense-tracker/internal/core/errors"
	core_ctx "github.com/loundxr/expense-tracker/internal/core/transport/http/context"
)

func (s *AnalyticsService) Summary(
	ctx context.Context,
	filter domain.SummaryFilter,
) (domain.Summary, error) {
	const op = "analytics.service.Summary"

	uid := core_ctx.GetUserID(ctx)
	role := core_ctx.GetUserRole(ctx)

	if role != domain.RoleAdmin {
		hasAccess, err := s.accChecker.HasAccess(ctx, uid, filter.AccountID)
		if err != nil {
			return domain.Summary{}, fmt.Errorf("%s: %w", op, err)
		}
		if !hasAccess {
			return domain.Summary{}, fmt.Errorf("%s: %w", op, core_errors.ErrForbidden)
		}
	}

	summary, err := s.analyticsRepository.Summary(ctx, filter)
	if err != nil {
		return domain.Summary{}, fmt.Errorf("%s: %w", op, err)
	}

	s.logger.Info(
		"user generated analytics summary",
		slog.Int64("user_id", uid),
		slog.Int64("account_id", filter.AccountID),
	)
	return summary, nil
}
