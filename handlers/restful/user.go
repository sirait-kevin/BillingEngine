package restful

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"

	"github.com/sirait-kevin/BillingEngine/entities"
	"github.com/sirait-kevin/BillingEngine/pkg/errs"
	"github.com/sirait-kevin/BillingEngine/pkg/helper"
)

func (h *UserHandler) GetUserByID(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	vars := mux.Vars(r)
	id, err := strconv.ParseInt(vars["id"], 10, 64)
	if err != nil {
		helper.JSON(w, ctx, nil, errs.NewWithMessage(http.StatusBadRequest, "Invalid user ID"))
		return
	}

	user, err := h.BillingUC.GetUserByID(ctx, id)
	if err != nil {
		helper.JSON(w, ctx, nil, err)
		return
	}

	helper.JSON(w, ctx, user, nil)
}

func (h *UserHandler) CreateUser(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	var user entities.User
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		helper.JSON(w, ctx, nil, errs.NewWithMessage(http.StatusBadRequest, "Invalid request payload"))
		return
	}

	id, err := h.BillingUC.CreateUser(ctx, &user)
	if err != nil {
		helper.JSON(w, ctx, nil, err)
		return
	}

	helper.JSON(w, ctx, map[string]int64{
		"user_id": id,
	}, nil)

}
