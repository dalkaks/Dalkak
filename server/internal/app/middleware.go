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
		var reqMap map[string]interface{}
		contentType := r.Header.Get("Content-Type")

		switch {
		case strings.Contains(contentType, "application/json"):
			reqMap = make(map[string]interface{})
			if err := json.NewDecoder(r.Body).Decode(&reqMap); err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}

		case strings.Contains(contentType, "application/x-www-form-urlencoded"),
			strings.Contains(contentType, "multipart/form-data"):
			const maxMemory = 32 << 20 // 32MB
			if err := r.ParseMultipartForm(maxMemory); err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}
			reqMap = make(map[string]interface{})
			for key, values := range r.Form {
				if len(values) > 0 {
					reqMap[key] = values[0]
				}
			}

		default:
			http.Error(w, "Unsupported content type", http.StatusUnsupportedMediaType)
			return
		}

		ctx := context.WithValue(r.Context(), "request", reqMap)
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
