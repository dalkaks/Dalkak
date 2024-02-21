package userobject

import timeutil "dalkak/pkg/utils/time"

type UserEntity struct {
	WalletAddress string
	Timestamp     int64
}

func NewUserEntity(walletAddress string) *UserEntity {
	return &UserEntity{
		WalletAddress: walletAddress,
		Timestamp:     timeutil.GetTimestamp(),
	}
}