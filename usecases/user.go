package usecases

import "github.com/sirait-kevin/BillingEngine/entities"

func (u *UserUseCase) GetUserByID(id int64) (*entities.User, error) {
	return u.UserRepository.GetByID(id)
}

func (u *UserUseCase) CreateUser(user *entities.User) (int64, error) {
	return u.UserRepository.Create(user)
}

func (u *UserUseCase) UpdateUser(user *entities.User) error {
	return u.UserRepository.Update(user)
}
