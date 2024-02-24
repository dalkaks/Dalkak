package mediavalueobject

import (
	"testing"
)

func TestNewPrefix(t *testing.T) {
	// 허용된 prefix 테스트
	t.Run("Allowed prefix", func(t *testing.T) {
		prefix, err := NewPrefix("board")
		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}
		if prefix.String() != "board" {
			t.Errorf("Expected prefix to be 'board', got '%s'", prefix)
		}
	})

	// 허용되지 않은 prefix 테스트
	t.Run("Disallowed prefix", func(t *testing.T) {
		_, err := NewPrefix("unallowed")
		if err == nil {
			t.Error("Expected error for disallowed prefix, got none")
		}
	})
}

func TestNewContentType(t *testing.T) {
	// 허용된 콘텐츠 유형 테스트
	t.Run("Allowed content types", func(t *testing.T) {
		tests := []struct {
			contentType string
			expected    string
		}{
			{"image/jpg", "image/jpg"},
			{"image/jpeg", "image/jpeg"},
			{"video/mp4", "video/mp4"},
		}

		for _, test := range tests {
			contentType, err := NewContentType(test.contentType)
			if err != nil {
				t.Errorf("Expected no error, got %v", err)
			}
			if contentType.String() != test.expected {
				t.Errorf("Expected %s, got %s", test.expected, contentType.String())
			}
		}
	})

	// 비허용 콘텐츠 유형 및 부적합한 입력 테스트
	t.Run("Disallowed content types", func(t *testing.T) {
		tests := []string{
			"text/plain",
			"image/test",
			"",
			"/",
			"image",
			"image/",
			"/jpg",
		}
		for _, test := range tests {
			_, err := NewContentType(test)
			if err == nil {
				t.Errorf("Expected error for disallowed content type, got none")
			}
		}
	})
}

func TestContentTypeString(t *testing.T) {
	// 문자열 변환 테스트
	t.Run("String conversion and validation", func(t *testing.T) {
		contentType, _ := NewContentType("image/jpeg")

		mediaType := contentType.ConvertToMediaType()
		if mediaType != "image" {
			t.Errorf("Expected 'image', got '%s'", mediaType)
		}

		extension := contentType.ConvertToExtension()
		if extension != "jpeg" {
			t.Errorf("Expected 'jpeg', got '%s'", extension)
		}
	})
}
