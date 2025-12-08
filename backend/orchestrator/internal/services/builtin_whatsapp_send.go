package services

import (
	"errors"
	"fmt"
	"log"

	wapp "github.com/Davanesh/auto-orchestrator/internal/executors"
)

func init() {
	RegisterExecutor("whatsapp_send", &WhatsAppSendExecutor{})
}

type WhatsAppSendExecutor struct{}

func (e *WhatsAppSendExecutor) Execute(n *ExecNode, g *ExecGraph) (string, error) {
	n.Status = "running"

	to := fmt.Sprintf("%v", n.Data["to"])
	if to == "" {
		// try to pull 'from' or 'recipient' keys
		to = fmt.Sprintf("%v", n.Data["recipient"])
	}
	if to == "" {
		n.Status = "failed"
		return "", errors.New("whatsapp_send missing 'to'")
	}

	// prefer explicit output > input
	body := fmt.Sprintf("%v", n.Data["output"])
	if body == "" {
		body = fmt.Sprintf("%v", n.Data["input"])
	}
	if body == "" {
		body = "(empty message)"
	}

	// normalize to format venom expects (we use SendWhatsAppMessage wrapper)
	// venom expects plain number (no prefix) or will handle +, whatsapp:
	err := wapp.SendWhatsAppMessage(to, body)
	if err != nil {
		n.Status = "failed"
		return "", fmt.Errorf("send failed: %w", err)
	}

	log.Printf("📤 WhatsApp sent to=%s body=%s", to, body)

	n.Status = "done"
	return "", nil
}
