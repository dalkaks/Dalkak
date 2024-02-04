package interfaces

import (
	"dalkak/pkg/dtos"
	"dalkak/pkg/payloads"
)

type BoardService interface {
	UploadImage(media *dtos.MediaDto, userInfo *dtos.UserInfo) (*payloads.BoardUploadMediaResponse, error)
}
