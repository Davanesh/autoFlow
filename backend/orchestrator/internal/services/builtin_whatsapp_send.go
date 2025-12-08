package services

import (
	"errors"
	"fmt"

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
		n.Status = "failed"
		return "", errors.New("missing 'to' in whatsapp_send node")
	}

	// GET message from previous node
	body := ""
	if v, exists := n.Data["input"]; exists {
		body = fmt.Sprintf("%v", v)
	}

	if body == "" {
		n.Status = "failed"
		return "", errors.New("whatsapp_send has empty body")
	}

	// Send message
	_, err := wapp.ExecuteWhatsAppSendNode(to, "static", "", "${body}", body)
	//                   regexPattern ↑   template ↑       input ↑
	if err != nil {
		n.Status = "failed"
		return "", err
	}

	n.Data["output"] = "Message sent"
	n.Status = "done"
	return "", nil
}
