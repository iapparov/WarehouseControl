package history

import (
	"github.com/google/uuid"
	"time"
	"warehousecontrol/internal/domain/item"
)

type History struct {
	ID              uuid.UUID     `json:"id"`
	ItemID          uuid.UUID     `json:"item_id"`
	Action          string        `json:"action"`
	ChangedBy       uuid.UUID     `json:"changed_by"`
	ChangedByLogin  string        `json:"changed_by_login"`
	ChangedAt       time.Time     `json:"changed_at"`
	OldItemSnapshot item.Item     `json:"old_item_snapshot"`
	NewItemSnapshot item.Item     `json:"new_item_snapshot"`
	ItemDiff        item.ItemDiff `json:"item_diff"`
}
