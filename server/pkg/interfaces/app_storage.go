package interfaces

import "dalkak/pkg/dtos"

type Storage interface {
	CreatePresignedURL(userId string, dto *dtos.UploadMediaDto) (*dtos.MediaMeta, string, error)
	GetHeadObject(key string) (*dtos.MediaHeadDto, error)
}
