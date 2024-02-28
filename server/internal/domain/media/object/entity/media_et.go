package mediaentity

import (
	generateutil "dalkak/pkg/utils/generate"
	timeutil "dalkak/pkg/utils/time"
)

type MediaEntity struct {
	Id        string `json:"id"`
	IsConfirm bool   `json:"isConfirm"`
	Timestamp int64  `json:"timestamp"`
}

type MediaEntityOption func(*MediaEntity)

func NewMediaEntity(options ...MediaEntityOption) *MediaEntity {
	me := &MediaEntity{
		Id:        generateutil.GenerateUUID(),
		IsConfirm: false,
		Timestamp: timeutil.GetTimestamp(),
	}
	for _, option := range options {
		option(me)
	}
	return me
}

func ConvertMediaEntity(Id string, IsConfirm bool, Timestamp int64) *MediaEntity {
	return &MediaEntity{
		Id:        Id,
		IsConfirm: IsConfirm,
		Timestamp: Timestamp,
	}
}

func (media *MediaEntity) IsPublic() bool {
	return media.IsConfirm
}

func (media *MediaEntity) CheckConfirm() bool {
	return media.IsConfirm
}

func (media *MediaEntity) SetConfirm() {
	media.IsConfirm = true
	media.Timestamp = timeutil.GetTimestamp()
}

func (media *MediaEntity) CheckId(id string) bool {
	return media.Id == id
}
