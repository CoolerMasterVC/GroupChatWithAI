package handlers

import (
	"bytes"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"os"

	"transport/internal/models"
)

var agentURL = getEnv("AGENT_SUMMARY_URL", "http://localhost:8001/summary")

func HandleSend(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		writeError(w, http.StatusBadRequest, "cannot read body")
		return
	}
	defer r.Body.Close()

	var req models.SendRequest
	if err := json.Unmarshal(body, &req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid json")
		return
	}

	log.Printf("[SEND] forwarding to agent: %+v", req)

	agentReq, _ := json.Marshal(req)
	resp, err := http.Post(agentURL, "application/json", bytes.NewBuffer(agentReq))
	if err != nil {
		log.Printf("[SEND] agent call failed: %v", err)
		writeError(w, http.StatusBadGateway, "agent unavailable")
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		log.Printf("[SEND] agent returned %d", resp.StatusCode)
		writeError(w, http.StatusBadGateway, "agent error")
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(models.SuccessResponse{Status: "forwarded to agent"})
}

func writeError(w http.ResponseWriter, code int, msg string) {
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(models.ErrorResponse{Error: msg})
}

func getEnv(key, def string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return def
}
