package mq

import (
	"context"
	"encoding/json"
	"log"

	"github.com/nsqio/go-nsq"

	"github.com/sirait-kevin/BillingEngine/entities"
	"github.com/sirait-kevin/BillingEngine/usecases"
)

type NSQHandler struct {
	UserUseCase *usecases.BillingUseCase
}

func (h *NSQHandler) HandleMessage(message *nsq.Message) error {
	ctx := context.Background()
	var user entities.User
	err := json.Unmarshal(message.Body, &user)
	if err != nil {
		log.Printf("Error unmarshalling message: %v", err)
		return err
	}

	err = h.UserUseCase.UpdateUser(ctx, &user)
	if err != nil {
		log.Printf("Error updating user: %v", err)
		return err
	}

	log.Printf("User updated successfully: %v", user)
	return nil
}
