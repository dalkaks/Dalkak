package app

import (
	"context"
	"encoding/json"
	"net/http"
	"strings"
)

func (app *APP) enableCORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", app.Origin)
		if r.Method == "OPTIONS" {
			w.Header().Set("Access-Control-Allow-Credentials", "true")
			w.Header().Set("Access-Control-Allow-Methods", "GET,POST,PUT,PATCH,DELETE,OPTIONS")
			w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, X-CSRF-Token, Authorization, x-client-id")
			return
		}
		next.ServeHTTP(w, r)
	})
}

func (app *APP) processData(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var req interface{}
		contentType := r.Header.Get("Content-Type")

		switch {
		case strings.Contains(contentType, "application/json"):
			if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}

		case strings.Contains(contentType, "application/x-www-form-urlencoded"),
			strings.Contains(contentType, "multipart/form-data"):
			if err := r.ParseForm(); err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}

		default:
			http.Error(w, "Unsupported content type", http.StatusUnsupportedMediaType)
			return
		}

		ctx := context.WithValue(r.Context(), "request", req)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func (app *APP) verifyMetaMaskSignature(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// 사인 데이터 추출 및 검증 로직

		// 사인 데이터가 유효하다면, 다음 핸들러로 요청을 전달
		if true /* isValidSignature(signature) */ {
			next.ServeHTTP(w, r)
		} else {
			http.Error(w, "Invalid MetaMask signature", http.StatusUnauthorized)
		}
	})
}
