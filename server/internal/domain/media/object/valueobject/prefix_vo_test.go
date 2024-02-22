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
