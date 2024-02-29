package dao

type TransactionDao struct {
	Id        string `json:"id"`
	Timestamp int64  `json:"timestamp"`
}

type RequestPageDao struct {
	Limit             int     `json:"limit"`
	ExclusiveStartKey *string `json:"startKey"`
}

type ResponsePageDao struct {
	Count             int     `json:"count"`
	ExclusiveStartKey *string `json:"startKey"`
}
