package blockchaintype

type NftMetadata struct {
	Name            string         `json:"name" validate:"required"`
	Description     string         `json:"description" validate:"required"`
	Image           string         `json:"image" validate:"required"`
	ExternalUrl     string         `json:"externalUrl"`
	Attributes      []NftAttribute `json:"attributes"`
	BackgroundColor string         `json:"backgroundColor"`
	AnimationUrl    string         `json:"animationUrl"`
}

/**
* NftAttribute is a struct that represents the attribute of NFT.
* 서비스 내에서 NftAttribute 설정은 다음과 같다.

1. 문자일 경우
TraitType, Value 사용
TraitType이 없다면 "PROPERTY"로 설정

2. 숫자일 경우
TraitType, Value, DisplayType(optional) 사용
value는 정수형 및 소수형을 모두 지원한다.
DisplayType는 Ranking폼(null),Boosts폼("boostNumber", "boostPercentage"),Stats폼("number")

3. 날짜일 경우
TraitType, Value, DisplayType 사용
Value는 Unix timestamp로 설정("1546360800")
DisplayType는 "date"로 설정
*/
type NftAttribute struct {
	TraitType   string      `json:"traitType,omitempty"`
	Value       interface{} `json:"value" validate:"required"`
	DisplayType string      `json:"displayType,omitempty"`
}


// todo param validation
func NewNftMetadata(name, description, image, externalUrl, backgroundColor, animationUrl string, attributes []NftAttribute) *NftMetadata {
	return &NftMetadata{
		Name:            name,
		Description:     description,
		Image:           image,
		ExternalUrl:     externalUrl,
		Attributes:      attributes,
		BackgroundColor: backgroundColor,
		AnimationUrl:    animationUrl,
	}
}