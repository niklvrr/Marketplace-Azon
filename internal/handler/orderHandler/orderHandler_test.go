package orderHandler

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

type mockOrderService struct {
	CreateOrderFn            func(ctx context.Context, req *model.CreateOrderRequest) (int64, error)
	GetOrdersByUserIdFn      func(ctx context.Context, req *model.GetOrdersByUserIdRequest) (*[]model.OrderResponse, error)
	GetOrderByIdFn           func(ctx context.Context, req *model.GetOrderByIdRequest) (*model.OrderResponse, error)
	GetOrderItemsByOrderIdFn func(ctx context.Context, req *model.GetOrderItemsByOrderIdRequest) (*[]model.OrderItemResponse, error)
	DeleteOrderByIdFn        func(ctx context.Context, req *model.DeleteOrderByIdRequest) error
}

func (m *mockOrderService) CreateOrder(ctx context.Context, req *model.CreateOrderRequest) (int64, error) {
	return m.CreateOrderFn(ctx, req)
}
func (m *mockOrderService) GetOrdersByUserId(ctx context.Context, req *model.GetOrdersByUserIdRequest) (*[]model.OrderResponse, error) {
	return m.GetOrdersByUserIdFn(ctx, req)
}
func (m *mockOrderService) GetOrderById(ctx context.Context, req *model.GetOrderByIdRequest) (*model.OrderResponse, error) {
	return m.GetOrderByIdFn(ctx, req)
}
func (m *mockOrderService) GetOrderItemsByOrderId(ctx context.Context, req *model.GetOrderItemsByOrderIdRequest) (*[]model.OrderItemResponse, error) {
	return m.GetOrderItemsByOrderIdFn(ctx, req)
}
func (m *mockOrderService) DeleteOrderById(ctx context.Context, req *model.DeleteOrderByIdRequest) error {
	return m.DeleteOrderByIdFn(ctx, req)
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

func TestOrderHandler_Create(t *testing.T) {
	tests := []struct {
		name           string
		body           string
		serviceID      int64
		serviceErr     error
		expectedStatus int
		expectedKey    string
	}{
		{
			"success",
			`{"user_id":1,"order_items":[{"product_id":1,"quantity":1}],"address":"addr"}`,
			123, nil, http.StatusCreated, "order_id",
		},
		{"bind error", `{"user_id":`, 0, nil, http.StatusBadRequest, ""},
		{
			"service error",
			`{"user_id":1,"order_items":[{"product_id":1,"quantity":1}], "address":"addr"}`,
			0, errors.New("svc"), http.StatusInternalServerError, "",
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			svc := &mockOrderService{
				CreateOrderFn: func(ctx context.Context, req *model.CreateOrderRequest) (int64, error) {
					return tt.serviceID, tt.serviceErr
				},
			}
			h := NewOrderHandler(svc)
			c, w := makeCtx(tt.body, http.MethodPost)
			h.Create(c)
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

func TestOrderHandler_GetOrdersByUserId(t *testing.T) {
	tests := []struct {
		name           string
		setUserId      bool
		serviceResp    *[]model.OrderResponse
		serviceErr     error
		expectedStatus int
	}{
		{"success", true, &[]model.OrderResponse{{}}, nil, http.StatusOK},
		{"no user id", false, nil, nil, http.StatusUnauthorized},
		{"service error", true, nil, errors.New("svc"), http.StatusInternalServerError},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			svc := &mockOrderService{
				GetOrdersByUserIdFn: func(ctx context.Context, req *model.GetOrdersByUserIdRequest) (*[]model.OrderResponse, error) {
					return tt.serviceResp, tt.serviceErr
				},
			}
			h := NewOrderHandler(svc)
			c, w := makeCtx("", http.MethodGet)
			if tt.setUserId {
				c.Set("userId", int64(9))
			}
			h.GetOrdersByUserId(c)
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
				if _, ok := out["orders"]; !ok {
					t.Fatalf("expected orders field")
				}
			}
		})
	}
}

func TestOrderHandler_GetOrderById(t *testing.T) {
	tests := []struct {
		name           string
		paramValue     string
		serviceResp    *model.OrderResponse
		serviceErr     error
		expectedStatus int
	}{
		{"success", "5", &model.OrderResponse{}, nil, http.StatusOK},
		{"bad param", "x", nil, nil, http.StatusBadRequest},
		{"service error", "7", nil, errors.New("svc"), http.StatusInternalServerError},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			svc := &mockOrderService{
				GetOrderByIdFn: func(ctx context.Context, req *model.GetOrderByIdRequest) (*model.OrderResponse, error) {
					return tt.serviceResp, tt.serviceErr
				},
			}
			h := NewOrderHandler(svc)
			c, w := makeCtx("", http.MethodGet)
			c.Params = gin.Params{{Key: "id", Value: tt.paramValue}}
			h.GetOrderById(c)
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
				if _, ok := out["order"]; !ok {
					t.Fatalf("expected order field")
				}
			}
		})
	}
}

func TestOrderHandler_GetOrderItemsByOrderId(t *testing.T) {
	tests := []struct {
		name           string
		paramValue     string
		serviceResp    *[]model.OrderItemResponse
		serviceErr     error
		expectedStatus int
	}{
		{"success", "6", &[]model.OrderItemResponse{{}}, nil, http.StatusOK},
		{"bad param", "bad", nil, nil, http.StatusBadRequest},
		{"service error", "9", nil, errors.New("svc"), http.StatusInternalServerError},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			svc := &mockOrderService{
				GetOrderItemsByOrderIdFn: func(ctx context.Context, req *model.GetOrderItemsByOrderIdRequest) (*[]model.OrderItemResponse, error) {
					return tt.serviceResp, tt.serviceErr
				},
			}
			h := NewOrderHandler(svc)
			c, w := makeCtx("", http.MethodGet)
			c.Params = gin.Params{{Key: "id", Value: tt.paramValue}}
			h.GetOrderItemsByOrderId(c)
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
				if _, ok := out["order_items"]; !ok {
					t.Fatalf("expected order_items field")
				}
			}
		})
	}
}

func TestOrderHandler_DeleteOrderById(t *testing.T) {
	tests := []struct {
		name           string
		paramValue     string
		serviceErr     error
		expectedStatus int
	}{
		{"success", "8", nil, http.StatusOK},
		{"bad param", "no", nil, http.StatusBadRequest},
		{"service error", "10", errors.New("svc"), http.StatusInternalServerError},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			svc := &mockOrderService{
				DeleteOrderByIdFn: func(ctx context.Context, req *model.DeleteOrderByIdRequest) error {
					return tt.serviceErr
				},
			}
			h := NewOrderHandler(svc)
			c, w := makeCtx("", http.MethodDelete)
			c.Params = gin.Params{{Key: "id", Value: tt.paramValue}}
			h.DeleteOrderById(c)
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
				if out["status"] != true {
					t.Fatalf("expected status true")
				}
			}
		})
	}
}
