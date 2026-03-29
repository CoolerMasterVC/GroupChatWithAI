package handlers

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"transport/internal/models"
)

// HandleTransfer обрабатывает POST /transfer
// @Summary Принять сегмент от агентного уровня
// @Description Принимает сегмент от агента и сохраняет для последующей сборки
// @Tags transport
// @Accept json
// @Produce json
// @Param segment body models.TransferRequest true "Сегмент сообщения"
// @Success 200 {object} models.SuccessResponse "Сегмент принят"
// @Failure 400 {object} models.ErrorResponse "Неверный формат запроса"
// @Router /transfer [post]
func HandleTransfer(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(models.ErrorResponse{Error: "cannot read body"})
		return
	}
	defer r.Body.Close()

	var seg models.TransferRequest
	if err := json.Unmarshal(body, &seg); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(models.ErrorResponse{Error: "invalid json"})
		return
	}

	log.Printf("[TRANSFER] segment %d/%d from %s", seg.SegmentNumber, seg.TotalSegments, seg.Username)

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(models.SuccessResponse{Status: "accepted"})
}
