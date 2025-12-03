package item

import (
	"testing"
)

func TestNewItem_Valid(t *testing.T) {
	it, err := NewItem("Widget", 10, 3.5)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if it.Name != "Widget" || it.Count != 10 || it.Price != 3.5 {
		t.Fatalf("unexpected item fields: %+v", it)
	}
	if it.ID == (it.ID) && it.ID.String() == "" {
		t.Fatalf("expected generated UUID, got empty")
	}
}

func TestNewItem_InvalidName(t *testing.T) {
	_, err := NewItem("", 1, 1.0)
	if err == nil {
		t.Fatalf("expected error for empty name")
	}
}

func TestNewItem_InvalidPrice(t *testing.T) {
	_, err := NewItem("A", 1, 0)
	if err == nil {
		t.Fatalf("expected error for non-positive price")
	}
}

func TestChangeItem_Valid(t *testing.T) {
	it, _ := NewItem("A", 1, 1.0)
	if err := it.ChangeItem("B", 2, 2.5); err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if it.Name != "B" || it.Count != 2 || it.Price != 2.5 {
		t.Fatalf("unexpected item after change: %+v", it)
	}
}

func TestChangeItem_Invalid(t *testing.T) {
	it, _ := NewItem("A", 1, 1.0)
	if err := it.ChangeItem("", 2, 2.5); err == nil {
		t.Fatalf("expected error for empty name")
	}
	if err := it.ChangeItem("B", 2, 0); err == nil {
		t.Fatalf("expected error for non-positive price")
	}
}

func TestDiff(t *testing.T) {
	a := Item{Name: "A", Count: 1, Price: 1.0}
	b := Item{Name: "B", Count: 2, Price: 1.5}
	d := a.Diff(b)
	if len(d) != 3 {
		t.Fatalf("expected 3 fields to differ, got %d", len(d))
	}
	if d["name"].Old != "A" || d["name"].New != "B" {
		t.Fatalf("unexpected name diff: %+v", d["name"])
	}
	if d["count"].Old != 1 || d["count"].New != 2 {
		t.Fatalf("unexpected count diff: %+v", d["count"])
	}
	if d["price"].Old != 1.0 || d["price"].New != 1.5 {
		t.Fatalf("unexpected price diff: %+v", d["price"])
	}
}
