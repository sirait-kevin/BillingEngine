package restful

import (
	"context"

	"github.com/sirait-kevin/BillingEngine/entities"
)

type UserUseCase interface {
	GetUserByID(ctx context.Context, id int64) (*entities.User, error)
	CreateUser(ctx context.Context, user *entities.User) (int64, error)
	UpdateUser(ctx context.Context, user *entities.User) error
}

type UserHandler struct {
	UserUseCase UserUseCase
}
