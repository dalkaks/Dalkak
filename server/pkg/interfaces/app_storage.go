package interfaces

import "dalkak/pkg/dtos"

type Storage interface {
	CreatePresignedURL(dto *dtos.UploadMediaDto) (*dtos.MediaMeta, error)
}
