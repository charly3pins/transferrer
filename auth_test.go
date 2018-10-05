package transferrer

import (
	"net/http"
	"testing"
)

func TestHandler_GenerateJWTOK(t *testing.T) {
	sample := `{"foo": "bar"}`
	req := newTestRequest(sample)
	c, rec := newTestContext(req)
	GenerateJWT(c)
	if http.StatusOK != rec.Result().StatusCode {
		t.Errorf("Expected: %v Actual: %v", http.StatusOK, rec.Result().StatusCode)
	}
}
