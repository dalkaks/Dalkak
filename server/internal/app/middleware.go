package app

import (
	"bytes"
	"context"
	appsecurity "dalkak/internal/security"
	"dalkak/pkg/dtos"
	"dalkak/pkg/utils/httputils"
	"encoding/json"
	"errors"
	"fmt"
	"io"
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
			httputils.ErrorJSON(w, errors.New("invalid auth header"), http.StatusUnauthorized)
			return
		}

		token := headerParts[1]
		sub, err := appsecurity.ParseTokenWithPublicKey(token, app.KmsSet)
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
		contentType := r.Header.Get("Content-Type")
		if contentType == "" {
			next.ServeHTTP(w, r)
			return
		}

		mediaType, _, err := mime.ParseMediaType(contentType)
		if err != nil {
			httputils.ErrorJSON(w, errors.New("Failed to parse media type"), http.StatusBadRequest)
			return
		}

		switch mediaType {
		case "application/json":
			next.ServeHTTP(w, r)

		default:
			httputils.ErrorJSON(w, errors.New("Unsupported content type"), http.StatusUnsupportedMediaType)
			return
		}
	})
}

func (app *APP) verifyMetaMaskSignature(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// 요청 본문을 버퍼에 읽어 저장
		bodyBytes, err := io.ReadAll(r.Body)
		if err != nil {
			httputils.ErrorJSON(w, errors.New("Failed to read request body"), http.StatusInternalServerError)
			return
		}

		// 요청 본문을 다시 설정하여 후속 처리에서도 사용할 수 있도록 함
		r.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))

		// 버퍼에서 읽은 본문 데이터를 사용하여 필요한 작업 수행
		var requestData map[string]interface{}
		if err := json.Unmarshal(bodyBytes, &requestData); err != nil {
			httputils.ErrorJSON(w, errors.New("Failed to parse request body"), http.StatusBadRequest)
			return
		}

		msg := "안전하게 지갑 연결"

		signatureStr, ok := requestData["signature"].(string)
		if !ok {
			httputils.ErrorJSON(w, errors.New("signature is missing or not a string"), http.StatusBadRequest)
			return
		}
		signature := common.FromHex(signatureStr)

		walletAddressStr, ok := requestData["walletAddress"].(string)
		if !ok {
			httputils.ErrorJSON(w, errors.New("walletAddress is missing or not a string"), http.StatusBadRequest)
			return
		}
		reqWalletAddr := common.HexToAddress(walletAddressStr)

		if signature[64] != 27 && signature[64] != 28 {
			httputils.ErrorJSON(w, errors.New("Invalid MetaMask signature: incorrect recovery id"), http.StatusUnauthorized)
			return
		}
		signature[64] -= 27

		recoveredAddr, err := recoverAddressFromSignature(signature, []byte(msg))
		if err != nil {
			httputils.ErrorJSON(w, errors.New("Invalid MetaMask signature: failed to recover address"), http.StatusUnauthorized)
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
