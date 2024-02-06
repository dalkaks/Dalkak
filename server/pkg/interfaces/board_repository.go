package interfaces

import "dalkak/pkg/dtos"

type BoardRepository interface {
  CreateBoardImage(dto *dtos.BoardImageDto, userId string) error
}
