package handlers

import (
	"encoding/json"
	"io"
	"log"
	"math/rand"
	"net/http"
	"os"
	"strconv"

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

		// потеря сегмента
		if rand.Intn(100) < lossProb {
			log.Printf("[TRANSFER] lost segment: %+v", seg)
			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(models.SuccessResponse{Status: "accepted (lost)"})
			return
		}

		// отправка в Kafka
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
