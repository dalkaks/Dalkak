package dao

import boardvalueobject "dalkak/internal/domain/board/object/valueobject"

type BoardDao struct {
	Id        string
	Status    string
	UserId    string
	Type      string
	Timestamp int64

	NftMetaName   string
	NftMetaDesc   string
	NftMetaExtUrl *string
	NftMetaBgCol  *string
	NftMetaAttrib []*boardvalueobject.NftAttribute

	NftImageExt *string
	NftVideoExt *string
}
