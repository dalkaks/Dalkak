package boardentity

import (
	generateutil "dalkak/pkg/utils/generate"
	responseutil "dalkak/pkg/utils/response"
	timeutil "dalkak/pkg/utils/time"
)

type BoardEntity struct {
	Id        string      `json:"id"`
	UserId    string      `json:"userId"`
	Timestamp int64       `json:"timestamp"`
	Status    BoardStatus `json:"status"`
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
	BoardCancelled       BoardStatus = "cancelled"
	BoardPosted          BoardStatus = "posted"
)

func NewBoardEntity(userId string) *BoardEntity {
	return &BoardEntity{
		Id:        generateutil.GenerateUUID(),
		UserId:    userId,
		Timestamp: timeutil.GetTimestamp(),
		Status:    BoardCreated,
	}
}

func ConvertBoardEntity(id, userId string, timestamp int64, statusStr string) (*BoardEntity, error) {
	status, err := NewBoardStatus(statusStr)
	if err != nil {
		return nil, err
	}
	return &BoardEntity{
		Id:        id,
		UserId:    userId,
		Timestamp: timestamp,
		Status:    status,
	}, nil
}

func (be BoardEntity) GetStatus() string {
	return string(be.Status)
}

func NewBoardStatus(statusStr string) (BoardStatus, error) {
	switch statusStr {
	case string(BoardCreated):
		return BoardCreated, nil
	case string(PaymentStatusPaid):
		return PaymentStatusPaid, nil
	case string(PaymentStatusFailed):
		return PaymentStatusFailed, nil
	case string(ContractUploaded):
		return ContractUploaded, nil
	case string(ContractUploadFailed):
		return ContractUploadFailed, nil
	case string(NFTUploaded):
		return NFTUploaded, nil
	case string(NFTUploadFailed):
		return NFTUploadFailed, nil
	case string(BoardCancelled):
		return BoardCancelled, nil
	case string(BoardPosted):
		return BoardPosted, nil
	default:
		return "", responseutil.NewAppError(responseutil.ErrCodeBadRequest, responseutil.ErrMsgRequestInvalid)
	}
}
