package handlers

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"transport/internal/models"
)

// HandleSend обрабатывает POST /send
// @Summary Передать сообщение от прикладного уровня на агентный уровень
// @Description Принимает полное сообщение от прикладного уровня и отправляет его целиком на агентный уровень (без разбиения на сегменты)
// @Tags transport
// @Accept json
// @Produce json
// @Param request body models.SendRequest true "Полное сообщение от пользователя"
// @Success 200 {object} models.SuccessResponse "Сообщение успешно передано"
// @Failure 400 {object} models.ErrorResponse "Неверный формат запроса"
// @Failure 502 {object} models.ErrorResponse "Ошибка при передаче на агентный уровень"
// @Router /send [post]
func HandleSend(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(models.ErrorResponse{Error: "cannot read body"})
		return
	}
	defer r.Body.Close()

	var req models.SendRequest
	if err := json.Unmarshal(body, &req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(models.ErrorResponse{Error: "invalid json"})
		return
	}

	log.Printf("[SEND] received from %s: %s", req.Username, req.Data)

	// TODO: реальная отправка на агентный уровень
	// if err := sendToAgent(req); err != nil {
	//     w.WriteHeader(http.StatusBadGateway)
	//     json.NewEncoder(w).Encode(models.ErrorResponse{Error: "agent unavailable"})
	//     return
	// }

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(models.SuccessResponse{Status: "accepted"})
}
