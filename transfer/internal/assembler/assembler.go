package assembler

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"transport/internal/models"
	"transport/internal/storage"
)

var (
	appReceiveURL  string
	intervalSec    int
	lossTimeoutSec int
)

func init() {
	appReceiveURL = os.Getenv("APP_RECEIVE_URL")
	if appReceiveURL == "" {
		appReceiveURL = "http://localhost:3000/receive"
	}
	intervalSec, _ = strconv.Atoi(os.Getenv("ASSEMBLY_INTERVAL_SEC"))
	if intervalSec <= 0 {
		intervalSec = 2
	}
	lossTimeoutSec, _ = strconv.Atoi(os.Getenv("LOSS_TIMEOUT_SEC"))
	if lossTimeoutSec <= 0 {
		lossTimeoutSec = 4
	}
}

func StartAssembler(store *storage.Storage) {
	ticker := time.NewTicker(time.Duration(intervalSec) * time.Second)
	go func() {
		for range ticker.C {
			scanAndSend(store)
		}
	}()
}

func scanAndSend(store *storage.Storage) {
	keys := store.GetAllKeys()
	now := time.Now().UTC()

	for _, sendTime := range keys {
		info, ok := store.Peek(sendTime)
		if !ok {
			continue
		}

		if info.Received == info.Total {
			fullText := ""
			for _, seg := range info.Segments {
				fullText += seg
			}
			// ВЫВОДИМ ПОЛНОЕ СООБЩЕНИЕ В КОНСОЛЬ
			log.Printf("[ASSEMBLER] FULL MESSAGE: %s", fullText)

			sendToApp(models.ReceiveRequest{
				Username: info.Username,
				Text:     fullText,
				SendTime: sendTime,
				Error:    "",
			})
			store.GetAndDelete(sendTime)
			log.Printf("[ASSEMBLER] sent full message for %s", sendTime)
			continue
		}

		if now.Sub(info.Last) > time.Duration(lossTimeoutSec)*time.Second {
			log.Printf("[ASSEMBLER] lost segments for %s", sendTime) // сообщение с ошибкой
			sendToApp(models.ReceiveRequest{
				Username: info.Username,
				Text:     "",
				SendTime: sendTime,
				Error:    "lost",
			})
			store.GetAndDelete(sendTime)
			log.Printf("[ASSEMBLER] sent error (lost) for %s", sendTime)
		}
	}
}

func sendToApp(req models.ReceiveRequest) {
	data, _ := json.Marshal(req)
	resp, err := http.Post(appReceiveURL, "application/json", bytes.NewBuffer(data))
	if err != nil {
		log.Printf("[ASSEMBLER] failed to send: %v", err)
		return
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		log.Printf("[ASSEMBLER] app responded with %d", resp.StatusCode)
	}
}
