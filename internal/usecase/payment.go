package usecase

import (
	"context"
	"fmt"
	"time"

	"github.com/bimaputraas/rest-api/internal/model"
	pkgerrors "github.com/bimaputraas/rest-api/pkg/errors"
	pkgvalidate "github.com/bimaputraas/rest-api/pkg/validate"
)

type (
	Payment struct {
		Amount  float64 `validate:"required" json:"amount"`
		Remarks string  `validate:"required" json:"remarks"`
	}
)

func (u *Usecase) Payment(ctx context.Context, userId uint, payment Payment) (model.Payment, error) {
	if err := pkgvalidate.Struct(payment); err != nil {
		return model.Payment{}, pkgerrors.InvalidArgument(err)
	}
	remarks := payment.Remarks
	amount := payment.Amount
	if amount < 1 {
		return model.Payment{}, pkgerrors.InvalidArgument(fmt.Errorf("invalid amount"))
	}

	userBalance, err := u.repo.Db.GetBalanceByUId(ctx, userId)
	if err != nil {
		return model.Payment{}, err
	}

	userBalanceBefore := userBalance.CurrentBalance
	userBalanceAfter := userBalanceBefore - amount
	now := time.Now().Format(time.DateTime)

	if userBalanceAfter < 0 {
		return model.Payment{}, pkgerrors.InvalidArgument(fmt.Errorf("balance is not enough"))
	}

	txRepo, err := u.repo.Db.BeginTx()
	if err != nil {
		return model.Payment{}, err
	}

	userBalance.CurrentBalance = userBalanceAfter
	userBalance.Updated = now
	err = txRepo.UpdateBalance(ctx, userBalance)
	if err != nil {
		errRB := txRepo.Rollback()
		if errRB != nil {
			return model.Payment{}, errRB
		}
		return model.Payment{}, err
	}

	data, err := txRepo.InsertPayment(ctx, model.Payment{
		UserID:        userId,
		Amount:        amount,
		Remarks:       remarks,
		BalanceBefore: userBalanceBefore,
		BalanceAfter:  userBalanceAfter,
		Created:       now,
	})

	if err != nil {
		errRB := txRepo.Rollback()
		if errRB != nil {
			return model.Payment{}, errRB
		}
		return model.Payment{}, err
	}

	transaction := model.Transaction{
		UserID:        userId,
		PaymentId:     &data.ID,
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
			return model.Payment{}, errRB
		}
		return model.Payment{}, err
	}

	data.UserID = 0
	return data, txRepo.Commit()
}
