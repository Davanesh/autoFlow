package executors

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"time"
)

// VENOM service URL
var VenomServiceURL = func() string {
	if v := os.Getenv("VENOM_SERVICE_URL"); v != "" {
		return v
	}
	return "http://localhost:5000"
}()

// SendWhatsAppMessage — ONLY Venom (NO Twilio)
func SendWhatsAppMessage(to, body string) error {
	client := &http.Client{Timeout: 15 * time.Second}

	// Normalise number
	to = strings.ReplaceAll(to, "whatsapp:", "")
	to = strings.TrimPrefix(to, "+")
	to = strings.TrimSpace(to)

	payload := map[string]interface{}{
		"to":      to,
		"message": body,
	}

	b, _ := json.Marshal(payload)

	req, _ := http.NewRequest("POST", VenomServiceURL+"/send", bytes.NewReader(b))
	req.Header.Set("Content-Type", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	data, _ := io.ReadAll(resp.Body)

	if resp.StatusCode >= 200 && resp.StatusCode < 300 {
		return nil
	}

	return fmt.Errorf("venom send error: status=%d body=%s", resp.StatusCode, string(data))
}
