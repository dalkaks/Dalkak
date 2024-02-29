package boarddto

import (
	boardaggregate "dalkak/internal/domain/board/object/aggregate"
	boardentity "dalkak/internal/domain/board/object/entity"
	boardvalueobject "dalkak/internal/domain/board/object/valueobject"
	orderaggregate "dalkak/internal/domain/order/object/aggregate"
)

type CreateBoardRequest struct {
	Name            string                           `json:"name" validate:"required"`
	Description     string                           `json:"description" validate:"required"`
	ImageId         *string                          `json:"imageId"`
	VideoId         *string                          `json:"videoId"`
	ExternalLink    *string                          `json:"externalLink"`
	BackgroundColor *string                          `json:"backgroundColor"`
	Attributes      []*boardvalueobject.NftAttribute `json:"attributes"`
}

type CreateBoardResponse struct {
	BoardId       string                  `json:"boardId"`
	BoardStatus   boardentity.BoardStatus `json:"boardStatus"`
	OrderId       string                  `json:"orderId"`
	OrderName     string                  `json:"orderName"`
	OriginPrice   int64                   `json:"originPrice"`
	DiscountPrice int64                   `json:"discountPrice"`
	PaymentPrice  int64                   `json:"paymentPrice"`
}

func NewCreateBoardResponse(board *boardaggregate.BoardAggregate, order *orderaggregate.OrderAggregate) *CreateBoardResponse {
	return &CreateBoardResponse{
		BoardId:       board.BoardEntity.Id,
		BoardStatus:   board.BoardEntity.Status,
		OrderId:       order.OrderEntity.Id,
		OrderName:     order.OrderEntity.Name,
		OriginPrice:   order.OrderPrice.OriginPrice,
		DiscountPrice: order.OrderPrice.DiscountPrice,
		PaymentPrice:  order.OrderPrice.PaymentPrice,
	}
}
