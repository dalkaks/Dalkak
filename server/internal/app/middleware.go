package app

import (
	"context"
	"encoding/json"
	"fmt"
	"mime"
	"net/http"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
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
		mediaType, _, err := mime.ParseMediaType(contentType)
		if err != nil {
			http.Error(w, "Invalid content type", http.StatusBadRequest)
			return
		}

		switch mediaType {
		case "application/json":
			reqMap = make(map[string]interface{})
			if err := json.NewDecoder(r.Body).Decode(&reqMap); err != nil {
				errMsg := fmt.Sprintf("JSON decoding error: %v", err)
				http.Error(w, errMsg, http.StatusBadRequest)
				return
			}

		case "application/x-www-form-urlencoded", "multipart/form-data":
			const maxMemory = 32 << 20 // 32MB
			if err := r.ParseMultipartForm(maxMemory); err != nil {
				errMsg := fmt.Sprintf("Form parsing error: %v", err)
				http.Error(w, errMsg, http.StatusBadRequest)
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
		msg := "안전하게 지갑 연결"

		signature := common.FromHex(r.Context().Value("request").(map[string]interface{})["Signature"].(string))
		reqWalletAddr := common.HexToAddress(r.Context().Value("request").(map[string]interface{})["WalletAddress"].(string))

		if signature[64] != 27 && signature[64] != 28 {
			http.Error(w, "Invalid MetaMask signature: incorrect recovery id", http.StatusUnauthorized)
			return
		}
		signature[64] -= 27

		recoveredAddr, err := recoverAddressFromSignature(signature, []byte(msg))
		if err != nil {
			http.Error(w, "Invalid MetaMask signature: "+err.Error(), http.StatusUnauthorized)
			return
		}

		if recoveredAddr.Hex() == reqWalletAddr.Hex() {
			next.ServeHTTP(w, r)
		} else {
			http.Error(w, "Invalid MetaMask signature: address mismatch", http.StatusUnauthorized)
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
