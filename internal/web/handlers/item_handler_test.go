package handlers_test

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	wbgin "github.com/wb-go/wbf/ginext"

	ditem "warehousecontrol/internal/domain/item"
	"warehousecontrol/internal/web/handlers"
)

type MockItemService struct {
	CreateFn   func(name string, count int, price float64, userID string, login string) (*ditem.Item, error)
	GetItemsFn func() ([]*ditem.Item, error)
	GetItemFn  func(id string) (*ditem.Item, error)
	PutItemFn  func(id string, name string, count int, price float64, userID string, login string) (*ditem.Item, error)
	DelItemFn  func(id string, userID string, login string) error
}

func (m *MockItemService) Create(name string, count int, price float64, userID string, login string) (*ditem.Item, error) {
	return m.CreateFn(name, count, price, userID, login)
}
func (m *MockItemService) GetItems() ([]*ditem.Item, error)       { return m.GetItemsFn() }
func (m *MockItemService) GetItem(id string) (*ditem.Item, error) { return m.GetItemFn(id) }
func (m *MockItemService) PutItem(id string, name string, count int, price float64, userID string, login string) (*ditem.Item, error) {
	return m.PutItemFn(id, name, count, price, userID, login)
}
func (m *MockItemService) DeleteItem(id string, userID string, login string) error {
	return m.DelItemFn(id, userID, login)
}

func performJSON(hf func(*wbgin.Context), method, path string, body any, setCtx func(*wbgin.Context)) *httptest.ResponseRecorder {
	var buf bytes.Buffer
	if body != nil {
		_ = json.NewEncoder(&buf).Encode(body)
	}
	req := httptest.NewRequest(method, path, &buf)
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(rr)
	ctx.Request = req
	if setCtx != nil {
		setCtx(ctx)
	}
	hf(ctx)
	return rr
}

func TestItemHandler_CreateItem_Success(t *testing.T) {
	mock := &MockItemService{
		CreateFn: func(name string, count int, price float64, userID string, login string) (*ditem.Item, error) {
			return &ditem.Item{ID: ditem.Item{}.ID, Name: name, Count: count, Price: price}, nil
		},
	}
	h := handlers.NewItemHandler(mock)
	body := map[string]any{"name": "A", "count": 1, "price": 1.0}
	rr := performJSON(h.CreateItem, http.MethodPost, "/api/items", body, func(c *wbgin.Context) { c.Set("userId", "uid"); c.Set("login", "john") })
	if rr.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rr.Code)
	}
}

func TestItemHandler_CreateItem_InvalidJSON(t *testing.T) {
	h := handlers.NewItemHandler(&MockItemService{})
	req := httptest.NewRequest(http.MethodPost, "/api/items", bytes.NewBufferString("{bad json"))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(rr)
	ctx.Request = req
	ctx.Request.Header = req.Header
	ctx.Request.Body = req.Body
	h.CreateItem(ctx)
	if rr.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", rr.Code)
	}
}

func TestItemHandler_CreateItem_MissingContext(t *testing.T) {
	h := handlers.NewItemHandler(&MockItemService{})
	body := map[string]any{"name": "A", "count": 1, "price": 1.0}
	rr := performJSON(h.CreateItem, http.MethodPost, "/api/items", body, nil)
	if rr.Code != http.StatusInternalServerError {
		t.Fatalf("expected 500, got %d", rr.Code)
	}
}

func TestItemHandler_CreateItem_ServiceError(t *testing.T) {
	mock := &MockItemService{CreateFn: func(name string, count int, price float64, userID string, login string) (*ditem.Item, error) {
		return nil, errors.New("svc err")
	}}
	h := handlers.NewItemHandler(mock)
	body := map[string]any{"name": "A", "count": 1, "price": 1.0}
	rr := performJSON(h.CreateItem, http.MethodPost, "/api/items", body, func(c *wbgin.Context) { c.Set("userId", "uid"); c.Set("login", "john") })
	if rr.Code != http.StatusInternalServerError {
		t.Fatalf("expected 500, got %d", rr.Code)
	}
}

func TestItemHandler_GetItems_Success(t *testing.T) {
	mock := &MockItemService{GetItemsFn: func() ([]*ditem.Item, error) { return []*ditem.Item{{Name: "A"}}, nil }}
	h := handlers.NewItemHandler(mock)
	rr := performJSON(h.GetItems, http.MethodGet, "/api/items", nil, nil)
	if rr.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rr.Code)
	}
}

func TestItemHandler_GetItems_ServiceError(t *testing.T) {
	mock := &MockItemService{GetItemsFn: func() ([]*ditem.Item, error) { return nil, errors.New("svc err") }}
	h := handlers.NewItemHandler(mock)
	rr := performJSON(h.GetItems, http.MethodGet, "/api/items", nil, nil)
	if rr.Code != http.StatusInternalServerError {
		t.Fatalf("expected 500, got %d", rr.Code)
	}
}

func TestItemHandler_GetItem_Success(t *testing.T) {
	mock := &MockItemService{GetItemFn: func(id string) (*ditem.Item, error) { return &ditem.Item{Name: "A"}, nil }}
	h := handlers.NewItemHandler(mock)
	req := httptest.NewRequest(http.MethodGet, "/api/items/123", nil)
	rr := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(rr)
	ctx.Request = req
	ctx.AddParam("id", "123")
	h.GetItem(ctx)
	if rr.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rr.Code)
	}
}

func TestItemHandler_GetItem_ServiceError(t *testing.T) {
	mock := &MockItemService{GetItemFn: func(id string) (*ditem.Item, error) { return nil, errors.New("svc err") }}
	h := handlers.NewItemHandler(mock)
	req := httptest.NewRequest(http.MethodGet, "/api/items/123", nil)
	rr := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(rr)
	ctx.Request = req
	ctx.AddParam("id", "123")
	h.GetItem(ctx)
	if rr.Code != http.StatusInternalServerError {
		t.Fatalf("expected 500, got %d", rr.Code)
	}
}

func TestItemHandler_PutItem_Success(t *testing.T) {
	mock := &MockItemService{PutItemFn: func(id string, name string, count int, price float64, userID string, login string) (*ditem.Item, error) {
		return &ditem.Item{Name: name, Count: count, Price: price}, nil
	}}
	h := handlers.NewItemHandler(mock)
	body := map[string]any{"name": "B", "count": 2, "price": 2.5}
	reqBody, _ := json.Marshal(body)
	req := httptest.NewRequest(http.MethodPut, "/api/items/123", bytes.NewReader(reqBody))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(rr)
	ctx.Request = req
	ctx.AddParam("id", "123")
	ctx.Set("userId", "uid")
	ctx.Set("login", "john")
	h.PutItem(ctx)
	if rr.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rr.Code)
	}
}

func TestItemHandler_PutItem_InvalidJSON(t *testing.T) {
	h := handlers.NewItemHandler(&MockItemService{})
	req := httptest.NewRequest(http.MethodPut, "/api/items/123", bytes.NewBufferString("{bad json"))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(rr)
	ctx.Request = req
	ctx.AddParam("id", "123")
	h.PutItem(ctx)
	if rr.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", rr.Code)
	}
}

func TestItemHandler_PutItem_MissingContext(t *testing.T) {
	h := handlers.NewItemHandler(&MockItemService{})
	body := map[string]any{"name": "B", "count": 2, "price": 2.5}
	rr := performJSON(h.PutItem, http.MethodPut, "/api/items/123", body, func(c *wbgin.Context) { c.AddParam("id", "123") })
	if rr.Code != http.StatusInternalServerError {
		t.Fatalf("expected 500, got %d", rr.Code)
	}
}

func TestItemHandler_PutItem_ServiceError(t *testing.T) {
	mock := &MockItemService{PutItemFn: func(id string, name string, count int, price float64, userID string, login string) (*ditem.Item, error) {
		return nil, errors.New("svc err")
	}}
	h := handlers.NewItemHandler(mock)
	body := map[string]any{"name": "B", "count": 2, "price": 2.5}
	rr := performJSON(h.PutItem, http.MethodPut, "/api/items/123", body, func(c *wbgin.Context) { c.AddParam("id", "123"); c.Set("userId", "uid"); c.Set("login", "john") })
	if rr.Code != http.StatusInternalServerError {
		t.Fatalf("expected 500, got %d", rr.Code)
	}
}

func TestItemHandler_DeleteItem_Success(t *testing.T) {
	mock := &MockItemService{DelItemFn: func(id string, userID string, login string) error { return nil }}
	h := handlers.NewItemHandler(mock)
	req := httptest.NewRequest(http.MethodDelete, "/api/items/123", nil)
	rr := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(rr)
	ctx.Request = req
	ctx.AddParam("id", "123")
	ctx.Set("userId", "uid")
	ctx.Set("login", "john")
	h.DeleteItem(ctx)
	if rr.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rr.Code)
	}
}

func TestItemHandler_DeleteItem_MissingContext(t *testing.T) {
	h := handlers.NewItemHandler(&MockItemService{})
	req := httptest.NewRequest(http.MethodDelete, "/api/items/123", nil)
	rr := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(rr)
	ctx.Request = req
	ctx.AddParam("id", "123")
	h.DeleteItem(ctx)
	if rr.Code != http.StatusInternalServerError {
		t.Fatalf("expected 500, got %d", rr.Code)
	}
}

func TestItemHandler_DeleteItem_ServiceError(t *testing.T) {
	mock := &MockItemService{DelItemFn: func(id string, userID string, login string) error { return errors.New("svc err") }}
	h := handlers.NewItemHandler(mock)
	req := httptest.NewRequest(http.MethodDelete, "/api/items/123", nil)
	rr := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(rr)
	ctx.Request = req
	ctx.AddParam("id", "123")
	ctx.Set("userId", "uid")
	ctx.Set("login", "john")
	h.DeleteItem(ctx)
	if rr.Code != http.StatusInternalServerError {
		t.Fatalf("expected 500, got %d", rr.Code)
	}
}
