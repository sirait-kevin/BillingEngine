package restful

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/sirait-kevin/BillingEngine/entities"
	"github.com/sirait-kevin/BillingEngine/pkg/errs"
	"github.com/sirait-kevin/BillingEngine/pkg/helper"
)

func (h *BillingHandler) GetUserByID(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	id, err := strconv.ParseInt(r.FormValue("user_id"), 10, 64)
	if err != nil {
		helper.JSON(w, ctx, nil, errs.NewWithMessage(http.StatusBadRequest, "Invalid user ID"))
		return
	}

	helper.JSON(w, ctx, id, nil)
}

func (h *BillingHandler) CreateUser(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	var user entities.User
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		helper.JSON(w, ctx, nil, errs.NewWithMessage(http.StatusBadRequest, "Invalid request payload"))
		return
	}

	helper.JSON(w, ctx, map[string]int64{
		"user_id": 1,
	}, nil)

}
