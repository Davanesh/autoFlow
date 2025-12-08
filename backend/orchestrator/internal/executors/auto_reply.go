package executors

import (
	"log"
	"strings"
)

// AutoReplySimple echoes incoming message back
func AutoReplySimple(message string, from string) {
	log.Println("🤖 Auto-reply triggered for:", from)

	msg := strings.TrimSpace(message)
	if msg == "" {
		msg = "(empty message)"
	}

	// quick normalization: if from like "whatsapp:+919..." keep as is
	to := from
	to = strings.ReplaceAll(to, "whatsapp:", "")
	to = strings.ReplaceAll(to, "+", "")

	err := SendWhatsAppMessage(to, msg)
	if err != nil {
		log.Println("❌ AutoReply send error:", err)
		return
	}
	log.Println("✅ Auto-replied ->", msg)
}
