package mediavalueobject

import (
	"testing"
)

func TestNewContentType(t *testing.T) {
	// 허용된 콘텐츠 유형 테스트
	t.Run("Allowed content types", func(t *testing.T) {
		tests := []struct {
			mediaType string
			extension string
			expected  string
		}{
			{"image", "jpg", "image/jpg"},
			{"video", "mp4", "video/mp4"},
			{"image", "png", "image/png"},
		}

		for _, test := range tests {
			contentType, err := NewContentType(test.mediaType, test.extension)
			if err != nil {
				t.Errorf("Expected no error, got %v", err)
			}
			if contentType.String() != test.expected {
				t.Errorf("Expected %s, got %s", test.expected, contentType.String())
			}
		}
	})

	// 비허용 콘텐츠 유형 테스트
	t.Run("Disallowed content types", func(t *testing.T) {
		_, err := NewContentType("text", "plain")
		if err == nil {
			t.Error("Expected error for disallowed content type, got none")
		}
	})

	// 부적합한 입력 테스트
	t.Run("Invalid input", func(t *testing.T) {
		_, err := NewContentType("", "")
		if err == nil {
			t.Error("Expected error for invalid input, got none")
		}
	})

	// jpeg jpg 구분 테스트
	t.Run("Jpeg and jpg distinction", func(t *testing.T) {
		contentType, _ := NewContentType("image", "jpeg")
		if contentType.String() != "image/jpeg" {
			t.Errorf("Expected 'image/jpeg', got '%s'", contentType.String())
		}
	})
}

func TestContentTypeString(t *testing.T) {
	// 문자열 변환 테스트
	t.Run("String conversion and validation", func(t *testing.T) {
		contentType, _ := NewContentType("image", "jpeg")

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
