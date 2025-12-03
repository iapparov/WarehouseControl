package handlers

import (
	"io"
	"net/http"
	"time"

	"warehousecontrol/internal/domain/history"

	wbgin "github.com/wb-go/wbf/ginext"
)

type HistoryHandler struct {
	Service HistoryIFace
}

type HistoryIFace interface {
	GetItems(id string, from, to time.Time, action string, login string) ([]*history.History, error)
	GetItemsCSV(id string, from, to time.Time, action string, login string, output io.Writer) error
}

func NewHistoryHandler(service HistoryIFace) *HistoryHandler {
	return &HistoryHandler{
		Service: service,
	}
}

// GetItems
// @Summary List history
// @Description Get item change history filtered by date range and optional filters
// @Tags history
// @Produce json
// @Param from query string true "From date (YYYY-MM-DD)"
// @Param to query string true "To date (YYYY-MM-DD)"
// @Param id query string false "Item UUID"
// @Param action query string false "created|updated|deleted"
// @Param login query string false "login substring"
// @Success 200 {array} history.History
// @Failure 400 {object} map[string]string
// @Security BearerAuth
// @Router /api/history [get]
func (h *HistoryHandler) GetItems(ctx *wbgin.Context) {
	from := ctx.Query("from")
	to := ctx.Query("to")
	action := ctx.Query("action")
	login := ctx.Query("login")
	id := ctx.Query("id")

	layout := "2006-01-02"
	fromParsed, err := time.ParseInLocation(layout, from, time.Local)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, wbgin.H{"error": "invalid from date format"})
		return
	}
	toParsed, err := time.ParseInLocation(layout, to, time.Local)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, wbgin.H{"error": "invalid to date format"})
		return
	}

	histories, err := h.Service.GetItems(id, fromParsed, toParsed, action, login)

	if err != nil {
		ctx.JSON(http.StatusBadRequest, wbgin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, histories)
}

// GetItemsCSV
// @Summary Export history CSV
// @Description Download item change history as CSV with same filters
// @Tags history
// @Produce text/csv
// @Param from query string true "From date (YYYY-MM-DD)"
// @Param to query string true "To date (YYYY-MM-DD)"
// @Param id query string false "Item UUID"
// @Param action query string false "created|updated|deleted"
// @Param login query string false "login substring"
// @Success 200 "CSV file"
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Security BearerAuth
// @Router /api/history/csv [get]
func (h *HistoryHandler) GetItemsCSV(ctx *wbgin.Context) {
	from := ctx.Query("from")
	to := ctx.Query("to")
	action := ctx.Query("action")
	login := ctx.Query("login")
	id := ctx.Query("id")

	layout := "2006-01-02"
	fromParsed, err := time.ParseInLocation(layout, from, time.Local)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, wbgin.H{"error": "invalid from date format"})
		return
	}
	toParsed, err := time.ParseInLocation(layout, to, time.Local)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, wbgin.H{"error": "invalid to date format"})
		return
	}

	ctx.Writer.Header().Set("Content-Disposition", "attachment; filename=transactions.csv")
	ctx.Writer.Header().Set("Content-Type", "text/csv")
	err = h.Service.GetItemsCSV(id, fromParsed, toParsed, action, login, ctx.Writer)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, wbgin.H{"error": err.Error()})
		return
	}
}
