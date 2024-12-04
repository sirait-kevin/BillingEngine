package middleware

import (
	"bytes"
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/sirupsen/logrus"

	"github.com/sirait-kevin/BillingEngine/pkg/errs"
	"github.com/sirait-kevin/BillingEngine/pkg/helper"
	"github.com/sirait-kevin/BillingEngine/pkg/logger"
)

// Mocked client keys and secrets for simplicity
var clientSecrets = map[string]string{
	"client1": "secret1",
	"client2": "secret2",
	// Add more clients and their secrets here
}

func VerifySignatureMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		clientKey := r.Header.Get("Client-Key")
		signature := r.Header.Get("X-Signature")
		if clientKey == "" || signature == "" {
			helper.JSON(w, ctx, nil, errs.NewWithMessage(http.StatusBadRequest, "Missing Client-Key or signature"))
			return
		}

		secret, ok := clientSecrets[clientKey]
		if !ok {
			helper.JSON(w, ctx, nil, errs.NewWithMessage(http.StatusUnauthorized, "Invalid Client-Key"))
			return
		}

		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			helper.JSON(w, ctx, nil, errs.NewWithMessage(http.StatusInternalServerError, "Error reading request body"))
			return
		}
		r.Body = ioutil.NopCloser(bytes.NewBuffer(body))

		hash := hmac.New(sha256.New, []byte(secret))
		hash.Write(body)
		expectedSignature := hex.EncodeToString(hash.Sum(nil))

		if !hmac.Equal([]byte(expectedSignature), []byte(signature)) {
			if os.Getenv("DEBUG_MODE") == "true" {
				helper.JSON(w, ctx, map[string]interface{}{
					"expected_signature": expectedSignature,
				}, errs.NewWithMessage(http.StatusUnauthorized, "Invalid signature"))
			} else {
				helper.JSON(w, ctx, nil, errs.NewWithMessage(http.StatusUnauthorized, "Invalid signature"))
			}
			return
		}

		next.ServeHTTP(w, r)
	})
}

func LoggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		logger := logger.Log.WithFields(logrus.Fields{
			"method": r.Method,
			"url":    r.URL.String(),
		})

		ctx := context.WithValue(r.Context(), "logger", logger)
		rw := &responseWriter{w, http.StatusOK}
		next.ServeHTTP(rw, r.WithContext(ctx))

		duration := time.Since(start)
		logger.WithFields(logrus.Fields{
			"status":   rw.statusCode,
			"duration": duration,
		}).Info("request completed")
	})
}

type responseWriter struct {
	http.ResponseWriter
	statusCode int
}

func (rw *responseWriter) WriteHeader(code int) {
	rw.statusCode = code
	rw.ResponseWriter.WriteHeader(code)
}

func ErrorHandlingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			ctx := r.Context()
			if err := recover(); err != nil {
				log.Println("Recovered from panic:", err)
				helper.JSON(w, ctx, nil, errs.NewWithMessage(http.StatusInternalServerError, "Internal Server Error"))
			}
		}()
		next.ServeHTTP(w, r)
	})
}
