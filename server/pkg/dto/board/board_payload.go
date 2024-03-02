package boarddto

import (
	boardaggregate "dalkak/internal/domain/board/object/aggregate"
	boardentity "dalkak/internal/domain/board/object/entity"
	boardvalueobject "dalkak/internal/domain/board/object/valueobject"
	mediaaggregate "dalkak/internal/domain/media/object/aggregate"
	orderaggregate "dalkak/internal/domain/order/object/aggregate"
	"dalkak/internal/infrastructure/database/dao"
)

type BoardProcessing struct {
	BoardId         string                           `json:"boardId"`
	Timestamp       int64                            `json:"timestamp"`
	Status          string                           `json:"status"`
	CategoryType    string                           `json:"categoryType"`
	CategoryId      string                           `json:"categoryId"`
	Network         string                           `json:"network"`
	Name            string                           `json:"name"`
	Description     string                           `json:"description"`
	ExternalLink    *string                          `json:"externalLink"`
	BackgroundColor *string                          `json:"backgroundColor"`
	Attributes      []*boardvalueobject.NftAttribute `json:"attributes"`
	ImageUrl        *string                          `json:"imageUrl"`
	VideoUrl        *string                          `json:"videoUrl"`
}

type CreateBoardRequest struct {
	Name            string                           `json:"name" validate:"required"`
	Description     string                           `json:"description" validate:"required"`
	CategoryType    string                           `json:"categoryType" validate:"required"`
	Network         string                           `json:"network" validate:"required"`
	ImageId         *string                          `json:"imageId" validate:"required"`
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

type GetBoardListProcessingRequest struct {
	CategoryType      string `query:"category-type" validate:"required"`
	CategoryId        string `query:"category-id"`
	Limit             int    `query:"limit" validate:"required"`
	ExclusiveStartKey string `query:"start-key"`
}

type GetBoardListProcessingResponse struct {
	Items            []*BoardProcessing `json:"items"`
	Count            int                `json:"count"`
	LastEvaluatedKey *string            `json:"lastEvaluatedKey"`
}

func NewGetBoardListProcessingResponse(items []*boardaggregate.BoardAggregate, media []*mediaaggregate.MediaNftAggregate, page *dao.ResponsePageDao) *GetBoardListProcessingResponse {
	processingItems := make([]*BoardProcessing, 0)
	for i, board := range items {
		processingItems = append(processingItems, &BoardProcessing{
			BoardId:         board.BoardEntity.Id,
			Timestamp:       board.BoardEntity.Timestamp,
			Status:          board.BoardEntity.Status.String(),
			CategoryType:    board.BoardCategory.GetCategoryType(),
			CategoryId:      board.BoardCategory.GetCategoryId(),
			Network:         board.BoardCategory.GetNetwork(),
			Name:            board.BoardMetadata.Name,
			Description:     board.BoardMetadata.Description,
			ExternalLink:    board.BoardMetadata.ExternalUrl,
			BackgroundColor: board.BoardMetadata.BackgroundColor,
			Attributes:      board.BoardMetadata.Attributes,
		})
		if media != nil {
			if media[i].MediaImageUrl != nil {
				processingItems[i].ImageUrl = &media[i].MediaImageUrl.AccessUrl
			}
			if media[i].MediaVideoResource != nil {
				processingItems[i].VideoUrl = &media[i].MediaVideoUrl.AccessUrl
			}
		}
	}
	return &GetBoardListProcessingResponse{
		Items:            processingItems,
		Count:            page.Count,
		LastEvaluatedKey: page.ExclusiveStartKey,
	}
}
