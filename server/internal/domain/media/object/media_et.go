package mediaobject

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
