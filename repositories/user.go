package repositories

import (
	"context"
	"database/sql"
	"fmt"
	"net/http"

	"github.com/sirupsen/logrus"

	"github.com/sirait-kevin/BillingEngine/entities"
	"github.com/sirait-kevin/BillingEngine/pkg/errs"
	"github.com/sirait-kevin/BillingEngine/pkg/helper"
)

type UserDB struct {
	ID             int64  `json:"id"`
	EncryptedName  []byte `json:"encrypted_name"`
	EncryptedEmail []byte `json:"encrypted_email"`
	HashedPassword string `json:"hashed_password"`
}

func (r *DBRepository) GetByID(ctx context.Context, id int64) (*entities.User, error) {
	fmt.Println("masuk sini")
	logger := ctx.Value("logger").(*logrus.Entry)
	logger.Debug("Fetching user from database by ID: ", id)

	user := &entities.User{}
	userDB := &UserDB{}

	err := r.DB.QueryRowContext(ctx, "SELECT id, encrypted_name, encrypted_email, hashed_password FROM users WHERE id = ?", id).Scan(&user.ID, &userDB.EncryptedName, &userDB.EncryptedEmail, &userDB.HashedPassword)
	if err != nil {
		if err == sql.ErrNoRows {
			err = errs.Wrap(http.StatusNotFound, err)
		}
		logger.Error("Error fetching user: ", err)
		return nil, err
	}
	user.Name, err = helper.Decrypt(userDB.EncryptedName)
	if err != nil {
		logger.Error("Error decrypting user name: ", err)
		return nil, err
	}
	user.Email, err = helper.Decrypt(userDB.EncryptedEmail)
	if err != nil {
		logger.Error("Error decrypting user email: ", err)
		return nil, err
	}
	return user, nil
}

func (r *DBRepository) Create(ctx context.Context, user *entities.User) (int64, error) {
	logger := ctx.Value("logger").(*logrus.Entry)
	logger.Debug("Inserting user into database: ", user)
	var (
		err    error
		userDB = &UserDB{}
	)
	userDB.EncryptedName, err = helper.Encrypt(user.Name)
	if err != nil {
		logger.Error("Error encrypting user name: ", err)
		return 0, err
	}
	userDB.EncryptedEmail, err = helper.Encrypt(user.Email)
	if err != nil {
		logger.Error("Error encrypting user email: ", err)
		return 0, err
	}
	userDB.HashedPassword, err = helper.HashPassword(user.Password)
	if err != nil {
		logger.Error("Error hashing user password: ", err)
		return 0, err
	}
	result, err := r.DB.ExecContext(ctx, "INSERT INTO users (encrypted_name, encrypted_email, hashed_password) VALUES (?, ?, ?)", userDB.EncryptedName, userDB.EncryptedEmail, userDB.HashedPassword)
	if err != nil {
		logger.Error("Error creating user: ", err)
		return 0, err
	}
	id, err := result.LastInsertId()
	if err != nil {
		logger.Error("Error getting last insert ID: ", err)
		return 0, err
	}
	return id, nil
}

func (r *DBRepository) Update(ctx context.Context, user *entities.User) error {
	logger := ctx.Value("logger").(*logrus.Logger)
	logger.Debug("Updating user in database: ", user)
	var (
		err    error
		userDB = &UserDB{}
	)

	userDB.EncryptedName, err = helper.Encrypt(user.Name)
	if err != nil {
		logger.Error("Error encrypting user name: ", err)
		return err
	}
	userDB.EncryptedEmail, err = helper.Encrypt(user.Email)
	if err != nil {
		logger.Error("Error encrypting user email: ", err)
		return err
	}
	userDB.HashedPassword, err = helper.HashPassword(user.Password)
	if err != nil {
		logger.Error("Error hashing user password: ", err)
		return err
	}
	_, err = r.DB.ExecContext(ctx, "UPDATE users SET encrypted_name = ?, encrypted_email = ?, hashed_password = ? WHERE id = ?", userDB.EncryptedName, userDB.EncryptedEmail, userDB.HashedPassword, user.ID)
	if err != nil {
		logger.Error("Error updating user: ", err)
		return err
	}
	return nil
}
