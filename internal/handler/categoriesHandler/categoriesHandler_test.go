package categoriesHandler

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

type mockCategoriesService struct {
	CreateFn  func(ctx context.Context, req *model.CreateCategoryRequest) (*model.CategoryResponse, error)
	GetByIdFn func(ctx context.Context, req *model.GetCategoryByIdRequest) (*model.CategoryResponse, error)
	UpdateFn  func(ctx context.Context, req *model.UpdateCategoryRequest) (*model.CategoryResponse, error)
	DeleteFn  func(ctx context.Context, req *model.DeleteCategoryRequest) error
	GetAllFn  func(ctx context.Context) (*[]model.CategoryResponse, error)
}

func (m *mockCategoriesService) Create(ctx context.Context, req *model.CreateCategoryRequest) (*model.CategoryResponse, error) {
	return m.CreateFn(ctx, req)
}
func (m *mockCategoriesService) GetById(ctx context.Context, req *model.GetCategoryByIdRequest) (*model.CategoryResponse, error) {
	return m.GetByIdFn(ctx, req)
}
func (m *mockCategoriesService) Update(ctx context.Context, req *model.UpdateCategoryRequest) (*model.CategoryResponse, error) {
	return m.UpdateFn(ctx, req)
}
func (m *mockCategoriesService) Delete(ctx context.Context, req *model.DeleteCategoryRequest) error {
	return m.DeleteFn(ctx, req)
}
func (m *mockCategoriesService) GetAll(ctx context.Context) (*[]model.CategoryResponse, error) {
	return m.GetAllFn(ctx)
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

func TestCategoriesHandler_Create(t *testing.T) {
	tests := []struct {
		name           string
		body           string
		serviceResp    *model.CategoryResponse
		serviceErr     error
		expectedStatus int
	}{
		{"success", `{"name":"Books"}`, &model.CategoryResponse{}, nil, http.StatusCreated},
		{"bind error", `{"name":`, nil, nil, http.StatusBadRequest},
		{"service error", `{"name":"Toys"}`, nil, errors.New("svc"), http.StatusInternalServerError},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			svc := &mockCategoriesService{
				CreateFn: func(ctx context.Context, req *model.CreateCategoryRequest) (*model.CategoryResponse, error) {
					return tt.serviceResp, tt.serviceErr
				},
			}
			h := NewCategoryHandler(svc)
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
			if tt.expectedStatus == http.StatusCreated {
				out := parseJSONBody(t, w)
				if _, ok := out["data"]; !ok {
					t.Fatalf("expected data field")
				}
			}
		})
	}
}

func TestCategoriesHandler_GetById(t *testing.T) {
	tests := []struct {
		name           string
		paramValue     string
		serviceResp    *model.CategoryResponse
		serviceErr     error
		expectedStatus int
	}{
		{"success", "3", &model.CategoryResponse{}, nil, http.StatusOK},
		{"bad param", "abc", nil, nil, http.StatusBadRequest},
		{"service error", "5", nil, errors.New("svc"), http.StatusInternalServerError},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			svc := &mockCategoriesService{
				GetByIdFn: func(ctx context.Context, req *model.GetCategoryByIdRequest) (*model.CategoryResponse, error) {
					return tt.serviceResp, tt.serviceErr
				},
			}
			h := NewCategoryHandler(svc)
			c, w := makeCtx("", http.MethodGet)
			c.Params = gin.Params{{Key: "id", Value: tt.paramValue}}
			h.GetById(c)
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

func TestCategoriesHandler_Update(t *testing.T) {
	tests := []struct {
		name           string
		paramValue     string
		body           string
		serviceResp    *model.CategoryResponse
		serviceErr     error
		expectedStatus int
	}{
		{"success", "4", `{"id":4,"name":"Electronics"}`, &model.CategoryResponse{}, nil, http.StatusOK},
		{"bad param", "notint", `{"id":5,"name":"X"}`, nil, nil, http.StatusBadRequest},
		{"bind error", "4", `{"name":`, nil, nil, http.StatusBadRequest},
		{"service error", "4", `{"id":4,"name":"Home"}`, nil, errors.New("svc"), http.StatusInternalServerError},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			svc := &mockCategoriesService{
				UpdateFn: func(ctx context.Context, req *model.UpdateCategoryRequest) (*model.CategoryResponse, error) {
					return tt.serviceResp, tt.serviceErr
				},
			}
			h := NewCategoryHandler(svc)
			c, w := makeCtx(tt.body, http.MethodPut)
			c.Params = gin.Params{{Key: "id", Value: tt.paramValue}}
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
					t.Fatalf("expected data field")
				}
			}
		})
	}
}

func TestCategoriesHandler_Delete(t *testing.T) {
	tests := []struct {
		name           string
		paramValue     string
		serviceErr     error
		expectedStatus int
	}{
		{"success", "6", nil, http.StatusOK},
		{"bad param", "x", nil, http.StatusBadRequest},
		{"service error", "7", errors.New("svc"), http.StatusInternalServerError},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			svc := &mockCategoriesService{
				DeleteFn: func(ctx context.Context, req *model.DeleteCategoryRequest) error {
					return tt.serviceErr
				},
			}
			h := NewCategoryHandler(svc)
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
				if out["data"] != true {
					t.Fatalf("expected data true")
				}
			}
		})
	}
}

func TestCategoriesHandler_GetAll(t *testing.T) {
	tests := []struct {
		name           string
		serviceResp    *[]model.CategoryResponse
		serviceErr     error
		expectedStatus int
	}{
		{"success", &[]model.CategoryResponse{{}}, nil, http.StatusOK},
		{"service error", nil, errors.New("svc"), http.StatusInternalServerError},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			svc := &mockCategoriesService{
				GetAllFn: func(ctx context.Context) (*[]model.CategoryResponse, error) {
					return tt.serviceResp, tt.serviceErr
				},
			}
			h := NewCategoryHandler(svc)
			c, w := makeCtx("", http.MethodGet)
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
				if _, ok := out["data"]; !ok {
					t.Fatalf("expected data field")
				}
			}
		})
	}
}
