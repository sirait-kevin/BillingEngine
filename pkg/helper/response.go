package helper

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/sirait-kevin/BillingEngine/pkg/errs"
)

type Response struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

func JSON(w http.ResponseWriter, ctx context.Context, data interface{}, err error) {
	var (
		responseCode = http.StatusOK
		responseMsg  = "SUCCESS"
	)

	if err != nil {
		responseCode = errs.GetHTTPCode(err)
		responseMsg = err.Error()
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(responseCode)
	json.NewEncoder(w).Encode(Response{Code: responseCode, Message: responseMsg, Data: data})
}
