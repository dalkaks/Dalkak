package mediadomain

import "dalkak/internal/infrastructure/database/dao"

type MediaRepository interface {
	FindMediaTemp(userId, mediaType, prefix string) (*dao.MediaTempDao, error)
}
