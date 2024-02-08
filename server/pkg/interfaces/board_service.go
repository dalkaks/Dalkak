package interfaces

import (
	"dalkak/pkg/dtos"
	"dalkak/pkg/payloads"
)

type BoardService interface {
	CreatePresignedURL(req *payloads.BoardUploadMediaRequest, userInfo *dtos.UserInfo) (*payloads.BoardUploadMediaResponse, error)
}
