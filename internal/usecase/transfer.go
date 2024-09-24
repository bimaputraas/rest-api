package usecase

import (
	"context"
	"fmt"
	"github.com/bimaputraas/rest-api/internal/model"
	pkgerrors "github.com/bimaputraas/rest-api/pkg/errors"
	pkgvalidate "github.com/bimaputraas/rest-api/pkg/validate"
	"time"
)

type (
	Transfer struct {
		TargetUser uint    `validate:"required" json:"target_user"`
		Amount     float64 `validate:"required" json:"amount"`
		Remarks    string  `validate:"required" json:"remarks"`
	}
)

func (u *Usecase) Transfer(ctx context.Context, userId uint, transfer Transfer) (model.Transfer, error) {
	if err := pkgvalidate.Struct(transfer); err != nil {
		return model.Transfer{}, pkgerrors.InvalidArgument(err)
	}

	remarks := transfer.Remarks
	amount := transfer.Amount
	targetUid := transfer.TargetUser
	now := time.Now().Format(time.DateTime)

	_, err := u.repo.GetUserById(ctx, targetUid)
	if pkgerrors.Code(err) == pkgerrors.ErrNotFound {
		return model.Transfer{}, pkgerrors.InvalidArgument(fmt.Errorf("target user is not exist"))
	}
	if err != nil {
		return model.Transfer{}, pkgerrors.Internal(err)
	}

	if amount < 1 {
		return model.Transfer{}, pkgerrors.InvalidArgument(fmt.Errorf("invalid amount"))
	}

	uBalance, err := u.repo.GetBalanceByUId(ctx, userId)
	if err != nil {
		return model.Transfer{}, err
	}

	balanceBefore := uBalance.CurrentBalance
	balanceAfter := balanceBefore - amount
	if balanceAfter < 0 {
		return model.Transfer{}, pkgerrors.InvalidArgument(fmt.Errorf("balance is not enough"))
	}

	uBalanceTarget, err := u.repo.GetBalanceByUId(ctx, targetUid)
	if err != nil {
		return model.Transfer{}, err
	}

	balanceBeforeTarget := uBalanceTarget.CurrentBalance
	balanceAfterTarget := balanceBeforeTarget + amount

	txRepo, err := u.repo.BeginTx()
	if err != nil {
		return model.Transfer{}, err
	}

	uBalance.CurrentBalance = balanceAfter
	uBalance.Updated = now
	err = txRepo.UpdateBalance(ctx, uBalance)
	if err != nil {
		errRB := txRepo.Rollback()
		if errRB != nil {
			return model.Transfer{}, errRB
		}
		return model.Transfer{}, err
	}

	uBalanceTarget.CurrentBalance = balanceAfterTarget
	uBalanceTarget.Updated = now
	err = txRepo.UpdateBalance(ctx, uBalanceTarget)
	if err != nil {
		errRB := txRepo.Rollback()
		if errRB != nil {
			return model.Transfer{}, errRB
		}
		return model.Transfer{}, err
	}

	data, err := txRepo.InsertTransfer(ctx, model.Transfer{
		UserID:        userId,
		Amount:        amount,
		Remarks:       remarks,
		BalanceBefore: balanceBefore,
		BalanceAfter:  balanceAfter,
		Created:       now,
	})
	if err != nil {
		errRB := txRepo.Rollback()
		if errRB != nil {
			return model.Transfer{}, errRB
		}
		return model.Transfer{}, err
	}

	transaction := model.Transaction{
		UserID:        userId,
		TransferId:    &data.ID,
		Amount:        data.Amount,
		Remarks:       data.Remarks,
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
			return model.Transfer{}, errRB
		}
		return model.Transfer{}, err
	}

	data.UserID = 0
	return data, txRepo.Commit()
}
