package usecase

import (
	"context"
	"fmt"
	"time"

	"github.com/bimaputraas/rest-api/internal/model"
	pkgerrors "github.com/bimaputraas/rest-api/pkg/errors"
)

func (u *Usecase) TopUp(ctx context.Context, userId uint, amountTopUp float64) (model.TopUp, error) {
	if amountTopUp < 1 {
		return model.TopUp{}, pkgerrors.InvalidArgument(fmt.Errorf("invalid amount"))
	}

	userBalance, err := u.repo.Storage.GetBalanceByUId(ctx, userId)
	if err != nil {
		return model.TopUp{}, err
	}

	userBalanceBefore := userBalance.CurrentBalance
	userBalanceAfter := userBalanceBefore + amountTopUp
	now := time.Now().Format("2006-01-02 15:04:05")

	txRepo, err := u.repo.Storage.BeginTx()
	if err != nil {
		return model.TopUp{}, err
	}

	userBalance.CurrentBalance = userBalanceAfter
	userBalance.Updated = now
	err = txRepo.UpdateBalance(ctx, userBalance)
	if err != nil {
		errRB := txRepo.Rollback()
		if errRB != nil {
			return model.TopUp{}, errRB
		}
		return model.TopUp{}, err
	}

	data, err := txRepo.InsertTopUp(ctx, model.TopUp{
		UserID:        userId,
		AmountTopUp:   amountTopUp,
		BalanceBefore: userBalanceBefore,
		BalanceAfter:  userBalanceAfter,
		Created:       now,
	})

	if err != nil {
		errRB := txRepo.Rollback()
		if errRB != nil {
			return model.TopUp{}, errRB
		}
		return model.TopUp{}, err
	}

	transaction := model.Transaction{
		UserID:        userId,
		TopUpId:       &data.ID,
		Amount:        data.AmountTopUp,
		Remarks:       "Top Up",
		BalanceBefore: data.BalanceBefore,
		BalanceAfter:  data.BalanceAfter,
		Status:        "SUCCESS",
		Created:       now,
	}

	transaction.MockRandType()
	_, err = txRepo.InsertTransaction(ctx, transaction)
	if err != nil {
		errRB := txRepo.Rollback()
		if errRB != nil {
			return model.TopUp{}, errRB
		}
		return model.TopUp{}, err
	}

	data.UserID = 0
	return data, txRepo.Commit()
}
