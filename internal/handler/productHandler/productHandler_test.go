package productHandler

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

type mockProductService struct {
	CreateFn     func(ctx context.Context, sellerId int64, req *model.CreateProductRequest) (model.ProductResponse, error)
	GetByIdFn    func(ctx context.Context, req *model.GetProductsRequest) (model.ProductResponse, error)
	UpdateByIdFn func(ctx context.Context, sellerId int64, req *model.UpdateProductRequest) (model.ProductResponse, error)
	DeleteByIdFn func(ctx context.Context, req *model.DeleteProductRequest) error
	GetAllFn     func(ctx context.Context, page, limit int) ([]model.ProductResponse, int64, error)
	SearchFn     func(ctx context.Context, page, limit int, req *model.SearchProductsRequest) ([]model.ProductResponse, int64, error)
}

func (m *mockProductService) Create(ctx context.Context, sellerId int64, req *model.CreateProductRequest) (model.ProductResponse, error) {
	return m.CreateFn(ctx, sellerId, req)
}
func (m *mockProductService) GetById(ctx context.Context, req *model.GetProductsRequest) (model.ProductResponse, error) {
	return m.GetByIdFn(ctx, req)
}
func (m *mockProductService) UpdateById(ctx context.Context, sellerId int64, req *model.UpdateProductRequest) (model.ProductResponse, error) {
	return m.UpdateByIdFn(ctx, sellerId, req)
}
func (m *mockProductService) DeleteById(ctx context.Context, req *model.DeleteProductRequest) error {
	return m.DeleteByIdFn(ctx, req)
}
func (m *mockProductService) GetAll(ctx context.Context, page, limit int) ([]model.ProductResponse, int64, error) {
	return m.GetAllFn(ctx, page, limit)
}
func (m *mockProductService) Search(ctx context.Context, page, limit int, req *model.SearchProductsRequest) ([]model.ProductResponse, int64, error) {
	return m.SearchFn(ctx, page, limit, req)
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

func TestProductHandler_Create(t *testing.T) {
	tests := []struct {
		name           string
		body           string
		setUserID      bool
		serviceResp    model.ProductResponse
		serviceErr     error
		expectedStatus int
	}{
		{"success", `{"name":"Product","price":100,"category_id":2,"stock":10}`, true, model.ProductResponse{}, nil, http.StatusCreated},
		{"bind error", `{"name":`, true, model.ProductResponse{}, nil, http.StatusBadRequest},
		{"unauthorized", `{"name":"Product","price":100,"category_id":2,"stock":10}`, false, model.ProductResponse{}, nil, http.StatusUnauthorized},
		{"service error", `{"name":"Product","price":100,"category_id":2,"stock":10}`, true, model.ProductResponse{}, errors.New("svc"), http.StatusInternalServerError},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			svc := &mockProductService{
				CreateFn: func(ctx context.Context, sellerId int64, req *model.CreateProductRequest) (model.ProductResponse, error) {
					return tt.serviceResp, tt.serviceErr
				},
			}
			h := NewProductsHandler(svc)
			c, w := makeCtx(tt.body, http.MethodPost)
			if tt.setUserID {
				c.Set("user_id", int64(99))
			}
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
			if tt.expectedStatus == http.StatusCreated {
				out := parseJSONBody(t, w)
				if _, ok := out["data"]; !ok {
					t.Fatalf("expected data field")
				}
			}
		})
	}
}

func TestProductHandler_Get(t *testing.T) {
	tests := []struct {
		name           string
		paramValue     string
		serviceResp    model.ProductResponse
		serviceErr     error
		expectedStatus int
	}{
		{"success", "5", model.ProductResponse{}, nil, http.StatusOK},
		{"empty id", "", model.ProductResponse{}, nil, http.StatusOK},
		{"service error", "7", model.ProductResponse{}, errors.New("svc"), http.StatusInternalServerError},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			svc := &mockProductService{
				GetByIdFn: func(ctx context.Context, req *model.GetProductsRequest) (model.ProductResponse, error) {
					return tt.serviceResp, tt.serviceErr
				},
			}
			h := NewProductsHandler(svc)
			c, w := makeCtx("", http.MethodGet)
			c.Params = gin.Params{{Key: "id", Value: tt.paramValue}}
			h.Get(c)
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
					t.Fatalf("expected data")
				}
			}
		})
	}
}

func TestProductHandler_Update(t *testing.T) {
	tests := []struct {
		name           string
		paramValue     string
		body           string
		setUserID      bool
		serviceResp    model.ProductResponse
		serviceErr     error
		expectedStatus int
	}{
		{"success", "10", `{"id":10,"name":"Up","price":150,"category_id":3,"stock":5}`, true, model.ProductResponse{}, nil, http.StatusOK},
		{"bind error", "10", `{"name":`, true, model.ProductResponse{}, nil, http.StatusBadRequest},
		{"unauthorized", "10", `{"id":10,"name":"Up","price":150,"category_id":3,"stock":5}`, false, model.ProductResponse{}, nil, http.StatusUnauthorized},
		{"service error", "10", `{"id":10,"name":"Up","price":150,"category_id":3,"stock":5}`, true, model.ProductResponse{}, errors.New("svc"), http.StatusInternalServerError},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			svc := &mockProductService{
				UpdateByIdFn: func(ctx context.Context, sellerId int64, req *model.UpdateProductRequest) (model.ProductResponse, error) {
					return tt.serviceResp, tt.serviceErr
				},
			}
			h := NewProductsHandler(svc)
			c, w := makeCtx(tt.body, http.MethodPut)
			c.Params = gin.Params{{Key: "id", Value: tt.paramValue}}
			if tt.setUserID {
				c.Set("user_id", int64(55))
			}
			h.Update(c)
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
					t.Fatalf("expected data")
				}
			}
		})
	}
}

func TestProductHandler_Delete(t *testing.T) {
	tests := []struct {
		name           string
		paramValue     string
		serviceErr     error
		expectedStatus int
	}{
		{"success", "8", nil, http.StatusOK},
		{"empty id", "", nil, http.StatusOK},
		{"service error", "11", errors.New("svc"), http.StatusInternalServerError},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			svc := &mockProductService{
				DeleteByIdFn: func(ctx context.Context, req *model.DeleteProductRequest) error {
					return tt.serviceErr
				},
			}
			h := NewProductsHandler(svc)
			c, w := makeCtx("", http.MethodDelete)
			c.Params = gin.Params{{Key: "id", Value: tt.paramValue}}
			h.Delete(c)
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
				if out["data"] != "product deleted" {
					t.Fatalf("expected product deleted")
				}
			}
		})
	}
}

func TestProductHandler_GetAll(t *testing.T) {
	tests := []struct {
		name           string
		query          string
		serviceResp    []model.ProductResponse
		serviceTotal   int64
		serviceErr     error
		expectedStatus int
		expectedPage   int
		expectedLimit  int
	}{
		{"defaults", "", []model.ProductResponse{{}}, 5, nil, http.StatusOK, 1, 20},
		{"with params", "page=2&limit=10", []model.ProductResponse{{}, {}}, 50, nil, http.StatusOK, 2, 10},
		{"service error", "", nil, 0, errors.New("svc"), http.StatusInternalServerError, 1, 20},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			svc := &mockProductService{
				GetAllFn: func(ctx context.Context, page, limit int) ([]model.ProductResponse, int64, error) {
					return tt.serviceResp, tt.serviceTotal, tt.serviceErr
				},
			}
			h := NewProductsHandler(svc)
			c, w := makeCtx("", http.MethodGet)
			if tt.query != "" {
				c.Request.URL.RawQuery = tt.query
			}
			h.GetAll(c)
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
				page := int(out["page"].(float64))
				limit := int(out["limit"].(float64))
				if page != tt.expectedPage || limit != tt.expectedLimit {
					t.Fatalf("page/limit got %d/%d want %d/%d", page, limit, tt.expectedPage, tt.expectedLimit)
				}
				if _, ok := out["data"]; !ok {
					t.Fatalf("expected data")
				}
			}
		})
	}
}

func TestProductHandler_Search(t *testing.T) {
	tests := []struct {
		name           string
		query          string
		serviceResp    []model.ProductResponse
		serviceTotal   int64
		serviceErr     error
		expectedStatus int
		expectedPage   int
		expectedLimit  int
	}{
		{"success", "q=phone&page=2&limit=5", []model.ProductResponse{{}}, 12, nil, http.StatusOK, 2, 5},
		{"defaults", "", []model.ProductResponse{{}}, 3, nil, http.StatusOK, 1, 20},
		{"service error", "q=phone", nil, 0, errors.New("svc"), http.StatusInternalServerError, 1, 20},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			svc := &mockProductService{
				SearchFn: func(ctx context.Context, page, limit int, req *model.SearchProductsRequest) ([]model.ProductResponse, int64, error) {
					return tt.serviceResp, tt.serviceTotal, tt.serviceErr
				},
			}
			h := NewProductsHandler(svc)
			c, w := makeCtx("", http.MethodGet)
			if tt.query != "" {
				c.Request.URL.RawQuery = tt.query
			}
			h.Search(c)
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
				page := int(out["page"].(float64))
				limit := int(out["limit"].(float64))
				if page != tt.expectedPage || limit != tt.expectedLimit {
					t.Fatalf("page/limit got %d/%d want %d/%d", page, limit, tt.expectedPage, tt.expectedLimit)
				}
				if _, ok := out["data"]; !ok {
					t.Fatalf("expected data")
				}
			}
		})
	}
}
