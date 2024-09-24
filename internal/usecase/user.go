package usecase

import (
	"context"
	"errors"
	"fmt"
	"github.com/bimaputraas/rest-api/internal/model"
	pkgerrors "github.com/bimaputraas/rest-api/pkg/errors"
	pkghash "github.com/bimaputraas/rest-api/pkg/hash"
	pkgjwt "github.com/bimaputraas/rest-api/pkg/jwt"
	pkgvalidate "github.com/bimaputraas/rest-api/pkg/validate"
	"github.com/golang-jwt/jwt/v5"
	"net/http"
	"time"
)

type (
	Login struct {
		PhoneNumber string `json:"phone_number" validate:"required,e164"`
		Pin         string `json:"pin" validate:"required"`
	}

	LoginResult struct {
		Token string `json:"token"`
	}
)

func (u *Usecase) Register(ctx context.Context, user model.User) (model.User, error) {
	if err := pkgvalidate.Struct(&user); err != nil {
		return model.User{}, pkgerrors.InvalidArgument(err)
	}

	uCheck, err := u.repo.GetUserByPhone(ctx, user.PhoneNumber)
	if err != nil {
		return model.User{}, err
	}

	if uCheck.PhoneNumber == user.PhoneNumber {
		return model.User{}, pkgerrors.InvalidArgument(fmt.Errorf("phone number already registered"))
	}

	if len(user.Pin) != 6 {
		return model.User{}, pkgerrors.InvalidArgument(fmt.Errorf("pin must be 6 digits"))
	}

	if !pkgvalidate.IsNumeric(user.Pin) {
		return model.User{}, pkgerrors.InvalidArgument(fmt.Errorf("pin must be a numbers"))
	}

	hashed, err := pkghash.FromString(user.Pin)
	if err != nil {
		return model.User{}, pkgerrors.InvalidArgument(err)
	}

	user.Pin = hashed
	user.Created = time.Now().Format(time.DateTime)
	txRepo, err := u.repo.BeginTx()
	if err != nil {
		return model.User{}, err
	}

	data, err := txRepo.InsertUser(ctx, user)
	if err != nil {
		errRB := txRepo.Rollback()
		if errRB != nil {
			return model.User{}, errRB
		}
		return model.User{}, err
	}

	err = txRepo.InsertBalance(ctx, model.Balance{
		UserID: user.ID,
	})

	if err != nil {
		errRB := txRepo.Rollback()
		if errRB != nil {
			return model.User{}, errRB
		}
		return model.User{}, err
	}

	data.Pin = ""
	return data, txRepo.Commit()
}

// Login generate token
func (u *Usecase) Login(ctx context.Context, login Login) (LoginResult, error) {
	var (
		plain  string
		hashed string
	)
	if err := pkgvalidate.Struct(&login); err != nil {
		return LoginResult{}, pkgerrors.InvalidArgument(err)
	}
	if len(login.Pin) != 6 {
		return LoginResult{}, pkgerrors.InvalidArgument(fmt.Errorf("pin must be 6 digits"))
	}

	if !pkgvalidate.IsNumeric(login.Pin) {
		return LoginResult{}, pkgerrors.InvalidArgument(fmt.Errorf("pin must be a numbers"))
	}

	plain = login.Pin

	user, err := u.repo.GetUserByPhone(ctx, login.PhoneNumber)
	if err != nil {
		code, _ := pkgerrors.Code(err)
		if code == http.StatusNotFound {
			return LoginResult{}, pkgerrors.InvalidArgument(fmt.Errorf("invalid phone or pin"))
		}
		return LoginResult{}, err
	}

	hashed = user.Pin

	if !pkghash.Check(plain, hashed) {
		return LoginResult{}, pkgerrors.InvalidArgument(fmt.Errorf("invalid phone or pin"))
	}

	secret := []byte(u.config.JWTSecret)
	token, err := pkgjwt.GenerateJWT(jwt.MapClaims{"user_id": user.ID}, secret)
	if err != nil {
		return LoginResult{}, pkgerrors.Internal(err)
	}

	return LoginResult{token}, nil
}

func (usecase *Usecase) Auth(ctx context.Context, token string) (uint, error) {
	if token == "" {
		return 0, pkgerrors.Illegal(errors.New("unauthenticated"))
	}

	secret := []byte(usecase.config.JWTSecret)
	claims, err := pkgjwt.ParseJWT(token, secret)
	if err != nil {
		return 0, pkgerrors.Illegal(errors.New("unauthenticated"))
	}

	userId, ok := claims["user_id"].(float64)
	if !ok {
		return 0, pkgerrors.Internal(errors.New("failed assert"))
	}

	user, err := usecase.repo.GetUserById(ctx, uint(userId))
	if err != nil {
		code, _ := pkgerrors.Code(err)
		if code == http.StatusNotFound {
			return 0, pkgerrors.Illegal(fmt.Errorf("unauthenticated"))
		}
		return 0, pkgerrors.Internal(err)
	}
	return user.ID, nil
}
