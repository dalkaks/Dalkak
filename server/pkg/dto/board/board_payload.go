package boarddto

import (
	boardentity "dalkak/internal/domain/board/object/entity"
	boardvalueobject "dalkak/internal/domain/board/object/valueobject"
)

type CreateBoardRequest struct {
	Name            string                          `json:"name" validate:"required"`
	Description     string                          `json:"description" validate:"required"`
	ImageId         string                          `json:"imageId"`
	VideoId         string                          `json:"videoId"`
	ExternalLink    string                          `json:"externalLink"`
	BackgroundColor string                          `json:"backgroundColor"`
	Attributes      []boardvalueobject.NftAttribute `json:"attributes"`
}

type CreateBoardResponse struct {
	OrderId       string                  `json:"orderId"`
	BoardStatus   boardentity.BoardStatus `json:"boardStatus"`
	OrderName     string                  `json:"orderName"`
	OriginPrice   int64                   `json:"originPrice"`
	DiscountPrice int64                   `json:"discountPrice"`
	PaymentPrice  int64                   `json:"paymentPrice"`
}

func NewCreateBoardResponse(orderId string, boardStatus boardentity.BoardStatus, orderName string, originPrice, discountPrice, paymentPrice int64) *CreateBoardResponse {
	return &CreateBoardResponse{
		OrderId:       orderId,
		BoardStatus:   boardStatus,
		OrderName:     orderName,
		OriginPrice:   originPrice,
		DiscountPrice: discountPrice,
		PaymentPrice:  paymentPrice,
	}
}
