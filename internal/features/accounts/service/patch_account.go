package service

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/loundxr/expense-tracker/internal/core/domain"
	core_errors "github.com/loundxr/expense-tracker/internal/core/errors"
	core_ctx "github.com/loundxr/expense-tracker/internal/core/transport/http/context"
)

func (s *AccountsService) PatchAccount(ctx context.Context, id int64, patch domain.AccountPatch) (domain.Account, error) {
	const op = "accounts.service.PatchAccount"
	uid := core_ctx.GetUserID(ctx)
	role := core_ctx.GetUserRole(ctx)

	if role != domain.RoleAdmin && uid != id {
		return domain.Account{}, fmt.Errorf("%s: %w", op, core_errors.ErrForbidden)
	}

	account, err := s.accountsRepository.GetAccountByID(ctx, id)
	if err != nil {
		return domain.Account{}, fmt.Errorf("%s: %w", op, err)
	}

	var patchedFields []string

	if patch.Name.Set {
		patchedFields = append(patchedFields, "name")
	}

	if err := account.ApplyPatch(patch); err != nil {
		return domain.Account{}, fmt.Errorf("%s: apply patch: %w", op, err)
	}

	patchedAccount, err := s.accountsRepository.PatchAccount(ctx, id, account)
	if err != nil {
		return domain.Account{}, fmt.Errorf("%s: %w", op, err)
	}

	s.logger.Info(
		"account patched",
		slog.Int64("actor_id", uid),
		slog.Int64("account_id", patchedAccount.ID),
		slog.Any("patched_fields", patchedFields),
	)
	return patchedAccount, nil
}
