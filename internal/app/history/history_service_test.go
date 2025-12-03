package history_test

import (
	"bytes"
	"errors"
	"testing"
	"time"

	"github.com/google/uuid"

	"warehousecontrol/internal/app/history"
	dhist "warehousecontrol/internal/domain/history"
	"warehousecontrol/internal/domain/item"
)

type fakeRepo struct{ called bool }

func (f *fakeRepo) GetItemsHistory(id string, from, to time.Time, action string, login string) ([]*dhist.History, error) {
	f.called = true
	return []*dhist.History{}, nil
}

func TestGetItems_ValidatesRange(t *testing.T) {
	svc := history.NewHistoryService(&fakeRepo{})
	from := time.Date(2025, 1, 2, 0, 0, 0, 0, time.UTC)
	to := time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC)
	if _, err := svc.GetItems("", from, to, "", ""); err == nil {
		t.Fatalf("expected error for from > to")
	}
}

func TestGetItems_ValidatesAction(t *testing.T) {
	svc := history.NewHistoryService(&fakeRepo{})
	if _, err := svc.GetItems("", time.Time{}, time.Time{}, "bad", ""); err == nil {
		t.Fatalf("expected invalid action error")
	}
}

func TestGetItems_ValidatesLoginMinLen(t *testing.T) {
	svc := history.NewHistoryService(&fakeRepo{})
	if _, err := svc.GetItems("", time.Time{}, time.Time{}, "", "ab"); err == nil {
		t.Fatalf("expected login length error")
	}
}

func TestGetItems_ValidatesUUID(t *testing.T) {
	svc := history.NewHistoryService(&fakeRepo{})
	if _, err := svc.GetItems("not-a-uuid", time.Time{}, time.Time{}, "", ""); err == nil {
		t.Fatalf("expected uuid parse error")
	}
}

func TestGetItems_PassesThroughOnValid(t *testing.T) {
	fr := &fakeRepo{}
	svc := history.NewHistoryService(fr)
	from := time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC)
	to := time.Date(2025, 1, 2, 0, 0, 0, 0, time.UTC)

	if _, err := svc.GetItems("", from, to, "created", "john"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !fr.called {
		t.Fatalf("expected repo to be called")
	}
}

func TestGetItems_AllowsEmptyLogin(t *testing.T) {
	fr := &fakeRepo{}
	svc := history.NewHistoryService(fr)

	if _, err := svc.GetItems("", time.Time{}, time.Time{}, "", ""); err != nil {
		t.Fatalf("did not expect error for empty login, got %v", err)
	}
}

func TestGetItems_AllowsValidUUID(t *testing.T) {
	fr := &fakeRepo{}
	svc := history.NewHistoryService(fr)

	validUUID := "123e4567-e89b-12d3-a456-426614174000"

	if _, err := svc.GetItems(validUUID, time.Time{}, time.Time{}, "", ""); err != nil {
		t.Fatalf("did not expect error for valid UUID: %v", err)
	}
}

func TestGetItems_ValidActions(t *testing.T) {
	actions := []string{"created", "updated", "deleted"}

	for _, a := range actions {
		fr := &fakeRepo{}
		svc := history.NewHistoryService(fr)

		if _, err := svc.GetItems("", time.Time{}, time.Time{}, a, ""); err != nil {
			t.Fatalf("expected action %s to be valid, got error %v", a, err)
		}
		if !fr.called {
			t.Fatalf("repo should be called for action %s", a)
		}
	}
}

func TestGetItems_NoFilters(t *testing.T) {
	fr := &fakeRepo{}
	svc := history.NewHistoryService(fr)

	if _, err := svc.GetItems("", time.Time{}, time.Time{}, "", ""); err != nil {
		t.Fatalf("unexpected error for empty filters: %v", err)
	}
	if !fr.called {
		t.Fatalf("expected repo to be called")
	}
}

type fakeRepoCSV struct {
	result []*dhist.History
	err    error
}

func (f *fakeRepoCSV) GetItemsHistory(id string, from, to time.Time, action string, login string) ([]*dhist.History, error) {
	return f.result, f.err
}

func TestGetItemsCSV_WritesHeader(t *testing.T) {
	repo := &fakeRepoCSV{result: []*dhist.History{}}
	svc := history.NewHistoryService(repo)

	var buf bytes.Buffer
	err := svc.GetItemsCSV("", time.Time{}, time.Time{}, "", "", &buf)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	output := buf.String()
	expectedHeader := "ID,ItemID,Action,ChangedBy,ChangedByLogin,ChangedAt,OldItemSnapshot,NewItemSnapshot"

	if output[:len(expectedHeader)] != expectedHeader {
		t.Fatalf("expected header %q, got %q", expectedHeader, output)
	}
}
func TestGetItemsCSV_WritesRows(t *testing.T) {
	h := &dhist.History{
		ID:              uuid.New(),
		ItemID:          uuid.New(),
		Action:          "created",
		ChangedBy:       uuid.New(),
		ChangedByLogin:  "john",
		ChangedAt:       time.Date(2025, 1, 1, 12, 0, 0, 0, time.UTC),
		OldItemSnapshot: item.Item{Name: "old"},
		NewItemSnapshot: item.Item{Name: "new"},
	}
	repo := &fakeRepoCSV{result: []*dhist.History{h}}
	svc := history.NewHistoryService(repo)

	var buf bytes.Buffer
	err := svc.GetItemsCSV("", time.Time{}, time.Time{}, "", "", &buf)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	out := buf.String()

	if !contains(out, h.ID.String()) ||
		!contains(out, h.ItemID.String()) ||
		!contains(out, "created") ||
		!contains(out, "john") ||
		!contains(out, "2025-01-01T12:00:00Z") {
		t.Fatalf("CSV row is missing expected fields:\n%s", out)
	}
}

func TestGetItemsCSV_ErrFromGetItems(t *testing.T) {
	repo := &fakeRepoCSV{err: errors.New("boom")}
	svc := history.NewHistoryService(repo)

	var buf bytes.Buffer
	err := svc.GetItemsCSV("", time.Time{}, time.Time{}, "", "", &buf)
	if err == nil {
		t.Fatalf("expected error from GetItems")
	}
}

func contains(s, substr string) bool {
	return bytes.Contains([]byte(s), []byte(substr))
}
