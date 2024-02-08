package app

import (
	"context"
	"dalkak/pkg/dtos"
	"dalkak/pkg/utils/httputils"
	"dalkak/pkg/utils/securityutils"
	"errors"
	"fmt"
	"mime"
	"net/http"
	"strings"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
)

func (app *APP) enableCORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", app.Origin)
		w.Header().Set("Access-Control-Allow-Credentials", "true")
		if r.Method == "OPTIONS" {
			w.Header().Set("Access-Control-Allow-Methods", "GET,POST,PUT,PATCH,DELETE,OPTIONS")
			w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, X-CSRF-Token, Authorization, x-client-id")
			return
		}
		next.ServeHTTP(w, r)
	})
}

func (app *APP) getTokenFromHeader(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			next.ServeHTTP(w, r)
			return
		}

		headerParts := strings.Split(authHeader, " ")
		if len(headerParts) != 2 || headerParts[0] != "Bearer" {
			httputils.ErrorJSON(w, errors.New("invalid auth header"), http.StatusBadRequest)
			return
		}

		token := headerParts[1]
		sub, err := securityutils.ParseTokenWithPublicKey(token, app.KmsSet)
		if err != nil {
			httputils.ErrorJSON(w, err, http.StatusUnauthorized)
			return
		}

		userInfo := dtos.UserInfo{WalletAddress: sub}

		ctx := context.WithValue(r.Context(), "user", userInfo)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func (app *APP) processData(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var reqMap map[string]interface{}

		contentType := r.Header.Get("Content-Type")
		if contentType == "" {
			next.ServeHTTP(w, r)
			return
		}

		mediaType, _, err := mime.ParseMediaType(contentType)
		if err != nil {
			httputils.ErrorJSON(w, err, http.StatusBadRequest)
			return
		}

		switch mediaType {
		case "application/json":
			if err := httputils.ReadJSON(w, r, &reqMap); err != nil {
				httputils.ErrorJSON(w, errors.New("JSON decoding error"), http.StatusBadRequest)
				return
			}

		// case "application/x-www-form-urlencoded", "multipart/form-data":
		// 	if err := r.ParseMultipartForm(config.MaxUploadSize); err != nil {
		// 		httputils.ErrorJSON(w, errors.New("Form parsing error"), http.StatusBadRequest)
		// 		return
		// 	}
		// 	reqMap = make(map[string]interface{})
		// 	for key, values := range r.Form {
		// 		if len(values) > 0 {
		// 			reqMap[key] = values[0]
		// 		}
		// 	}

		default:
			httputils.ErrorJSON(w, errors.New("Unsupported content type"), http.StatusUnsupportedMediaType)
			return
		}

		ctx := context.WithValue(r.Context(), "request", reqMap)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func (app *APP) verifyMetaMaskSignature(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		msg := "안전하게 지갑 연결"

		signature := common.FromHex(r.Context().Value("request").(map[string]interface{})["Signature"].(string))
		reqWalletAddr := common.HexToAddress(r.Context().Value("request").(map[string]interface{})["WalletAddress"].(string))

		if signature[64] != 27 && signature[64] != 28 {
			httputils.ErrorJSON(w, errors.New("Invalid MetaMask signature: incorrect recovery id"), http.StatusUnauthorized)
			return
		}
		signature[64] -= 27

		recoveredAddr, err := recoverAddressFromSignature(signature, []byte(msg))
		if err != nil {
			httputils.ErrorJSON(w, errors.New("Invalid MetaMask signature: "+err.Error()), http.StatusUnauthorized)
			return
		}

		if recoveredAddr.Hex() == reqWalletAddr.Hex() {
			next.ServeHTTP(w, r)
		} else {
			httputils.ErrorJSON(w, errors.New("Invalid MetaMask signature: address mismatch"), http.StatusUnauthorized)
			return
		}
	})
}

func signHash(data []byte) []byte {
	prefix := fmt.Sprintf("\x19Ethereum Signed Message:\n%d", len(data))
	return crypto.Keccak256([]byte(prefix), data)
}

func recoverAddressFromSignature(signature []byte, data []byte) (common.Address, error) {
	publicKey, err := crypto.SigToPub(signHash(data), signature)
	if err != nil {
		return common.Address{}, err
	}
	return crypto.PubkeyToAddress(*publicKey), nil
}
