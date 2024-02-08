package interfaces

import "dalkak/pkg/dtos"

type Storage interface {
	CreatePresignedURL(mediaType dtos.MediaType, ext string) (*dtos.MediaMeta, error)
}
