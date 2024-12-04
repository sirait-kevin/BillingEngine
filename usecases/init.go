package usecases

import "github.com/sirait-kevin/BillingEngine/entities"

type UserRepository interface {
	GetByID(id int64) (*entities.User, error)
	Create(user *entities.User) (int64, error)
	Update(user *entities.User) error
}

type UserUseCase struct {
	UserRepository UserRepository
}
