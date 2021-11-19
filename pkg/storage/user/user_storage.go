package user

import (
	"context"
	"database/sql"
	"kanggo/pkg/entity/schema"

	"gorm.io/gorm"
)

//go:generate mockery --name UserStorage --case snake --output ../../mocks --disable-version-string

type (
	UserStorage interface {
		Insert(ctx context.Context, data schema.User) error
		GetByEmail(ctx context.Context, email string) (*schema.User, error)
	}

	userStorage struct {
		Native *sql.DB
		Gorm   *gorm.DB
	}
)

func NewUserStorage(native *sql.DB, gorm *gorm.DB) UserStorage {
	return &userStorage{
		Native: native,
		Gorm:   gorm,
	}
}

func (m *userStorage) Insert(ctx context.Context, data schema.User) error {
	if err := m.Gorm.WithContext(ctx).Create(&data).Error; err != nil {
		return err
	}

	return nil
}

func (m *userStorage) GetByEmail(ctx context.Context, email string) (*schema.User, error) {

	user := schema.User{}
	qry := `SELECT id, name, email, password FROM users WHERE email = ? `
	res := m.Native.QueryRowContext(ctx, qry, email)
	if err := res.Scan(&user.Base.Id, &user.Name, &user.Email, &user.Password); err != nil {
		if err == sql.ErrNoRows {
			return nil, err
		}
		return nil, err
	}

	return &user, nil
}
