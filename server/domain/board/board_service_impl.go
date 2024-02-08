package board

import (
	"dalkak/pkg/dtos"
	"dalkak/pkg/interfaces"
	"dalkak/pkg/payloads"
	"errors"
)

type BoardServiceImpl struct {
	mode    string
	domain  string
	db      interfaces.BoardRepository
	storage interfaces.Storage
}

const boardStoragePath = "board"

func NewBoardService(mode string, domain string, db interfaces.Database, storage interfaces.Storage) *BoardServiceImpl {
	boardRepo := NewBoardRepository(db)

	return &BoardServiceImpl{
		mode:    mode,
		domain:  domain,
		db:      boardRepo,
		storage: storage,
	}
}

func (service *BoardServiceImpl) CreatePresignedURL(dto *payloads.BoardUploadMediaRequest, userInfo *dtos.UserInfo) (*payloads.BoardUploadMediaResponse, error) {
	// req 검증
	if dto.IsValid() == false {
		return nil, errors.New("invalid request")
	}

	// s3 presigned url 생성
	mediaType, err := dtos.ToMediaType(dto.MediaType)
	if err != nil {
		return nil, err
	}
	mediaMeta, err := service.storage.CreatePresignedURL(mediaType, dto.Ext)
	if err != nil {
		return nil, err
	}

	return &payloads.BoardUploadMediaResponse{
		Id:  mediaMeta.ID,
		Url: mediaMeta.URL,
	}, nil
}
