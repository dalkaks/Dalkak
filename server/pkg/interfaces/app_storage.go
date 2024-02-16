package interfaces

import "dalkak/pkg/dtos"

type Storage interface {
	GetHeadObject(key string) (*dtos.MediaHeadDto, error)
	DeleteObject(key string) error
	
	CreatePresignedURL(userId string, dto *dtos.UploadMediaDto) (*dtos.MediaMeta, string, error)
}
