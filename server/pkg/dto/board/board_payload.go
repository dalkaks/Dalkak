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
