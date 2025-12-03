package handlers_test

import (
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/gin-gonic/gin"

	dhist "warehousecontrol/internal/domain/history"
	"warehousecontrol/internal/web/handlers"
)

type MockHistoryService struct {
	GetItemsFn    func(id string, from, to time.Time, action string, login string) ([]*dhist.History, error)
	GetItemsCSVFn func(id string, from, to time.Time, action string, login string, output io.Writer) error
}

func (m *MockHistoryService) GetItems(id string, from, to time.Time, action string, login string) ([]*dhist.History, error) {
	return m.GetItemsFn(id, from, to, action, login)
}
func (m *MockHistoryService) GetItemsCSV(id string, from, to time.Time, action string, login string, output io.Writer) error {
	return m.GetItemsCSVFn(id, from, to, action, login, output)
}

func TestHistoryHandler_GetItems_InvalidFromDate(t *testing.T) {
	h := handlers.NewHistoryHandler(&MockHistoryService{})
	rr := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(rr)
	req := httptest.NewRequest(http.MethodGet, "/api/history?from=bad&to=2025-01-02", nil)
	ctx.Request = req
	h.GetItems(ctx)
	if rr.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", rr.Code)
	}
}

func TestHistoryHandler_GetItems_InvalidToDate(t *testing.T) {
	h := handlers.NewHistoryHandler(&MockHistoryService{})
	rr := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(rr)
	req := httptest.NewRequest(http.MethodGet, "/api/history?from=2025-01-01&to=bad", nil)
	ctx.Request = req
	h.GetItems(ctx)
	if rr.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", rr.Code)
	}
}

func TestHistoryHandler_GetItems_ServiceError(t *testing.T) {
	mock := &MockHistoryService{GetItemsFn: func(id string, from, to time.Time, action string, login string) ([]*dhist.History, error) {
		return nil, assertErr("boom")
	}}
	h := handlers.NewHistoryHandler(mock)
	rr := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(rr)
	req := httptest.NewRequest(http.MethodGet, "/api/history?from=2025-01-01&to=2025-01-02", nil)
	ctx.Request = req
	h.GetItems(ctx)
	if rr.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", rr.Code)
	}
}

func TestHistoryHandler_GetItems_Success(t *testing.T) {
	mock := &MockHistoryService{GetItemsFn: func(id string, from, to time.Time, action string, login string) ([]*dhist.History, error) {
		return []*dhist.History{}, nil
	}}
	h := handlers.NewHistoryHandler(mock)
	rr := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(rr)
	req := httptest.NewRequest(http.MethodGet, "/api/history?from=2025-01-01&to=2025-01-02&action=updated&login=john&id=", nil)
	ctx.Request = req
	h.GetItems(ctx)
	if rr.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rr.Code)
	}
}

func TestHistoryHandler_GetItemsCSV_InvalidFromDate(t *testing.T) {
	h := handlers.NewHistoryHandler(&MockHistoryService{})
	rr := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(rr)
	req := httptest.NewRequest(http.MethodGet, "/api/history/csv?from=bad&to=2025-01-02", nil)
	ctx.Request = req
	h.GetItemsCSV(ctx)
	if rr.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", rr.Code)
	}
}

func TestHistoryHandler_GetItemsCSV_InvalidToDate(t *testing.T) {
	h := handlers.NewHistoryHandler(&MockHistoryService{})
	rr := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(rr)
	req := httptest.NewRequest(http.MethodGet, "/api/history/csv?from=2025-01-01&to=bad", nil)
	ctx.Request = req
	h.GetItemsCSV(ctx)
	if rr.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", rr.Code)
	}
}

func TestHistoryHandler_GetItemsCSV_ServiceError(t *testing.T) {
	mock := &MockHistoryService{GetItemsCSVFn: func(id string, from, to time.Time, action string, login string, output io.Writer) error {
		return assertErr("boom")
	}}
	h := handlers.NewHistoryHandler(mock)
	rr := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(rr)
	req := httptest.NewRequest(http.MethodGet, "/api/history/csv?from=2025-01-01&to=2025-01-02", nil)
	ctx.Request = req
	h.GetItemsCSV(ctx)
	if rr.Code != http.StatusInternalServerError {
		t.Fatalf("expected 500, got %d", rr.Code)
	}
}

func TestHistoryHandler_GetItemsCSV_Success(t *testing.T) {
	mock := &MockHistoryService{GetItemsCSVFn: func(id string, from, to time.Time, action string, login string, output io.Writer) error {
		_, _ = output.Write([]byte("id,item_id,action\n"))
		_, _ = output.Write([]byte("1,abc,updated\n"))
		return nil
	}}
	h := handlers.NewHistoryHandler(mock)
	rr := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(rr)
	req := httptest.NewRequest(http.MethodGet, "/api/history/csv?from=2025-01-01&to=2025-01-02&action=updated&login=john", nil)
	ctx.Request = req
	h.GetItemsCSV(ctx)
	if rr.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rr.Code)
	}
	ct := rr.Header().Get("Content-Type")
	if !strings.Contains(ct, "text/csv") {
		t.Fatalf("expected text/csv, got %s", ct)
	}
	disp := rr.Header().Get("Content-Disposition")
	if !strings.Contains(disp, "attachment") {
		t.Fatalf("expected attachment disposition, got %s", disp)
	}
	body := rr.Body.String()
	if !strings.Contains(body, "id,item_id,action") {
		t.Fatalf("expected csv header, got %s", body)
	}
}

// minimal error type
type assertErr string

func (e assertErr) Error() string { return string(e) }
