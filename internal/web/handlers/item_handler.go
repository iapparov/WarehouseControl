package handlers

import (
	"net/http"

	"warehousecontrol/internal/domain/item"
	"warehousecontrol/internal/web/dto"

	wbgin "github.com/wb-go/wbf/ginext"
)

type ItemHandler struct {
	Service ItemIFace
}

type ItemIFace interface {
	Create(name string, count int, price float64, userID string, login string) (*item.Item, error)
	GetItems() ([]*item.Item, error)
	GetItem(id string) (*item.Item, error)
	PutItem(id string, name string, count int, price float64, userID string, login string) (*item.Item, error)
	DeleteItem(id string, userID string, login string) error
}

func NewItemHandler(service ItemIFace) *ItemHandler {
	return &ItemHandler{
		Service: service,
	}
}

// CreateItem
// @Summary Create item
// @Description Create a new inventory item (admin only)
// @Tags items
// @Accept json
// @Produce json
// @Param body body dto.ItemCreateRequest true "Item payload"
// @Success 200 {object} item.Item
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Security BearerAuth
// @Router /api/items [post]
func (h *ItemHandler) CreateItem(ctx *wbgin.Context) {
	var req dto.ItemCreateRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, wbgin.H{"error": err.Error()})
		return
	}
	userID, ok := ctx.Get("userId")
	if !ok {
		ctx.JSON(http.StatusInternalServerError, wbgin.H{"error": "userId not found in context"})
		return
	}
	login, ok := ctx.Get("login")
	if !ok {
		ctx.JSON(http.StatusInternalServerError, wbgin.H{"error": "login not found in context"})
		return
	}
	item, err := h.Service.Create(req.Name, req.Count, req.Price, userID.(string), login.(string))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, wbgin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, item)
}

// GetItems
// @Summary List items
// @Description Get list of items
// @Tags items
// @Produce json
// @Success 200 {array} item.Item
// @Failure 500 {object} map[string]string
// @Security BearerAuth
// @Router /api/items [get]
func (h *ItemHandler) GetItems(ctx *wbgin.Context) {
	items, err := h.Service.GetItems()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, wbgin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, items)
}

// GetItem
// @Summary Get item
// @Description Get a single item by ID
// @Tags items
// @Produce json
// @Param id path string true "Item UUID"
// @Success 200 {object} item.Item
// @Failure 500 {object} map[string]string
// @Security BearerAuth
// @Router /api/items/{id} [get]
func (h *ItemHandler) GetItem(ctx *wbgin.Context) {
	id := ctx.Param("id")
	item, err := h.Service.GetItem(id)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, wbgin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, item)
}

// PutItem
// @Summary Update item
// @Description Update item fields (admin/manager)
// @Tags items
// @Accept json
// @Produce json
// @Param id path string true "Item UUID"
// @Param body body dto.ItemUpdateRequest true "Updated fields"
// @Success 200 {object} item.Item
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Security BearerAuth
// @Router /api/items/{id} [put]
func (h *ItemHandler) PutItem(ctx *wbgin.Context) {
	id := ctx.Param("id")
	var req dto.ItemUpdateRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, wbgin.H{"error": err.Error()})
		return
	}
	userID, ok := ctx.Get("userId")
	if !ok {
		ctx.JSON(http.StatusInternalServerError, wbgin.H{"error": "userId not found in context"})
		return
	}
	login, ok := ctx.Get("login")
	if !ok {
		ctx.JSON(http.StatusInternalServerError, wbgin.H{"error": "login not found in context"})
		return
	}
	item, err := h.Service.PutItem(id, req.Name, req.Count, req.Price, userID.(string), login.(string))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, wbgin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, item)
}

// DeleteItem
// @Summary Delete item
// @Description Delete item (admin only)
// @Tags items
// @Produce json
// @Param id path string true "Item UUID"
// @Success 200 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Security BearerAuth
// @Router /api/items/{id} [delete]
func (h *ItemHandler) DeleteItem(ctx *wbgin.Context) {
	id := ctx.Param("id")
	userID, ok := ctx.Get("userId")
	if !ok {
		ctx.JSON(http.StatusInternalServerError, wbgin.H{"error": "userId not found in context"})
		return
	}
	login, ok := ctx.Get("login")
	if !ok {
		ctx.JSON(http.StatusInternalServerError, wbgin.H{"error": "login not found in context"})
		return
	}
	err := h.Service.DeleteItem(id, userID.(string), login.(string))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, wbgin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, wbgin.H{"message": "item deleted"})
}
