package item_test

import (
	"errors"
	"testing"

	"warehousecontrol/internal/app/item"
	"warehousecontrol/internal/config"
	domain "warehousecontrol/internal/domain/item"

	"github.com/google/uuid"
)

type fakeRepo struct {
	createItemCalled bool
	getItemCalled    bool
	putItemCalled    bool
	deleteItemCalled bool

	itemToReturn *domain.Item
	errToReturn  error
}

func (f *fakeRepo) CreateItem(i *domain.Item, userID string, login string) error {
	f.createItemCalled = true
	return f.errToReturn
}
func (f *fakeRepo) GetItems() ([]*domain.Item, error) {
	return []*domain.Item{f.itemToReturn}, f.errToReturn
}
func (f *fakeRepo) GetItem(id string) (*domain.Item, error) {
	f.getItemCalled = true
	return f.itemToReturn, f.errToReturn
}
func (f *fakeRepo) PutItem(i *domain.Item, userID string, login string) error {
	f.putItemCalled = true
	return f.errToReturn
}
func (f *fakeRepo) DeleteItem(id string, userID string, login string) error {
	f.deleteItemCalled = true
	return f.errToReturn
}

func testCfg() *config.AppConfig {
	return &config.AppConfig{
		ItemConfig: config.ItemConfig{
			NameMinLength: 2,
			NameMaxLegth:  10,
		},
	}
}

func TestCreate_Success(t *testing.T) {
	repo := &fakeRepo{}
	svc := item.NewItemService(repo, testCfg())

	obj, err := svc.Create("Apple", 5, 10.0, "uid", "login")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if obj == nil {
		t.Fatalf("expected item, got nil")
	}
	if !repo.createItemCalled {
		t.Fatalf("expected CreateItem to be called")
	}
}

func TestCreate_InvalidName(t *testing.T) {
	repo := &fakeRepo{}
	svc := item.NewItemService(repo, testCfg())

	_, err := svc.Create("A", 5, 10.0, "uid", "login")
	if err == nil {
		t.Fatalf("expected validation error")
	}
}

func TestGetItem_InvalidUUID(t *testing.T) {
	repo := &fakeRepo{}
	svc := item.NewItemService(repo, testCfg())

	_, err := svc.GetItem("not-uuid")
	if err == nil {
		t.Fatalf("expected uuid error")
	}
}

func TestGetItem_Success(t *testing.T) {
	repo := &fakeRepo{itemToReturn: &domain.Item{ID: uuid.New(), Name: "Test", Count: 1, Price: 1}}
	svc := item.NewItemService(repo, testCfg())

	_, err := svc.GetItem(uuid.NewString())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !repo.getItemCalled {
		t.Fatalf("expected GetItem to be called")
	}
}

func TestPutItem_Success(t *testing.T) {
	repo := &fakeRepo{
		itemToReturn: &domain.Item{
			ID:    uuid.New(),
			Name:  "OldName",
			Count: 1,
			Price: 1,
		},
	}
	svc := item.NewItemService(repo, testCfg())

	_, err := svc.PutItem(repo.itemToReturn.ID.String(), "NewName", 10, 5.5, "uid", "login")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !repo.getItemCalled || !repo.putItemCalled {
		t.Fatalf("expected GetItem and PutItem to be called")
	}
}

func TestPutItem_InvalidUUID(t *testing.T) {
	repo := &fakeRepo{}
	svc := item.NewItemService(repo, testCfg())

	_, err := svc.PutItem("bad-uuid", "GoodName", 1, 1, "uid", "login")
	if err == nil {
		t.Fatalf("expected uuid error")
	}
}

func TestPutItem_InvalidName(t *testing.T) {
	repo := &fakeRepo{}
	svc := item.NewItemService(repo, testCfg())

	_, err := svc.PutItem(uuid.NewString(), "A", 1, 1, "uid", "login")
	if err == nil {
		t.Fatalf("expected name error")
	}
}

func TestPutItem_RepoGetError(t *testing.T) {
	repo := &fakeRepo{errToReturn: errors.New("fail")}
	svc := item.NewItemService(repo, testCfg())

	_, err := svc.PutItem(uuid.NewString(), "ValidName", 1, 1, "uid", "login")
	if err == nil {
		t.Fatalf("expected repo get error")
	}
}

func TestDeleteItem_Success(t *testing.T) {
	repo := &fakeRepo{}
	svc := item.NewItemService(repo, testCfg())

	err := svc.DeleteItem(uuid.NewString(), "uid", "login")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !repo.deleteItemCalled {
		t.Fatalf("expected DeleteItem to be called")
	}
}

func TestDeleteItem_InvalidUUID(t *testing.T) {
	repo := &fakeRepo{}
	svc := item.NewItemService(repo, testCfg())

	err := svc.DeleteItem("bad-uuid", "uid", "login")
	if err == nil {
		t.Fatalf("expected uuid error")
	}
}

func TestIsNameValid(t *testing.T) {
	svc := item.NewItemService(&fakeRepo{}, testCfg())

	if _, err := svc.Create("X", 1, 1, "u", "l"); err == nil {
		t.Fatalf("expected name too short error")
	}

	if _, err := svc.Create("VeryLongNameHere", 1, 1, "u", "l"); err == nil {
		t.Fatalf("expected name too long error")
	}

	if _, err := svc.Create("OkName", 1, 1, "u", "l"); err != nil {
		t.Fatalf("unexpected error for valid name: %v", err)
	}
}
