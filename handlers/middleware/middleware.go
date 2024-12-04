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
		clientKey := r.Header.Get("Client-Key")
		signature := r.Header.Get("X-Signature")
		if clientKey == "" || signature == "" {
			helper.JSON(w, http.StatusUnauthorized, "Missing Client-Key or signature", nil)
			return
		}

		secret, ok := clientSecrets[clientKey]
		if !ok {
			helper.JSON(w, http.StatusUnauthorized, "Invalid Client-Key", nil)
			return
		}

		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			helper.JSON(w, http.StatusInternalServerError, "Error reading request body", nil)
			return
		}
		r.Body = ioutil.NopCloser(bytes.NewBuffer(body))

		hash := hmac.New(sha256.New, []byte(secret))
		hash.Write(body)
		expectedSignature := hex.EncodeToString(hash.Sum(nil))

		if !hmac.Equal([]byte(expectedSignature), []byte(signature)) {
			if os.Getenv("DEBUG_MODE") == "true" {
				helper.JSON(w, http.StatusUnauthorized, "Invalid signature", map[string]interface{}{
					"expected_signature": expectedSignature,
				})
			} else {
				helper.JSON(w, http.StatusUnauthorized, "Invalid signature", nil)
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
			if err := recover(); err != nil {
				log.Println("Recovered from panic:", err)
				helper.JSON(w, http.StatusInternalServerError, "Internal Server Error", nil)
			}
		}()
		next.ServeHTTP(w, r)
	})
}
