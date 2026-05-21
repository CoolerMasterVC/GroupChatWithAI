package handlers

import (
	"bytes"
	"encoding/json"
	"io"
	"log"
	"math/rand"
	"net/http"
	"os"
	"strconv"
	"time"

	"transport/internal/kafka"
	"transport/internal/models"
	"transport/internal/storage"
)

var lossProb = getLossProb()

func getLossProb() int {
	v := os.Getenv("LOSS_PROBABILITY")
	if v == "" {
		return 1
	}
	p, err := strconv.Atoi(v)
	if err != nil {
		return 1
	}
	return p
}

// sendErrorToApp немедленно отправляет ошибку на прикладной уровень
func sendErrorToApp(username string, sendTime time.Time) {
	appReceiveURL := os.Getenv("APP_RECEIVE_URL")
	if appReceiveURL == "" {
		appReceiveURL = "http://localhost:3000/receive"
	}

	errMsg := models.ReceiveRequest{
		Username: username,
		Text:     "",
		SendTime: sendTime,
		Error:    "lost",
	}
	data, _ := json.Marshal(errMsg)
	resp, err := http.Post(appReceiveURL, "application/json", bytes.NewBuffer(data))
	if err != nil {
		log.Printf("[ERROR] failed to send error to app: %v", err)
		return
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		log.Printf("[ERROR] app returned %d for error notification", resp.StatusCode)
	}
	log.Printf("[ERROR] sent 'lost' for %s (agent generation failed)", sendTime)
}

// HandleTransfer – обратите внимание: store передаётся через замыкание в main.go,
// здесь мы получаем его как аргумент. Сигнатура нестандартная, поэтому регистрация в main через адаптер.
// Для удобства используем замыкание, см. main.go.
func HandleTransfer(store *storage.Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		body, err := io.ReadAll(r.Body)
		if err != nil {
			writeError(w, http.StatusBadRequest, "cannot read body")
			return
		}
		defer r.Body.Close()

		var seg models.TransferRequest
		if err := json.Unmarshal(body, &seg); err != nil {
			writeError(w, http.StatusBadRequest, "invalid json")
			return
		}

		// Специальный случай: segment_number = 0 означает ошибку генерации у агента
		if seg.SegmentNumber == 0 {
			log.Printf("[TRANSFER] agent reported generation error for %s (segment 0)", seg.SendTime)
			// Немедленно уведомляем прикладной уровень об ошибке
			sendErrorToApp(seg.Username, seg.SendTime)
			// Отвечаем агенту успехом, чтобы он не повторял
			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(models.SuccessResponse{Status: "error acknowledged"})
			return
		}

		// Валидация номера сегмента (должен быть >=1 и <= total)
		if seg.SegmentNumber < 1 || seg.SegmentNumber > seg.TotalSegments {
			log.Printf("[TRANSFER] invalid segment number %d (total %d), returning 400", seg.SegmentNumber, seg.TotalSegments)
			writeError(w, http.StatusBadRequest, "invalid segment number")
			return
		}

		// Симуляция потери сегмента (только для нормальных сегментов)
		if rand.Intn(100) < lossProb {
			log.Printf("[TRANSFER] lost segment %d/%d for %s", seg.SegmentNumber, seg.TotalSegments, seg.SendTime)
			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(models.SuccessResponse{Status: "accepted (lost)"})
			return
		}

		// Отправка в Kafka (только для нормальных сегментов)
		if err := kafka.ProduceSegment(seg); err != nil {
			log.Printf("[TRANSFER] kafka produce error: %v", err)
			writeError(w, http.StatusInternalServerError, "kafka error")
			return
		}

		log.Printf("[TRANSFER] produced segment %d/%d for %s", seg.SegmentNumber, seg.TotalSegments, seg.SendTime)
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(models.SuccessResponse{Status: "accepted"})
	}
}
