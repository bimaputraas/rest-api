package usecase

import (
	"context"

	"github.com/bimaputraas/rest-api/internal/model"
)

func (u *Usecase) GetAllUserTransactions(ctx context.Context, userId uint) ([]model.Transaction, error) {
	data, err := u.repo.Db.GetTransactionsByUId(ctx, userId)
	if err != nil {
		return nil, err
	}

	return data, nil
}
