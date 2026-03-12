package handlers

import (
	"encoding/json"
	"io"
	"net/http"

	"transport/internal/models"
)

// ReceiveSegment обработчик POST /segment
// @Summary Принять сегмент от агента
// @Description Принимает один сегмент сообщения. В реальной реализации здесь должна быть отправка в Kafka, потеря с вероятностью R% и т.д. В данной заглушке просто возвращается 200 OK.
// @Tags transport
// @Accept json
// @Produce json
// @Param segment body models.Segment true "Данные сегмента"
// @Success 200 {object} map[string]string "Сегмент принят (заглушка)"
// @Failure 400 {object} map[string]string "Неверный формат запроса"
// @Router /segment [post]
func ReceiveSegment(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, `{"error":"cannot read body"}`, http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	var seg models.Segment
	if err := json.Unmarshal(body, &seg); err != nil {
		http.Error(w, `{"error":"invalid json"}`, http.StatusBadRequest)
		return
	}

	// Здесь в реальном коде была бы логика:
	// - с вероятностью R% (номер варианта) "теряем" сегмент
	// - иначе отправляем в Kafka
	// - сегменты читаются из Kafka, накапливаются, собираются и отправляются на прикладной уровень

	// Заглушка: просто логируем полученный сегмент
	// (можно заменить на fmt.Printf или использование логгера)

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"status": "accepted (stub)"})
}
