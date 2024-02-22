package mediavalueobject

import responseutil "dalkak/pkg/utils/response"

var AllowedPrefixes = map[string]bool{
	"board": true,
}

type Prefix string

func NewPrefix(prefix string) (Prefix, error) {
	newPrefix := Prefix(prefix)
	if !newPrefix.IsAllowedPrefix() {
		return "", responseutil.NewAppError(responseutil.ErrCodeBadRequest, responseutil.ErrMsgRequestInvalid)
	}
	return newPrefix, nil
}

func (prefix Prefix) String() string {
	return string(prefix)
}

func (prefix Prefix) IsAllowedPrefix() bool {
	_, ok := AllowedPrefixes[prefix.String()]
return ok
}
