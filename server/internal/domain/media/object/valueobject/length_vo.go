package mediavalueobject

import responseutil "dalkak/pkg/utils/response"

const MaxUploadSize = 32 << 20 // 32MB

type Length int64

func NewLength(lengthInt int64) (Length, error) {
	length := Length(lengthInt)
	if !length.IsAllowedLength() {
		return 0, responseutil.NewAppError(responseutil.ErrCodeBadRequest, responseutil.ErrMsgRequestInvalid)
	}
	return length, nil
}

func (length Length) Int64() int64 {
	return int64(length)
}

func (length Length) IsAllowedLength() bool {
	return length.Int64() <= MaxUploadSize && length.Int64() > 0
}
