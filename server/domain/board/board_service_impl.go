package board

import (
	"dalkak/pkg/dtos"
	"dalkak/pkg/interfaces"
	"dalkak/pkg/payloads"
)

type BoardServiceImpl struct {
	mode    string
	domain  string
	db      interfaces.BoardRepository
	storage interfaces.Storage
}

func NewBoardService(mode string, domain string, db interfaces.Database, storage interfaces.Storage) *BoardServiceImpl {
	boardRepo := NewBoardRepository(db)

	return &BoardServiceImpl{
		mode:    mode,
		domain:  domain,
		db:      boardRepo,
		storage: storage,
	}
}

func (service *BoardServiceImpl) UploadImage(media *dtos.MediaDto, userInfo *dtos.UserInfo) (*payloads.BoardUploadMediaResponse, error) {
  // 이미지 업로드
  newMedia, err := service.storage.Upload(media)
  if err != nil {
    return nil, err
  }

  // 데이터베이스 저장

  // 실패 시 이미지 삭제

	// 이미지 업로드 결과 반환
	return &payloads.BoardUploadMediaResponse{
		Id: "1",
    Url: "link",
	}, nil
}
