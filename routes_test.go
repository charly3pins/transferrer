package transferrer

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
)

type StoreOK struct{}

func (s StoreOK) Account(email string) (Account, error) {
	return Account{Balance: 10000}, nil
}

func (s StoreOK) Move(t Transfer) error {
	return nil
}

type StoreKO struct{}

func (s StoreKO) Account(email string) (Account, error) {
	return Account{}, fmt.Errorf("Account error")
}

func (s StoreKO) Move(t Transfer) error {
	return fmt.Errorf("Transfer error")
}

func TestHandler_BalanceOK(t *testing.T) {
	sample := `{"foo": "bar"}`
	handler := Handler{Store: StoreOK{}}
	req := newTestRequest(sample)
	req.Header.Set("Authorization", "Bearer SuperToken")
	c, rec := newTestContext(req)
	c.Set(emailContextKey, "foo@email.com")
	handler.Balance(c)
	if http.StatusOK != rec.Result().StatusCode {
		t.Errorf("Expected: %v Actual: %v", http.StatusOK, rec.Result().StatusCode)
	}
}

func TestHandler_BalanceEmailContextKO(t *testing.T) {
	sample := `{"foo": "bar"}`
	handler := Handler{Store: StoreKO{}}
	req := newTestRequest(sample)
	c, rec := newTestContext(req)
	handler.Balance(c)
	if http.StatusUnauthorized != rec.Result().StatusCode {
		t.Errorf("Expected: %v Actual: %v", http.StatusUnauthorized, rec.Result().StatusCode)
	}
}

func TestHandler_BalanceKO(t *testing.T) {
	sample := `{"foo": "bar"}`
	handler := Handler{Store: StoreKO{}}
	req := newTestRequest(sample)
	req.Header.Set("Authorization", "Bearer SuperToken")
	c, rec := newTestContext(req)
	c.Set(emailContextKey, "foo@email.com")
	handler.Balance(c)
	if http.StatusInternalServerError != rec.Result().StatusCode {
		t.Errorf("Expected: %v Actual: %v", http.StatusInternalServerError, rec.Result().StatusCode)
	}
}

func TestHandler_TransferBindKO(t *testing.T) {
	sample := `{"foo":"bar"}`
	handler := Handler{Store: StoreOK{}}
	req := newTestRequest(sample)
	req.Header.Set("Authorization", "Bearer SuperToken")
	c, rec := newTestContext(req)
	c.Set(emailContextKey, "foo@email.com")
	handler.Transfer(c)
	if http.StatusBadRequest != rec.Result().StatusCode {
		t.Errorf("Expected: %v Actual: %v", http.StatusBadRequest, rec.Result().StatusCode)
	}
}

func TestHandler_TransferBalanceNotEnough(t *testing.T) {
	sample := `{
		"originUser": "sender@mail.com",
		"originNumber": "fake1",
		"destinationUser": "receiver@email.com",
		"destinationNumber": "fake2",
		"amount": 10000000000
	}`
	handler := Handler{Store: StoreOK{}}
	req := newTestRequest(sample)
	req.Header.Set("Authorization", "Bearer SuperToken")
	c, rec := newTestContext(req)
	c.Set(emailContextKey, "foo@email.com")
	handler.Transfer(c)
	if http.StatusInternalServerError != rec.Result().StatusCode {
		t.Errorf("Expected: %v Actual: %v", http.StatusInternalServerError, rec.Result().StatusCode)
	}
}

func TestHandler_TransferOK(t *testing.T) {
	sample := `{
		"originUser": "sender@mail.com",
		"originNumber": "fake1",
		"destinationUser": "receiver@email.com",
		"destinationNumber": "fake2",
		"amount": 1000
	}`
	handler := Handler{Store: StoreOK{}}
	req := newTestRequest(sample)
	req.Header.Set("Authorization", "Bearer SuperToken")
	c, rec := newTestContext(req)
	c.Set(emailContextKey, "foo@email.com")
	handler.Transfer(c)
	if http.StatusOK != rec.Result().StatusCode {
		t.Errorf("Expected: %v Actual: %v", http.StatusOK, rec.Result().StatusCode)
	}
}

func newTestContext(request *http.Request) (*gin.Context, *httptest.ResponseRecorder) {
	gin.SetMode(gin.TestMode)
	rec := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(rec)
	c.Request = request
	return c, rec
}

func newTestRequest(data string) *http.Request {
	req, _ := http.NewRequest("POST", "http://test.localhost:8080", strings.NewReader(data))
	req.Header.Set("Content-Type", "application/json")
	return req
}
