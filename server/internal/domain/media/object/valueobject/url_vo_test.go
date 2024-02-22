package mediavalueobject

import (
	"testing"
)

// mediaTempKey 생성 테스트
func TestGenerateMediaTempKey(t *testing.T) {
	userId := "0x1234567890abcdef"
	prefix, _ := NewPrefix("board")
	contentType, _ := NewContentType("image", "jpg")

	expectedKey := "temp/0x1234567890abcdef/board/image/image.jpg"
	generatedKey := GenerateMediaTempKey(userId, prefix, contentType)

	if generatedKey != expectedKey {
		t.Errorf("Expected key to be %s, got %s", expectedKey, generatedKey)
	}
}
