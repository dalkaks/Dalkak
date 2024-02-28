package boardvalueobject

import (
	parseutil "dalkak/pkg/utils/parse"
	responseutil "dalkak/pkg/utils/response"
	"net/url"
	"regexp"
)

const (
	maxMetadataNameLength        = 30
	maxMetadataDescriptionLength = 200
	maxMetadataAttributesLength  = 10
	maxMetadataAttributesValue   = 1e8
	maxMetadataAttributesPercent = 100
	maxAttributeTraitTypeLength  = 30
	maxAttributeValueTypeLength  = 30
)

// image, video는 media로 분리
type NftMetadata struct {
	Name            string          `json:"name" validate:"required"`
	Description     string          `json:"description" validate:"required"`
	ExternalUrl     *string         `json:"externalUrl"`
	Attributes      *[]NftAttribute `json:"attributes"`
	BackgroundColor *string         `json:"backgroundColor"`
}

/*
*
* NftAttribute is a struct that represents the attribute of NFT.
* 서비스 내에서 NftAttribute 설정은 다음과 같다.

1. 문자일 경우
TraitType, Value 사용
TraitType이 없다면 "PROPERTY"로 설정

2. 숫자일 경우
TraitType, Value, DisplayType(optional) 사용
value는 정수형 및 소수형을 모두 지원한다.
DisplayType는 Ranking폼(null),Boosts폼("boost_number", "boost_percentage"),Stats폼("number")

3. 날짜일 경우
TraitType, Value, DisplayType 사용
Value는 Unix timestamp로 설정("1546360800")
DisplayType는 "date"로 설정
*/
type NftAttribute struct {
	TraitType   *string     `json:"traitType,omitempty"`
	Value       interface{} `json:"value" validate:"required"`
	DisplayType *string     `json:"displayType,omitempty"`
}

func NewNftMetadata(name, description string, externalUrl, backgroundColor *string, attributes *[]NftAttribute) (*NftMetadata, error) {
	if err := validateMetadataName(name); err != nil {
		return nil, err
	}
	if err := validateMetadataDescription(description); err != nil {
		return nil, err
	}
	if externalUrl != nil {
		if err := validateMetadataExternalUrl(*externalUrl); err != nil {
			return nil, err
		}
	}
	if backgroundColor != nil {
		if err := validateMetadataBackgroundColor(*backgroundColor); err != nil {
			return nil, err
		}
	}
	if attributes != nil {
		if len(*attributes) > maxMetadataAttributesLength {
			return nil, responseutil.NewAppError(responseutil.ErrCodeBadRequest, responseutil.ErrMsgBoardAttributesInvalid)
		}
		for _, attr := range *attributes {
			if err := validateNftAttribute(&attr); err != nil {
				return nil, err
			}
		}
	}

	return &NftMetadata{
		Name:            name,
		Description:     description,
		ExternalUrl:     externalUrl,
		Attributes:      attributes,
		BackgroundColor: backgroundColor,
	}, nil
}

func validateMetadataName(name string) error {
	if len(name) > maxMetadataNameLength || len(name) == 0 {
		return responseutil.NewAppError(responseutil.ErrCodeBadRequest, responseutil.ErrMsgBoardNameInvalid)
	}
	return nil
}

func validateMetadataDescription(description string) error {
	if len(description) > maxMetadataDescriptionLength || len(description) == 0 {
		return responseutil.NewAppError(responseutil.ErrCodeBadRequest, responseutil.ErrMsgBoardDescriptionInvalid)
	}
	return nil
}

func validateMetadataExternalUrl(externalUrl string) error {
	parsedUrl, err := url.Parse(externalUrl)
	if err != nil {
		return responseutil.NewAppError(responseutil.ErrCodeBadRequest, responseutil.ErrMsgBoardExternalUrlInvalid)
	}

	if parsedUrl.Scheme == "" || parsedUrl.Host == "" {
		return responseutil.NewAppError(responseutil.ErrCodeBadRequest, responseutil.ErrMsgBoardExternalUrlInvalid)
	}

	if parsedUrl.Scheme != "http" && parsedUrl.Scheme != "https" {
		return responseutil.NewAppError(responseutil.ErrCodeBadRequest, responseutil.ErrMsgBoardExternalUrlInvalid)
	}

	return nil
}

func validateMetadataBackgroundColor(backgroundColor string) error {
	match, _ := regexp.MatchString("^[0-9a-fA-F]{6}$", backgroundColor)
	if !match {
		return responseutil.NewAppError(responseutil.ErrCodeBadRequest, responseutil.ErrMsgBoardBackgroundColorInvalid)
	}
	return nil
}

func validateNftAttribute(attr *NftAttribute) error {
	switch attr.Value.(type) {
	case string:
		if attr.DisplayType != nil {
			return responseutil.NewAppError(responseutil.ErrCodeBadRequest, responseutil.ErrMsgBoardAttributesInvalid)
		}
		if attr.TraitType != nil && (len(*attr.TraitType) == 0 || len(*attr.TraitType) > maxAttributeTraitTypeLength) {
			return responseutil.NewAppError(responseutil.ErrCodeBadRequest, responseutil.ErrMsgBoardAttributesInvalid)
		}
		if attr.Value == "" || len(attr.Value.(string)) > maxAttributeValueTypeLength {
			return responseutil.NewAppError(responseutil.ErrCodeBadRequest, responseutil.ErrMsgBoardAttributesInvalid)
		}

	case float64, float32, int, int64, int32, int16, int8, uint, uint64, uint32, uint16, uint8:
		if attr.TraitType == nil || *attr.TraitType == "" {
			return responseutil.NewAppError(responseutil.ErrCodeBadRequest, responseutil.ErrMsgBoardAttributesInvalid)
		}
		number, err := parseutil.ToFloat64(attr.Value)
		if err != nil {
			return responseutil.NewAppError(responseutil.ErrCodeBadRequest, responseutil.ErrMsgBoardAttributesInvalid)
		}

		if attr.DisplayType != nil && *attr.DisplayType == "date" {
			if _, err := parseutil.ToUnixTimestamp(attr.Value); err != nil {
				return responseutil.NewAppError(responseutil.ErrCodeBadRequest, responseutil.ErrMsgBoardAttributesInvalid)
			}
		}
		if err := validateDisplayTypeForNumber(attr.DisplayType, number); err != nil {
			return err
		}

	default:
		return responseutil.NewAppError(responseutil.ErrCodeBadRequest, responseutil.ErrMsgBoardAttributesInvalid)
	}

	return nil
}

func validateDisplayTypeForNumber(displayType *string, number float64) error {
	if displayType == nil {
		return validateNumberRange(number, maxMetadataAttributesValue)
	}

	switch *displayType {
	case "date":
		if _, err := parseutil.ToUnixTimestamp(number); err != nil {
			return responseutil.NewAppError(responseutil.ErrCodeBadRequest, responseutil.ErrMsgBoardAttributesInvalid)
		}
	case "number", "boost_number":
		return validateNumberRange(number, maxMetadataAttributesValue)
	case "boost_percentage":
		return validateNumberRange(number, maxMetadataAttributesPercent)
	default:
		return responseutil.NewAppError(responseutil.ErrCodeBadRequest, responseutil.ErrMsgBoardAttributesInvalid)
	}
	return nil
}

func validateNumberRange(number float64, max float64) error {
	if number < 0 || number > max {
		return responseutil.NewAppError(responseutil.ErrCodeBadRequest, responseutil.ErrMsgBoardAttributesInvalid)
	}
	return nil
}
