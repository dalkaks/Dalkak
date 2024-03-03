package mediafactory

import (
	mediaaggregate "dalkak/internal/domain/media/object/aggregate"
	mediaentity "dalkak/internal/domain/media/object/entity"
	mediavalueobject "dalkak/internal/domain/media/object/valueobject"
	"dalkak/internal/infrastructure/database/dao"
	mediadto "dalkak/pkg/dto/media"
	responseutil "dalkak/pkg/utils/response"
	"fmt"
)

const boardPrefix = "board"

type MediaNftAggregateFactory interface {
	CreateMediaNftAggregateFromDto(dto *mediadto.CreateMediaNftDto, mediaImage *mediaaggregate.MediaTempAggregate, mediaVideo *mediaaggregate.MediaTempAggregate) (*mediaaggregate.MediaNftAggregate, error)
	CreateMediaNftAggregateFromBoardDaos(daos []*dao.BoardDao) ([]*mediaaggregate.MediaNftAggregate, error)
}

type CreateMediaNftFactory struct {
	staticLink string
}

func NewCreateMediaNftFactory(staticLink string) *CreateMediaNftFactory {
	return &CreateMediaNftFactory{
		staticLink: staticLink,
	}
}

func (factory *CreateMediaNftFactory) CreateMediaNftAggregate(dto *mediadto.CreateMediaNftDto, mediaImage *mediaaggregate.MediaTempAggregate, mediaVideo *mediaaggregate.MediaTempAggregate) (*mediaaggregate.MediaNftAggregate, error) {
	media := mediaentity.ConvertMediaEntity(dto.PrefixId, true, dto.Timestamp)

	var options []mediaaggregate.MediaNftAggregateOption
	if mediaImage != nil {
		options = append(options, mediaaggregate.WithMediaNftImageResource(&mediaImage.MediaResource))
		if mediaImage.MediaUrl != nil {
			newImageUrl, err := mediaImage.MediaUrl.ConvertMediaTempToFormalUrl(factory.staticLink, dto.PrefixId)
			if err != nil {
				return nil, err
			}
			options = append(options, mediaaggregate.WithMediaNftImageUrl(newImageUrl))
		}
	}

	if mediaVideo != nil {
		options = append(options, mediaaggregate.WithMediaNftVideoResource(&mediaVideo.MediaResource))
		if mediaVideo.MediaUrl != nil {
			newVideoUrl, err := mediaVideo.MediaUrl.ConvertMediaTempToFormalUrl(factory.staticLink, dto.PrefixId)
			if err != nil {
				return nil, err
			}
			options = append(options, mediaaggregate.WithMediaNftVideoUrl(newVideoUrl))
		}
	}

	mediaNftAggregate, err := mediaaggregate.NewMediaNftAggregate(media, options...)
	if err != nil {
		return nil, err
	}
	return mediaNftAggregate, nil
}

func (factory *CreateMediaNftFactory) CreateMediaNftAggregateFromBoardDaos(daos []*dao.BoardDao) ([]*mediaaggregate.MediaNftAggregate, error) {
	var mediaNfts []*mediaaggregate.MediaNftAggregate
	for _, dao := range daos {
		media := mediaentity.ConvertMediaEntity(dao.Id, true, dao.Timestamp)

		var options []mediaaggregate.MediaNftAggregateOption
		if dao.NftImageExt != nil && *dao.NftImageExt != "" {
			imageOptions, err := createMediaNftAggregateOptionFromExt(factory.staticLink, "image", *dao.NftImageExt, boardPrefix, dao.Id)
			if err != nil {
				return nil, err
			}
			options = append(options, imageOptions...)
		}
		if dao.NftVideoExt != nil && *dao.NftVideoExt != "" {
			videoOptions, err := createMediaNftAggregateOptionFromExt(factory.staticLink, "video", *dao.NftVideoExt, boardPrefix, dao.Id)
			if err != nil {
				return nil, err
			}
			options = append(options, videoOptions...)
		}

		mediaNftAggregate, err := mediaaggregate.NewMediaNftAggregate(media, options...)
		if err != nil {
			return nil, err
		}
		mediaNfts = append(mediaNfts, mediaNftAggregate)
	}
	return mediaNfts, nil
}

func createMediaNftAggregateOptionFromExt(staticLink, mediaTypeStr, ext, prefix, id string) ([]mediaaggregate.MediaNftAggregateOption, error) {
	var options []mediaaggregate.MediaNftAggregateOption
	contentType, err := mediavalueobject.NewContentType(fmt.Sprintf("%s/%s", mediaTypeStr, ext))
	if err != nil {
		return nil, err
	}
	resource, err := mediavalueobject.NewMediaResource(prefix, contentType.String())
	if err != nil {
		return nil, err
	}
	mediaUrl := mediavalueobject.ConvertMediaUrl(staticLink, prefix, id, mediaTypeStr, ext)

	if resource.GetMediaType() == "image" {
		options = append(options, mediaaggregate.WithMediaNftImageResource(resource))
		options = append(options, mediaaggregate.WithMediaNftImageUrl(mediaUrl))
		return options, nil
	} else if resource.GetMediaType() == "video" {
		options = append(options, mediaaggregate.WithMediaNftVideoResource(resource))
		options = append(options, mediaaggregate.WithMediaNftVideoUrl(mediaUrl))
		return options, nil
	}
	return nil, responseutil.NewAppError(responseutil.ErrCodeBadRequest, responseutil.ErrMsgRequestInvalid)
}
