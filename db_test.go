package transferrer

import (
	"testing"

	_ "github.com/lib/pq"
)

func TestNewDB_KO(t *testing.T) {
	db, err := NewDB("")
	if err == nil {
		t.Errorf("Expecting error")
	}
	if db != nil {
		t.Errorf("Expecting db nil")
	}
}
