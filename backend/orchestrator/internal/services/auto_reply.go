package services

import (
	"log"
	"time"

	wapp "github.com/Davanesh/auto-orchestrator/internal/executors"
)

func AutoReplyFlow(message string, from string) {
	log.Println("Auto-reply triggered for:", from)

	// build mini flow: input -> reply -> send

	output := "You said: " + message

	// delay optional
	time.Sleep(300 * time.Millisecond)

	// send reply
	err := wapp.SendWhatsAppMessage(from, output)
	if err != nil {
		log.Println("Auto-reply failed:", err)
		return
	}

	log.Println("Auto-replied:", output)
}
