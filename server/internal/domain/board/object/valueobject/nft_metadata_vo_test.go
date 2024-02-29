package boardvalueobject

import (
	"testing"
)

func TestNewNftMetadata(t *testing.T) {
	// validateMetadataName 테스트
	t.Run("validateMetadataName", func(t *testing.T) {
		tests := []struct {
			name     string
			expected bool
		}{
			{"ValidName", true},
			{"", false},
			{"ThisNameIsWayTooLongForTheGivenLimitAndShouldFailValidation", false},
		}

		for _, test := range tests {
			err := validateMetadataName(test.name)
			if (err == nil) != test.expected {
				t.Errorf("validateMetadataName(%s) expected %v, got %v", test.name, test.expected, err == nil)
			}
		}
	})

	// validateMetadataDescription 테스트
	t.Run("validateMetadataDescription", func(t *testing.T) {
		tests := []struct {
			description string
			expected    bool
		}{
			{"This is a valid description.", true},
			{"", false},
			{"This description is way too long and should fail the validation because it exceeds the maximum allowed length for a description in the system. It keeps going and going and going and going and just doesn't stop, because it's way too long.", false},
		}

		for _, test := range tests {
			err := validateMetadataDescription(test.description)
			if (err == nil) != test.expected {
				t.Errorf("validateMetadataDescription(%s) expected %v, got %v", test.description, test.expected, err == nil)
			}
		}
	})

	// validateMetadataExternalUrl 테스트
	t.Run("validateMetadataExternalUrl", func(t *testing.T) {
		tests := []struct {
			url      string
			expected bool
		}{
			{"http://www.example.com", true},
			{"https://example.com", true},
			{"https://www.example.com", true},
			{"https:/www.example.com", false},
			{"ftp://www.example.com", false},
		}

		for _, test := range tests {
			err := validateMetadataExternalUrl(test.url)
			if (err == nil) != test.expected {
				t.Errorf("validateMetadataExternalUrl(%s) expected %v, got %v", test.url, test.expected, err == nil)
			}
		}
	})

	// validateMetadataBackgroundColor 테스트
	t.Run("validateMetadataBackgroundColor", func(t *testing.T) {
		tests := []struct {
			color    string
			expected bool
		}{
			{"000000", true},
			{"ffffff", true},
			{"#000000", false},
			{"#fffff", false},
			{"&00000", false},
		}

		for _, test := range tests {
			err := validateMetadataBackgroundColor(test.color)
			if (err == nil) != test.expected {
				t.Errorf("validateMetadataBackgroundColor(%s) expected %v, got %v", test.color, test.expected, err == nil)
			}
		}
	})

	// NewNftMetadata 테스트
	t.Run("NewNftMetadata", func(t *testing.T) {
		tests := []struct {
			name            string
			description     string
			externalUrl     *string
			attributes      []*NftAttribute
			backgroundColor *string
			expected        bool
		}{
			{"ValidName", "ValidDescription", nil, nil, nil, true},
			{"", "ValidDescription", nil, nil, nil, false},
			{"ValidName", "", nil, nil, nil, false},
		}

		for _, test := range tests {
			_, err := NewNftMetadata(test.name, test.description, test.externalUrl, test.backgroundColor, test.attributes)
			if (err == nil) != test.expected {
				t.Errorf("NewNftMetadata(%s, %s, %v, %v, %v) expected %v, got %v", test.name, test.description, test.externalUrl, test.backgroundColor, test.attributes, test.expected, err == nil)
			}
		}
	})
}

func TestValidateNftAttribute(t *testing.T) {
	// validateNftAttribute 문자 테스트
	t.Run("validateNftAttribute string", func(t *testing.T) {
		test := []struct {
			attr     NftAttribute
			expected bool
		}{
			{NftAttribute{TraitType: strPtr("trait"), Value: "value"}, true},
			{NftAttribute{Value: "value"}, true},
			{NftAttribute{TraitType: strPtr(""), Value: "value"}, false},
			{NftAttribute{TraitType: strPtr("trait"), Value: ""}, false},
			{NftAttribute{TraitType: strPtr("trait"), Value: "value", DisplayType: strPtr("date")}, false},
		}

		for _, test := range test {
			err := validateNftAttribute(&test.attr)
			if (err == nil) != test.expected {
				t.Errorf("validateNftAttribute(%v) expected %v, got %v", test.attr, test.expected, err == nil)
			}
		}
	})

	// validateNftAttribute 날짜 테스트
	t.Run("validateNftAttribute date", func(t *testing.T) {
		test := []struct {
			attr     NftAttribute
			expected bool
		}{
			{NftAttribute{TraitType: strPtr("trait"), Value: 1546360800, DisplayType: strPtr("date")}, true},
			{NftAttribute{TraitType: strPtr("trait"), Value: "1546360800", DisplayType: strPtr("date")}, false},
			{NftAttribute{TraitType: strPtr("trait"), Value: "value", DisplayType: strPtr("date")}, false},
		}

		for _, test := range test {
			err := validateNftAttribute(&test.attr)
			if (err == nil) != test.expected {
				t.Errorf("validateNftAttribute(%v) expected %v, got %v", test.attr, test.expected, err == nil)
			}
		}
	})

	// validateNftAttribute 숫자 테스트
	t.Run("validateNftAttribute number", func(t *testing.T) {
		test := []struct {
			attr     NftAttribute
			expected bool
		}{
			{NftAttribute{TraitType: strPtr("trait"), Value: 10}, true},
			{NftAttribute{TraitType: strPtr("trait"), Value: 10, DisplayType: strPtr("number")}, true},
			{NftAttribute{TraitType: strPtr("trait"), Value: 10, DisplayType: strPtr("boost_number")}, true},
			{NftAttribute{TraitType: strPtr("trait"), Value: 10, DisplayType: strPtr("boost_percentage")}, true},
			{NftAttribute{TraitType: strPtr("trait"), Value: 10, DisplayType: strPtr("invalid")}, false},
			{NftAttribute{Value: 10}, false},
			{NftAttribute{Value: 10, DisplayType: strPtr("number")}, false},
			{NftAttribute{TraitType: strPtr("trait"), Value: 10.5}, true},
			{NftAttribute{TraitType: strPtr("trait"), Value: 0}, true},
			{NftAttribute{TraitType: strPtr("trait"), Value: -10}, false},
			{NftAttribute{TraitType: strPtr("trait"), Value: 110, DisplayType: strPtr("boost_percentage")}, false},
		}

		for _, test := range test {
			err := validateNftAttribute(&test.attr)
			if (err == nil) != test.expected {
				t.Errorf("validateNftAttribute(%v) expected %v, got %v", test.attr, test.expected, err == nil)
			}
		}
	})
}

func strPtr(s string) *string {
	return &s
}
