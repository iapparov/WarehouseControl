package item

import (
	"warehousecontrol/internal/config"
	"warehousecontrol/internal/domain/item"

	"github.com/google/uuid"
	wbzlog "github.com/wb-go/wbf/zlog"

	"fmt"
	"unicode/utf8"
)

type ItemService struct {
	repo ItemStorageProvider
	cfg  *config.AppConfig
}

type ItemStorageProvider interface {
	CreateItem(item *item.Item, userID string, login string) error
	GetItems() ([]*item.Item, error)
	GetItem(uuid string) (*item.Item, error)
	PutItem(item *item.Item, userID string, login string) error
	DeleteItem(uuid string, userID string, login string) error
}

func NewItemService(repo ItemStorageProvider, cfg *config.AppConfig) *ItemService {
	return &ItemService{
		repo: repo,
		cfg:  cfg,
	}
}

func (s *ItemService) Create(name string, count int, price float64, userID string, login string) (*item.Item, error) {
	err := s.isNameValid(name)
	if err != nil {
		wbzlog.Logger.Warn().Err(err).Msg("item name is invalid")
		return nil, err
	}

	item, err := item.NewItem(name, count, price)
	if err != nil {
		wbzlog.Logger.Warn().Err(err).Msg("cant create item")
		return nil, err
	}
	err = s.repo.CreateItem(item, userID, login)
	if err != nil {
		return nil, err
	}

	return item, nil
}

func (s *ItemService) GetItems() ([]*item.Item, error) {
	return s.repo.GetItems()
}

func (s *ItemService) GetItem(id string) (*item.Item, error) {
	_, err := uuid.Parse(id)
	if err != nil {
		wbzlog.Logger.Warn().Err(err).Msg("invalid UUID format")
		return nil, fmt.Errorf("invalid UUID format: %w", err)
	}
	return s.repo.GetItem(id)
}

func (s *ItemService) PutItem(id string, name string, count int, price float64, userID string, login string) (*item.Item, error) {
	_, err := uuid.Parse(id)
	if err != nil {
		wbzlog.Logger.Warn().Err(err).Msg("invalid UUID format")
		return nil, fmt.Errorf("invalid UUID format: %w", err)
	}

	err = s.isNameValid(name)
	if err != nil {
		wbzlog.Logger.Warn().Err(err).Msg("item name is invalid")
		return nil, err
	}

	item, err := s.repo.GetItem(id)
	if err != nil {
		return nil, err
	}

	err = item.ChangeItem(name, count, price)
	if err != nil {
		wbzlog.Logger.Warn().Err(err).Msg("cant change item")
		return nil, err
	}

	err = s.repo.PutItem(item, userID, login)
	if err != nil {
		return nil, err
	}

	return item, nil
}

func (s *ItemService) DeleteItem(id string, userID string, login string) error {
	_, err := uuid.Parse(id)
	if err != nil {
		wbzlog.Logger.Warn().Err(err).Msg("invalid UUID format")
		return fmt.Errorf("invalid UUID format: %w", err)
	}
	return s.repo.DeleteItem(id, userID, login)
}

func (s *ItemService) isNameValid(name string) error {
	if name == "" || utf8.RuneCountInString(name) < s.cfg.ItemConfig.NameMinLength || utf8.RuneCountInString(name) > s.cfg.ItemConfig.NameMaxLegth {
		return fmt.Errorf("name cant be empty, should be bigger than %d, and smaller than %d", s.cfg.ItemConfig.NameMinLength, s.cfg.ItemConfig.NameMaxLegth)
	}
	return nil
}
