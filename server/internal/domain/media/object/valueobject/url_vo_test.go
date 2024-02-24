package mediavalueobject

import (
	"testing"
)

// mediaTempKey 생성 테스트
func TestGenerateMediaTempKey(t *testing.T) {
	userId := "0x1234567890abcdef"
	resource, err := NewMediaResource("board", "image/jpg")
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	expectedKey := "temp/0x1234567890abcdef/board/image/image.jpg"
	generatedKey, err := GenerateMediaTempKey(userId, resource)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if generatedKey != expectedKey {
		t.Errorf("Expected key to be %s, got %s", expectedKey, generatedKey)
	}
}
