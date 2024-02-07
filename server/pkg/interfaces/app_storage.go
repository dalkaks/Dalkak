package interfaces

import "dalkak/pkg/dtos"

type Storage interface {
  Upload(media *dtos.MediaDto, path string) (*dtos.MediaMeta, error)
}
