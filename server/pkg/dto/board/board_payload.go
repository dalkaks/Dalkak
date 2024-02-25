package boarddto

import blockchaintype "dalkak/internal/infrastructure/blockchain/type"

type CreateBoardRequest struct {
	Title           string                      `json:"title" validate:"required"`
	Content         string                      `json:"content" validate:"required"`
	ImageId         string                      `json:"imageId"`
	VideoId         string                      `json:"videoId"`
	ExternalLink    string                      `json:"externalLink"`
	BackgroundColor string                      `json:"backgroundColor"`
	Attributes      blockchaintype.NftAttribute `json:"attributes"`
}

type CreateBoardResponse struct {
	OrderId       string `json:"orderId"`
	OrderStatus   string `json:"orderStatus"`
	OrderName     string `json:"orderName"`
	OriginPrice   int64  `json:"originPrice"`
	DiscountPrice int64  `json:"discountPrice"`
	PaymentPrice  int64  `json:"paymentPrice"`
}

func NewCreateBoardResponse(orderId, orderStatus, orderName string, originPrice, discountPrice, paymentPrice int64) *CreateBoardResponse {
	return &CreateBoardResponse{
		OrderId:       orderId,
		OrderStatus:   orderStatus,
		OrderName:     orderName,
		OriginPrice:   originPrice,
		DiscountPrice: discountPrice,
		PaymentPrice:  paymentPrice,
	}
}
