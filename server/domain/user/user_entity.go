package user

type UserTable struct {
	WalletAddress string
	Timestamp     int64
}

const UserTableName = "user"
const WalletAddressKey = "WalletAddress"
