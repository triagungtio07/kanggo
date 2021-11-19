package user

import (
	"context"
	"kanggo/pkg/entity/model"
	"kanggo/pkg/entity/schema"
	storage "kanggo/pkg/storage/user"
	"kanggo/utils"

	"github.com/jinzhu/copier"
)

//go:generate mockery --name UserUsecase --case snake --output ../../mocks --disable-version-string

type (
	UserUsecase interface {
		Insert(context.Context, model.RegisterRequest) error
		GetByEmail(context.Context, string) (*model.UserResponse, error)
	}

	userUsecase struct {
		userStorage storage.UserStorage
	}
)

func NewUserUsecase(userStorage storage.UserStorage) UserUsecase {
	return &userUsecase{
		userStorage: userStorage,
	}
}

func (u *userUsecase) Insert(ctx context.Context, req model.RegisterRequest) error {
	hashedPassword := utils.HashPassword(req.Password)
	req.Password = hashedPassword
	data := schema.User{}

	copier.Copy(&data, &req)

	if err := u.userStorage.Insert(ctx, data); err != nil {
		return err
	}

	return nil
}

func (u *userUsecase) GetByEmail(ctx context.Context, email string) (*model.UserResponse, error) {

	res, err := u.userStorage.GetByEmail(ctx, email)
	if err != nil {
		return nil, err
	}

	user := model.UserResponse{
		Id:       int(res.Id),
		Name:     res.Name,
		Email:    res.Email,
		Password: res.Password,
	}

	return &user, nil
}
