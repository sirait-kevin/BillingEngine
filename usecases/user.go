package usecases

import (
	"context"

	"github.com/sirait-kevin/BillingEngine/entities"
)

func (u *UserUseCase) GetUserByID(ctx context.Context, id int64) (*entities.User, error) {
	return u.DBRepo.GetByID(ctx, id)
}

func (u *UserUseCase) CreateUser(ctx context.Context, user *entities.User) (int64, error) {
	return u.DBRepo.Create(ctx, user)
}

func (u *UserUseCase) UpdateUser(ctx context.Context, user *entities.User) error {
	return u.DBRepo.Update(ctx, user)
}
