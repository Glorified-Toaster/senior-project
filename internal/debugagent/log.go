package debugagent

import (
	"encoding/json"
	"os"
	"time"
)

const logPath = "/home/potato/Dev/senior-project/.cursor/debug-6333be.log"

// Log appends one NDJSON line for debug-mode analysis. Do not log secrets or PII.
func Log(sessionID, hypothesisID, runID, location, message string, data map[string]any) {
	payload := map[string]any{
		"sessionId":    sessionID,
		"hypothesisId": hypothesisID,
		"runId":        runID,
		"location":     location,
		"message":      message,
		"data":         data,
		"timestamp":    time.Now().UnixMilli(),
	}
	b, err := json.Marshal(payload)
	if err != nil {
		return
	}
	f, err := os.OpenFile(logPath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0o644)
	if err != nil {
		return
	}
	_, _ = f.Write(append(b, '\n'))
	_ = f.Close()
}
