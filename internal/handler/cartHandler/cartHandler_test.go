package cartHandler

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/niklvrr/myMarketplace/internal/model"
)

type mockCartService struct {
	GetCartByUserIdFn      func(ctx context.Context, req *model.GetCartByUserIdRequest) (*model.CartResponse, error)
	GetCartItemsByCartIdFn func(ctx context.Context, req *model.GetCartItemsByCartIdRequest) (*[]model.CartItemResponse, error)
	AddItemFn              func(ctx context.Context, req *model.AddItemRequest) (int64, error)
	RemoveItemFn           func(ctx context.Context, req *model.RemoveItemRequest) error
	ClearCartFn            func(ctx context.Context, req *model.ClearCartRequest) error
}

func (m *mockCartService) GetCartByUserId(ctx context.Context, req *model.GetCartByUserIdRequest) (*model.CartResponse, error) {
	return m.GetCartByUserIdFn(ctx, req)
}
func (m *mockCartService) GetCartItemsByCartId(ctx context.Context, req *model.GetCartItemsByCartIdRequest) (*[]model.CartItemResponse, error) {
	return m.GetCartItemsByCartIdFn(ctx, req)
}
func (m *mockCartService) AddItem(ctx context.Context, req *model.AddItemRequest) (int64, error) {
	return m.AddItemFn(ctx, req)
}
func (m *mockCartService) RemoveItem(ctx context.Context, req *model.RemoveItemRequest) error {
	return m.RemoveItemFn(ctx, req)
}
func (m *mockCartService) ClearCart(ctx context.Context, req *model.ClearCartRequest) error {
	return m.ClearCartFn(ctx, req)
}

func init() {
	gin.SetMode(gin.ReleaseMode)
}

func makeCtx(body string, method string) (*gin.Context, *httptest.ResponseRecorder) {
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	var req *http.Request
	if body != "" {
		req = httptest.NewRequest(method, "/", bytes.NewBufferString(body))
		req.Header.Set("Content-Type", "application/json")
	} else {
		req = httptest.NewRequest(method, "/", nil)
	}
	c.Request = req
	return c, w
}

func parseJSONBody(t *testing.T, b *httptest.ResponseRecorder) map[string]interface{} {
	var out map[string]interface{}
	err := json.Unmarshal(b.Body.Bytes(), &out)
	if err != nil {
		t.Fatalf("failed to unmarshal body: %v, body: %s", err, b.Body.String())
	}
	return out
}

func TestCartHandler_GetCartByUserId(t *testing.T) {
	tests := []struct {
		name           string
		setUserID      bool
		serviceCart    *model.CartResponse
		serviceErr     error
		expectedStatus int
	}{
		{"success", true, &model.CartResponse{}, nil, http.StatusOK},
		{"no user id", false, nil, nil, http.StatusBadRequest},
		{"service error", true, nil, errors.New("svc"), http.StatusInternalServerError},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			svc := &mockCartService{
				GetCartByUserIdFn: func(ctx context.Context, req *model.GetCartByUserIdRequest) (*model.CartResponse, error) {
					return tt.serviceCart, tt.serviceErr
				},
			}
			h := NewCartHandler(svc)
			c, w := makeCtx("", http.MethodGet)
			if tt.setUserID {
				c.Set("user_id", int64(11))
			}
			h.GetCartByUserId(c)
			if tt.expectedStatus == http.StatusInternalServerError {
				if w.Code < 500 || w.Code >= 600 {
					t.Fatalf("expected 5xx status, got %d", w.Code)
				}
				return
			}
			if w.Code != tt.expectedStatus {
				t.Fatalf("status got %d want %d body: %s", w.Code, tt.expectedStatus, w.Body.String())
			}
			if tt.expectedStatus == http.StatusOK {
				out := parseJSONBody(t, w)
				if _, ok := out["data"]; !ok {
					t.Fatalf("expected data field")
				}
			}
		})
	}
}

func TestCartHandler_GetCartItemsByCartId(t *testing.T) {
	tests := []struct {
		name           string
		paramValue     string
		serviceItems   *[]model.CartItemResponse
		serviceErr     error
		expectedStatus int
	}{
		{"success", "5", &[]model.CartItemResponse{{}}, nil, http.StatusOK},
		{"bad param", "abc", nil, nil, http.StatusBadRequest},
		{"service error", "7", nil, errors.New("svc"), http.StatusInternalServerError},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			svc := &mockCartService{
				GetCartItemsByCartIdFn: func(ctx context.Context, req *model.GetCartItemsByCartIdRequest) (*[]model.CartItemResponse, error) {
					return tt.serviceItems, tt.serviceErr
				},
			}
			h := NewCartHandler(svc)
			c, w := makeCtx("", http.MethodGet)
			c.Params = gin.Params{{Key: "cart_id", Value: tt.paramValue}}
			h.GetCartItemsByCartId(c)
			if tt.expectedStatus == http.StatusInternalServerError {
				if w.Code < 500 || w.Code >= 600 {
					t.Fatalf("expected 5xx status, got %d", w.Code)
				}
				return
			}
			if w.Code != tt.expectedStatus {
				t.Fatalf("status got %d want %d body: %s", w.Code, tt.expectedStatus, w.Body.String())
			}
			if tt.expectedStatus == http.StatusOK {
				out := parseJSONBody(t, w)
				if _, ok := out["data"]; !ok {
					t.Fatalf("expected data field")
				}
			}
		})
	}
}

func TestCartHandler_AddItem(t *testing.T) {
	tests := []struct {
		name           string
		body           string
		serviceID      int64
		serviceErr     error
		expectedStatus int
		expectedKey    string
	}{
		{"success", `{"cart_id":1,"product_id":2,"quantity":1}`, 42, nil, http.StatusOK, "cart_item_id"},
		{"bind error", `{"cart_id":`, 0, nil, http.StatusBadRequest, ""},
		{"service error", `{"cart_id":1,"product_id":2,"quantity":1}`, 0, errors.New("svc"), http.StatusInternalServerError, ""},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			svc := &mockCartService{
				AddItemFn: func(ctx context.Context, req *model.AddItemRequest) (int64, error) {
					return tt.serviceID, tt.serviceErr
				},
			}
			h := NewCartHandler(svc)
			c, w := makeCtx(tt.body, http.MethodPost)
			h.AddItem(c)
			if tt.expectedStatus == http.StatusInternalServerError {
				if w.Code < 500 || w.Code >= 600 {
					t.Fatalf("expected 5xx status, got %d", w.Code)
				}
				return
			}
			if w.Code != tt.expectedStatus {
				t.Fatalf("status got %d want %d body: %s", w.Code, tt.expectedStatus, w.Body.String())
			}
			if tt.expectedKey != "" {
				out := parseJSONBody(t, w)
				if _, ok := out[tt.expectedKey]; !ok {
					t.Fatalf("expected key %s in response", tt.expectedKey)
				}
			}
		})
	}
}

func TestCartHandler_RemoveItem(t *testing.T) {
	tests := []struct {
		name           string
		body           string
		serviceErr     error
		expectedStatus int
	}{
		{"success", `{"id":3,"cart_id":1}`, nil, http.StatusOK},
		{"bind error", `{"cart_item_id":`, nil, http.StatusBadRequest},
		{"service error", `{"id":3,"cart_id":1}`, errors.New("svc"), http.StatusInternalServerError},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			svc := &mockCartService{
				RemoveItemFn: func(ctx context.Context, req *model.RemoveItemRequest) error {
					return tt.serviceErr
				},
			}
			h := NewCartHandler(svc)
			c, w := makeCtx(tt.body, http.MethodPost)
			h.RemoveItem(c)
			if tt.expectedStatus == http.StatusInternalServerError {
				if w.Code < 500 || w.Code >= 600 {
					t.Fatalf("expected 5xx status, got %d", w.Code)
				}
				return
			}
			if w.Code != tt.expectedStatus {
				t.Fatalf("status got %d want %d body: %s", w.Code, tt.expectedStatus, w.Body.String())
			}
			if tt.expectedStatus == http.StatusOK {
				out := parseJSONBody(t, w)
				if _, ok := out["status"]; !ok {
					t.Fatalf("expected status field")
				}
			}
		})
	}
}

func TestCartHandler_ClearCart(t *testing.T) {
	tests := []struct {
		name           string
		body           string
		serviceErr     error
		expectedStatus int
	}{
		{"success", `{"cart_id":1}`, nil, http.StatusOK},
		{"bind error", `{"cart_id":`, nil, http.StatusBadRequest},
		{"service error", `{"cart_id":1}`, errors.New("svc"), http.StatusInternalServerError},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			svc := &mockCartService{
				ClearCartFn: func(ctx context.Context, req *model.ClearCartRequest) error {
					return tt.serviceErr
				},
			}
			h := NewCartHandler(svc)
			c, w := makeCtx(tt.body, http.MethodPost)
			h.ClearCart(c)
			if tt.expectedStatus == http.StatusInternalServerError {
				if w.Code < 500 || w.Code >= 600 {
					t.Fatalf("expected 5xx status, got %d", w.Code)
				}
				return
			}
			if w.Code != tt.expectedStatus {
				t.Fatalf("status got %d want %d body: %s", w.Code, tt.expectedStatus, w.Body.String())
			}
			if tt.expectedStatus == http.StatusOK {
				out := parseJSONBody(t, w)
				if _, ok := out["status"]; !ok {
					t.Fatalf("expected status field")
				}
			}
		})
	}
}
