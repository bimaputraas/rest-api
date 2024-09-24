package usecase

import (
	"context"
	"fmt"
	"github.com/bimaputraas/rest-api/internal/model"
	pkgerrors "github.com/bimaputraas/rest-api/pkg/errors"
	"time"
)

func (u *Usecase) TopUp(ctx context.Context, userId uint, amountTopUp float64) (model.TopUp, error) {
	if amountTopUp < 1 {
		return model.TopUp{}, pkgerrors.InvalidArgument(fmt.Errorf("invalid amount"))
	}

	uBalance, err := u.repo.GetBalanceByUId(ctx, userId)
	if err != nil {
		return model.TopUp{}, err
	}

	balanceBefore := uBalance.CurrentBalance
	balanceAfter := balanceBefore + amountTopUp
	now := time.Now().Format("2006-01-02 15:04:05")

	txRepo, err := u.repo.BeginTx()
	if err != nil {
		return model.TopUp{}, err
	}

	uBalance.CurrentBalance = balanceAfter
	uBalance.Updated = now
	err = txRepo.UpdateBalance(ctx, uBalance)
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
		BalanceBefore: balanceBefore,
		BalanceAfter:  balanceAfter,
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
		TopUpId:       data.ID,
		Amount:        data.AmountTopUp,
		Remarks:       "",
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
