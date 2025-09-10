package userHandler

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/niklvrr/myMarketplace/internal/model"
)

type mockService struct {
	SignUpFn          func(ctx context.Context, req *model.SighUpRequest) (string, error)
	LoginFn           func(ctx context.Context, req *model.LoginRequest) (string, error)
	GetUserByIdFn     func(ctx context.Context, req *model.GetUserByIdRequest) (model.UserResponse, error)
	UpdateUserByIdFn  func(ctx context.Context, req *model.UpdateUserByIdRequest, role string) (model.UserResponse, error)
	GetUserByEmailFn  func(ctx context.Context, req *model.GetUserByEmailRequest) (model.UserResponse, error)
	BlockUserByIdFn   func(ctx context.Context, req *model.BlockUserByIdRequest) error
	UnblockUserByIdFn func(ctx context.Context, req *model.UnblockUserByIdRequest) error
	GetAllUsersFn     func(ctx context.Context) ([]model.UserResponse, error)
	UpdateUserRoleFn  func(ctx context.Context, req *model.UpdateUserRoleRequest) error
	ApproveProductFn  func(ctx context.Context, req *model.ApproveProductRequest) error
	LogoutFn          func(ctx context.Context, req *model.LogoutRequest) error
}

func (m *mockService) SignUp(ctx context.Context, req *model.SighUpRequest) (string, error) {
	return m.SignUpFn(ctx, req)
}
func (m *mockService) Login(ctx context.Context, req *model.LoginRequest) (string, error) {
	return m.LoginFn(ctx, req)
}
func (m *mockService) GetUserById(ctx context.Context, req *model.GetUserByIdRequest) (model.UserResponse, error) {
	return m.GetUserByIdFn(ctx, req)
}
func (m *mockService) UpdateUserById(ctx context.Context, req *model.UpdateUserByIdRequest, role string) (model.UserResponse, error) {
	return m.UpdateUserByIdFn(ctx, req, role)
}
func (m *mockService) GetUserByEmail(ctx context.Context, req *model.GetUserByEmailRequest) (model.UserResponse, error) {
	return m.GetUserByEmailFn(ctx, req)
}
func (m *mockService) BlockUserById(ctx context.Context, req *model.BlockUserByIdRequest) error {
	return m.BlockUserByIdFn(ctx, req)
}
func (m *mockService) UnblockUserById(ctx context.Context, req *model.UnblockUserByIdRequest) error {
	return m.UnblockUserByIdFn(ctx, req)
}
func (m *mockService) GetAllUsers(ctx context.Context) ([]model.UserResponse, error) {
	return m.GetAllUsersFn(ctx)
}
func (m *mockService) UpdateUserRole(ctx context.Context, req *model.UpdateUserRoleRequest) error {
	return m.UpdateUserRoleFn(ctx, req)
}
func (m *mockService) ApproveProduct(ctx context.Context, req *model.ApproveProductRequest) error {
	return m.ApproveProductFn(ctx, req)
}
func (m *mockService) Logout(ctx context.Context, req *model.LogoutRequest) error {
	return m.LogoutFn(ctx, req)
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

func TestUserHandler_SignUp(t *testing.T) {
	tests := []struct {
		name           string
		body           string
		serviceToken   string
		serviceErr     error
		expectedStatus int
		expectedData   interface{}
	}{
		{"success", `{"name":"Alice","email":"a@b.com","password":"password123"}`, "tok", nil, http.StatusCreated, "tok"},
		{"bind error", `{"email":`, "", nil, http.StatusBadRequest, nil},
		{"service error", `{"name":"Bob","email":"b@b.com","password":"password123"}`, "", errors.New("svc"), http.StatusInternalServerError, nil},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			svc := &mockService{
				SignUpFn: func(ctx context.Context, req *model.SighUpRequest) (string, error) {
					return tt.serviceToken, tt.serviceErr
				},
			}
			h := NewUserHandler(svc)
			c, w := makeCtx(tt.body, http.MethodPost)
			h.SignUp(c)
			if tt.expectedStatus == http.StatusInternalServerError {
				if w.Code < 500 || w.Code >= 600 {
					t.Fatalf("expected 5xx status, got %d", w.Code)
				}
				return
			}
			if w.Code != tt.expectedStatus {
				t.Fatalf("status got %d want %d body: %s", w.Code, tt.expectedStatus, w.Body.String())
			}
			if tt.expectedData != nil {
				out := parseJSONBody(t, w)
				if out["data"] != tt.expectedData {
					t.Fatalf("data got %v want %v", out["data"], tt.expectedData)
				}
			}
		})
	}
}

func TestUserHandler_Login(t *testing.T) {
	tests := []struct {
		name           string
		body           string
		serviceToken   string
		serviceErr     error
		expectedStatus int
		expectedKey    string
	}{
		{"success", `{"email":"a@b.com","password":"password123"}`, "tok", nil, http.StatusOK, "token"},
		{"bind error", `{bad`, "", nil, http.StatusBadRequest, ""},
		{"service error", `{"email":"a@b.com","password":"password123"}`, "", errors.New("svc"), http.StatusInternalServerError, ""},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			svc := &mockService{
				LoginFn: func(ctx context.Context, req *model.LoginRequest) (string, error) {
					return tt.serviceToken, tt.serviceErr
				},
			}
			h := NewUserHandler(svc)
			c, w := makeCtx(tt.body, http.MethodPost)
			h.Login(c)
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

func TestUserHandler_GetUserById(t *testing.T) {
	tests := []struct {
		name           string
		setUserID      bool
		serviceUser    model.UserResponse
		serviceErr     error
		expectedStatus int
	}{
		{"success", true, model.UserResponse{}, nil, http.StatusOK},
		{"no user id", false, model.UserResponse{}, nil, http.StatusBadRequest},
		{"service error", true, model.UserResponse{}, errors.New("svc"), http.StatusInternalServerError},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			svc := &mockService{
				GetUserByIdFn: func(ctx context.Context, req *model.GetUserByIdRequest) (model.UserResponse, error) {
					return tt.serviceUser, tt.serviceErr
				},
			}
			h := NewUserHandler(svc)
			c, w := makeCtx("", http.MethodGet)
			if tt.setUserID {
				c.Set("user_id", int64(5))
			}
			h.GetUserById(c)
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

func TestUserHandler_GetUserByEmail(t *testing.T) {
	tests := []struct {
		name           string
		body           string
		serviceUser    model.UserResponse
		serviceErr     error
		expectedStatus int
	}{
		{"success", `{"email":"a@b.com"}`, model.UserResponse{}, nil, http.StatusOK},
		{"bind error", `{bad`, model.UserResponse{}, nil, http.StatusBadRequest},
		{"service error", `{"email":"a@b.com"}`, model.UserResponse{}, errors.New("svc"), http.StatusInternalServerError},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			svc := &mockService{
				GetUserByEmailFn: func(ctx context.Context, req *model.GetUserByEmailRequest) (model.UserResponse, error) {
					return tt.serviceUser, tt.serviceErr
				},
			}
			h := NewUserHandler(svc)
			c, w := makeCtx(tt.body, http.MethodPost)
			h.GetUserByEmail(c)
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

func TestUserHandler_UpdateUserById(t *testing.T) {
	tests := []struct {
		name           string
		setUserID      bool
		setRole        bool
		body           string
		serviceUser    model.UserResponse
		serviceErr     error
		expectedStatus int
	}{
		{"success", true, true, `{"id":7}`, model.UserResponse{}, nil, http.StatusOK},
		{"no user id", false, true, `{}`, model.UserResponse{}, nil, http.StatusBadRequest},
		{"no role", true, false, `{}`, model.UserResponse{}, nil, http.StatusBadRequest},
		{"bind error", true, true, `{"x"`, model.UserResponse{}, nil, http.StatusBadRequest},
		{"service error", true, true, `{"id":7}`, model.UserResponse{}, errors.New("svc"), http.StatusInternalServerError},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			svc := &mockService{
				UpdateUserByIdFn: func(ctx context.Context, req *model.UpdateUserByIdRequest, role string) (model.UserResponse, error) {
					return tt.serviceUser, tt.serviceErr
				},
			}
			h := NewUserHandler(svc)
			c, w := makeCtx(tt.body, http.MethodPut)
			if tt.setUserID {
				c.Set("user_id", int64(7))
			}
			if tt.setRole {
				c.Set("role", "adm")
			}
			h.UpdateUserById(c)
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

func TestUserHandler_Logout(t *testing.T) {
	tests := []struct {
		name           string
		setUserID      bool
		serviceErr     error
		expectedStatus int
	}{
		{"success", true, nil, http.StatusOK},
		{"no user id", false, nil, http.StatusBadRequest},
		{"service error", true, errors.New("svc"), http.StatusInternalServerError},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			var capturedKey string
			svc := &mockService{
				LogoutFn: func(ctx context.Context, req *model.LogoutRequest) error {
					capturedKey = req.BlockKey
					return tt.serviceErr
				},
			}
			h := NewUserHandler(svc)
			c, w := makeCtx("", http.MethodPost)
			if tt.setUserID {
				c.Set("user_id", int64(12))
			}
			h.Logout(c)
			if tt.expectedStatus == http.StatusInternalServerError {
				if w.Code < 500 || w.Code >= 600 {
					t.Fatalf("expected 5xx status, got %d", w.Code)
				}
				return
			}
			if w.Code != tt.expectedStatus {
				t.Fatalf("status got %d want %d body: %s", w.Code, tt.expectedStatus, w.Body.String())
			}
			if tt.setUserID && tt.expectedStatus == http.StatusOK {
				expectedKey := "blacklist_user:" + strconv.Itoa(12)
				if capturedKey != expectedKey {
					t.Fatalf("logout block key got %q want %q", capturedKey, expectedKey)
				}
			}
		})
	}
}

func TestUserHandler_BlockUnblockUserById(t *testing.T) {
	t.Run("block success and logout called", func(t *testing.T) {
		var logoutKey string
		svc := &mockService{
			BlockUserByIdFn: func(ctx context.Context, req *model.BlockUserByIdRequest) error {
				return nil
			},
			LogoutFn: func(ctx context.Context, req *model.LogoutRequest) error {
				logoutKey = req.BlockKey
				return nil
			},
		}
		h := NewUserHandler(svc)
		c, w := makeCtx(`{"id":3}`, http.MethodPost)
		h.BlockUserById(c)
		if w.Code != http.StatusOK {
			t.Fatalf("status got %d want %d body: %s", w.Code, http.StatusOK, w.Body.String())
		}
		if logoutKey != "blocked_user:3" {
			t.Fatalf("logoutKey got %q want %q", logoutKey, "blocked_user:3")
		}
	})

	t.Run("block bind error", func(t *testing.T) {
		svc := &mockService{}
		h := NewUserHandler(svc)
		c, w := makeCtx(`{bad`, http.MethodPost)
		h.BlockUserById(c)
		if w.Code != http.StatusBadRequest {
			t.Fatalf("status got %d want %d", w.Code, http.StatusBadRequest)
		}
	})

	t.Run("unblock success", func(t *testing.T) {
		svc := &mockService{
			UnblockUserByIdFn: func(ctx context.Context, req *model.UnblockUserByIdRequest) error {
				return nil
			},
		}
		h := NewUserHandler(svc)
		c, w := makeCtx(`{"id":4}`, http.MethodPost)
		h.UnblockUserById(c)
		if w.Code != http.StatusOK {
			t.Fatalf("status got %d want %d", w.Code, http.StatusOK)
		}
	})

	t.Run("unblock bind error", func(t *testing.T) {
		svc := &mockService{}
		h := NewUserHandler(svc)
		c, w := makeCtx(`{bad`, http.MethodPost)
		h.UnblockUserById(c)
		if w.Code != http.StatusBadRequest {
			t.Fatalf("status got %d want %d", w.Code, http.StatusBadRequest)
		}
	})
}

func TestUserHandler_GetAllUsers(t *testing.T) {
	tests := []struct {
		name           string
		serviceUsers   []model.UserResponse
		serviceErr     error
		expectedStatus int
	}{
		{"success", []model.UserResponse{{}}, nil, http.StatusOK},
		{"service error", nil, errors.New("svc"), http.StatusInternalServerError},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			svc := &mockService{
				GetAllUsersFn: func(ctx context.Context) ([]model.UserResponse, error) {
					return tt.serviceUsers, tt.serviceErr
				},
			}
			h := NewUserHandler(svc)
			c, w := makeCtx("", http.MethodGet)
			h.GetAllUsers(c)
			if tt.expectedStatus == http.StatusInternalServerError {
				if w.Code < 500 || w.Code >= 600 {
					t.Fatalf("expected 5xx status, got %d", w.Code)
				}
				return
			}
			if w.Code != tt.expectedStatus {
				t.Fatalf("status got %d want %d", w.Code, tt.expectedStatus)
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

func TestUserHandler_UpdateUserRole(t *testing.T) {
	tests := []struct {
		name           string
		body           string
		serviceErr     error
		expectedStatus int
	}{
		{"success", `{"id":1,"role":"admin"}`, nil, http.StatusOK},
		{"bind error", `{bad`, nil, http.StatusBadRequest},
		{"service error", `{"id":1,"role":"admin"}`, errors.New("svc"), http.StatusInternalServerError},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			svc := &mockService{
				UpdateUserRoleFn: func(ctx context.Context, req *model.UpdateUserRoleRequest) error {
					return tt.serviceErr
				},
			}
			h := NewUserHandler(svc)
			c, w := makeCtx(tt.body, http.MethodPost)
			h.UpdateUserRole(c)
			if tt.expectedStatus == http.StatusInternalServerError {
				if w.Code < 500 || w.Code >= 600 {
					t.Fatalf("expected 5xx status, got %d", w.Code)
				}
				return
			}
			if w.Code != tt.expectedStatus {
				t.Fatalf("status got %d want %d body: %s", w.Code, tt.expectedStatus, w.Body.String())
			}
		})
	}
}

func TestUserHandler_ApproveProduct(t *testing.T) {
	tests := []struct {
		name           string
		body           string
		serviceErr     error
		expectedStatus int
	}{
		{"success", `{"product_id":1,"approved":true}`, nil, http.StatusOK},
		{"bind error", `{bad`, nil, http.StatusBadRequest},
		{"service error", `{"product_id":1}`, errors.New("svc"), http.StatusInternalServerError},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			svc := &mockService{
				ApproveProductFn: func(ctx context.Context, req *model.ApproveProductRequest) error {
					return tt.serviceErr
				},
			}
			h := NewUserHandler(svc)
			c, w := makeCtx(tt.body, http.MethodPost)
			h.ApproveProduct(c)
			if tt.expectedStatus == http.StatusInternalServerError {
				if w.Code < 500 || w.Code >= 600 {
					t.Fatalf("expected 5xx status, got %d", w.Code)
				}
				return
			}
			if w.Code != tt.expectedStatus {
				t.Fatalf("status got %d want %d body: %s", w.Code, tt.expectedStatus, w.Body.String())
			}
		})
	}
}
