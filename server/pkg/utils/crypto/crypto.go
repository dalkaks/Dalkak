package cryptoutil

import (
	responseutil "dalkak/pkg/utils/response"
	"fmt"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
)

type MetaMaskSignature struct {
	Signature     string `json:"signature"`
	WalletAddress string `json:"walletAddress"`
}

func VerifyMetaMaskSignature(dto *MetaMaskSignature) error {
	// todo message
	msg := "안전하게 지갑 연결"
	signature := common.FromHex(dto.Signature)
	walletAddress := common.HexToAddress(dto.WalletAddress)

	if signature[64] != 27 && signature[64] != 28 {
		return responseutil.NewAppError(responseutil.ErrCodeUnauthorized, responseutil.ErrMsgMetaMaskInvalidSignature)
	}
	signature[64] -= 27

	recoveredAddr, err := recoverAddressFromSignature(signature, []byte(msg))
	if err != nil {
		return responseutil.NewAppError(responseutil.ErrCodeUnauthorized, responseutil.ErrMsgMetaMaskInvalidSignature)
	}

	if recoveredAddr.Hex() != walletAddress.Hex() {
		return responseutil.NewAppError(responseutil.ErrCodeUnauthorized, responseutil.ErrMsgMetaMaskNotMatchAddress)
	}
	return nil
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
