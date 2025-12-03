package history

import (
	"warehousecontrol/internal/domain/history"

	"encoding/csv"
	"encoding/json"
	"fmt"
	"io"
	"time"

	"github.com/google/uuid"
	wbzlog "github.com/wb-go/wbf/zlog"
)

type HistoryService struct {
	repo HistoryStorageProvider
}

type HistoryStorageProvider interface {
	GetItemsHistory(id string, from, to time.Time, action string, login string) ([]*history.History, error)
}

func NewHistoryService(repo HistoryStorageProvider) *HistoryService {
	return &HistoryService{repo: repo}
}

func (s *HistoryService) GetItems(id string, from, to time.Time, action string, login string) ([]*history.History, error) {
	if !from.IsZero() && !to.IsZero() && from.After(to) {
		err := fmt.Errorf("'from' date cannot be after 'to'")
		wbzlog.Logger.Warn().Err(err).Msg("invalid date range in history request")
		return nil, err
	}
	if action != "" && action != "created" && action != "updated" && action != "deleted" {
		err := fmt.Errorf("invalid action filter")
		wbzlog.Logger.Warn().Err(err).Msg("invalid action filter in history request")
		return nil, err
	}
	if login != "" && len(login) < 3 {
		err := fmt.Errorf("login filter must be at least 3 characters long")
		wbzlog.Logger.Warn().Err(err).Msg("invalid login filter in history request")
		return nil, err
	}
	if id != "" {
		_, err := uuid.Parse(id)
		if err != nil {
			wbzlog.Logger.Warn().Err(err).Msg("invalid UUID format in history request")
			return nil, fmt.Errorf("invalid UUID format: %w", err)
		}
	}
	return s.repo.GetItemsHistory(id, from, to, action, login)
}

func (s *HistoryService) GetItemsCSV(id string, from, to time.Time, action string, login string, output io.Writer) error {
	histories, err := s.GetItems(id, from, to, action, login)
	if err != nil {
		wbzlog.Logger.Error().Err(err).Msg("repo get history error")
		return err
	}

	writer := csv.NewWriter(output)
	defer writer.Flush()

	headers := []string{"ID", "ItemID", "Action", "ChangedBy", "ChangedByLogin", "ChangedAt", "OldItemSnapshot", "NewItemSnapshot", "ItemDiff"}
	if err := writer.Write(headers); err != nil {
		wbzlog.Logger.Error().Err(err).Msg("Error writing CSV headers")
		return err
	}

	for _, group := range histories {
		itemDiffJSON, err := json.Marshal(group.ItemDiff)
		if err != nil {
			wbzlog.Logger.Error().Err(err).Msg("Error marshalling ItemDiff to JSON")
			return err
		}
		row := []string{
			group.ID.String(),
			group.ItemID.String(),
			group.Action,
			group.ChangedBy.String(),
			group.ChangedByLogin,
			group.ChangedAt.Format(time.RFC3339),
			fmt.Sprintf("%v", group.OldItemSnapshot),
			fmt.Sprintf("%v", group.NewItemSnapshot),
			string(itemDiffJSON),
		}
		if err := writer.Write(row); err != nil {
			wbzlog.Logger.Error().Err(err).Msg("Error writing CSV row")
			return err
		}
	}

	wbzlog.Logger.Info().Msg("CSV report generation completed")
	return nil
}
