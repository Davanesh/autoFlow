package executors

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"regexp"
	"strings"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
)

//
// ────────────────────────────────────────────────────────────────
//    GLOBAL WAITERS (USED BY WORKFLOW ENGINE)
// ────────────────────────────────────────────────────────────────
//

var waiters = struct {
	m sync.Map // map[string]chan string
}{}

func waiterKey(runID, nodeID string) string {
	return runID + ":" + nodeID
}

func RegisterWaiter(runID, nodeID string, timeoutSeconds int) (chan string, error) {
	key := waiterKey(runID, nodeID)
	ch := make(chan string, 1)
	waiters.m.Store(key, ch)

	if timeoutSeconds > 0 {
		go func() {
			<-time.After(time.Duration(timeoutSeconds) * time.Second)
			if val, ok := waiters.m.Load(key); ok {
				if c, ok2 := val.(chan string); ok2 {
					select { case c <- "__TIMEOUT__": default: }
					waiters.m.Delete(key)
				}
			}
		}()
	}

	return ch, nil
}

func WaitForWhatsAppMessage(runID, nodeID string, timeoutSeconds int) (string, error) {
	ch, _ := RegisterWaiter(runID, nodeID, timeoutSeconds)
	msg := <-ch
	if msg == "__TIMEOUT__" {
		return "", fmt.Errorf("waiter timeout")
	}
	return msg, nil
}

func deliverMessageToWaiter(runID, nodeID, text string) bool {
	key := waiterKey(runID, nodeID)
	if val, ok := waiters.m.Load(key); ok {
		if ch, ok2 := val.(chan string); ok2 {
			select { case ch <- text: default: }
			waiters.m.Delete(key)
			return true
		}
	}
	return false
}

//
// ────────────────────────────────────────────────────────────────
//    PAYLOAD STRUCT FOR VENOM → GO
// ────────────────────────────────────────────────────────────────
//

type incomingPayload struct {
	From string `json:"From"`
	Body string `json:"Body"`
}

//
// ────────────────────────────────────────────────────────────────
//    MAIN WEBHOOK HANDLER
// ────────────────────────────────────────────────────────────────
//

func HandleWhatsAppWebhook(w http.ResponseWriter, r *http.Request) {
	var payload incomingPayload

	// Read raw body
	bodyBytes, _ := io.ReadAll(r.Body)
	r.Body.Close()

	// Try JSON parse
	if len(bodyBytes) > 0 {
		_ = json.Unmarshal(bodyBytes, &payload)
	}

	// Fallback: try form parsing
	if payload.Body == "" {
		_ = r.ParseForm()
		payload.Body = r.FormValue("Body")
		payload.From = r.FormValue("From")
	}

	if payload.From == "" {
		w.WriteHeader(400)
		w.Write([]byte("missing From"))
		return
	}

	log.Printf("📩 Venom incoming from=%s body=%s", payload.From, payload.Body)

	//
	// ─── 1) DIRECT run:xxx node:yyy ROUTING ─────────────────────
	//
	text := strings.ToLower(payload.Body)
	runRe := regexp.MustCompile(`run[:=]\s*([^\s]+)`)
	nodeRe := regexp.MustCompile(`node[:=]\s*([^\s]+)`)

	runID := ""
	nodeID := ""

	if m := runRe.FindStringSubmatch(text); len(m) > 1 {
		runID = m[1]
	}
	if m := nodeRe.FindStringSubmatch(text); len(m) > 1 {
		nodeID = m[1]
	}

	if runID != "" && nodeID != "" {
		if deliverMessageToWaiter(runID, nodeID, payload.Body) {
			w.WriteHeader(200)
			w.Write([]byte("Delivered"))
			return
		}
	}

	//
	// ─── 2) AUTO-WAKE: Send message to first active waiter ─────
	//
	delivered := false
	waiters.m.Range(func(k, v interface{}) bool {
		if ch, ok := v.(chan string); ok {
			select { case ch <- payload.Body: default: }
			waiters.m.Delete(k)
			delivered = true
			return false
		}
		return true
	})

	if delivered {
		w.WriteHeader(200)
		w.Write([]byte("Delivered (auto)"))
		return
	}

	//
	// ─── 3) AUTO-REPLY (ECHO MODE) ─────────────────────────────
	//
	go func() {
		SendWhatsAppMessage(payload.From, payload.Body)
	}()

	w.WriteHeader(200)
	w.Write([]byte("Auto replied"))
}

//
// ────────────────────────────────────────────────────────────────
//    GIN WRAPPER — REQUIRED BY main.go
// ────────────────────────────────────────────────────────────────
//

func HandleWhatsAppWebhookGin(c *gin.Context) {
	HandleWhatsAppWebhook(c.Writer, c.Request)
}

