package usecases

import (
	"context"

	"github.com/sirait-kevin/BillingEngine/entities"
)

type UserRepository interface {
	GetByID(ctx context.Context, id int64) (*entities.User, error)
	Create(ctx context.Context, user *entities.User) (int64, error)
	Update(ctx context.Context, user *entities.User) error
}

type UserUseCase struct {
	UserRepository UserRepository
}
