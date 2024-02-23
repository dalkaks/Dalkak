package mediaentity

import (
	generateutil "dalkak/pkg/utils/generate"
	timeutil "dalkak/pkg/utils/time"
)

type MediaEntity struct {
	Id          string      `json:"id"`
	IsConfirm   bool        `json:"isConfirm"`
	Timestamp   int64       `json:"timestamp"`
}

func NewMediaEntity() *MediaEntity {
	return &MediaEntity{
		Id:          generateutil.GenerateUUID(),
		IsConfirm:   false,
		Timestamp:   timeutil.GetTimestamp(),
	}
}

func ConvertMediaEntity(Id string, IsConfirm bool, Timestamp int64) *MediaEntity {
	return &MediaEntity{
		Id:          Id,
		IsConfirm:   IsConfirm,
		Timestamp:   Timestamp,
	}
}

func (media *MediaEntity) IsPublic() bool {
	return media.IsConfirm
}

func (media *MediaEntity) CheckId(id string) bool {
	return media.Id == id	
}
