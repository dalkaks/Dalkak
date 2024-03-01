package boardentity

import (
	generateutil "dalkak/pkg/utils/generate"
	timeutil "dalkak/pkg/utils/time"
)

type BoardEntity struct {
	Id        string      `json:"id"`
	UserId    string      `json:"userId"`
	Timestamp int64       `json:"timestamp"`
	Status    BoardStatus `json:"status"`
	Type      BoardType   `json:"type"`
}

type BoardStatus string

const (
	BoardCreated         BoardStatus = "created"
	PaymentStatusPaid    BoardStatus = "paid"
	PaymentStatusFailed  BoardStatus = "payFailed"
	ContractUploaded     BoardStatus = "contractUpload"
	ContractUploadFailed BoardStatus = "contractUploadFailed"
	NFTUploaded          BoardStatus = "nftUpload"
	NFTUploadFailed      BoardStatus = "nftUploadFailed"
	BoardPosted          BoardStatus = "posted"
	BoardDeleted         BoardStatus = "deleted"
)

type BoardType string

const (
	BoardDefaultNft BoardType = "defaultNft"
	BoardCustomNft  BoardType = "customNft"
)

func NewBoardEntity(boardType BoardType, userId string) *BoardEntity {
	return &BoardEntity{
		Id:        generateutil.GenerateUUID(),
		UserId:    userId,
		Timestamp: timeutil.GetTimestamp(),
		Status:    BoardCreated,
		Type:      boardType,
	}
}

func (bs BoardStatus) String() string {
	return string(bs)
}

func (bt BoardType) String() string {
	return string(bt)
}
