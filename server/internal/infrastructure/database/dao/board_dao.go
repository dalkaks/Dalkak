package dao

import boardvalueobject "dalkak/internal/domain/board/object/valueobject"

type BoardDao struct {
	Id        string
	Status    string
	UserId    string
	Timestamp int64

	Type    string
	TypeId  string
	Network string

	NftMetaName   string
	NftMetaDesc   string
	NftMetaExtUrl *string
	NftMetaBgCol  *string
	NftMetaAttrib []*boardvalueobject.NftAttribute

	NftImageExt *string
	NftVideoExt *string
}

type BoardFindFilter struct {
	UserId         string
	StatusIncluded *string
	StatusExcluded *string
	CategoryType   *string
	CategoryId     *string
}
