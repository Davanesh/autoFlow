package services

import (
	"fmt"
	"strings"
	"errors"

	wapp "github.com/Davanesh/auto-orchestrator/internal/executors"
)

func init() {
	RegisterExecutor("whatsapp_quick_reply", &WhatsAppQuickReplyExecutor{})
}

type WhatsAppQuickReplyExecutor struct{}

func (e *WhatsAppQuickReplyExecutor) Execute(n *ExecNode, g *ExecGraph) (string, error) {
	n.Status = "running"

	runID := g.RunID
	nodeID := n.ID

	//---------------------------------------
	// 1) Load options (expected replies)
	//---------------------------------------
	rawOptions, ok := n.Data["options"]
	if !ok {
		return "", errors.New("quick reply node missing 'options'")
	}

	optionsMap, ok := rawOptions.(map[string]interface{})
	if !ok {
		return "", errors.New("'options' must be a map")
	}

	//---------------------------------------
	// 2) Load fallback message
	//---------------------------------------
	invalidMessage := fmt.Sprintf("%v", n.Data["invalid_message"])
	if invalidMessage == "" {
		invalidMessage = "Invalid reply. Please try again."
	}

	//---------------------------------------
	// 3) Wait for actual WhatsApp message
	//---------------------------------------
	msg, err := wapp.WaitForWhatsAppMessage(runID, nodeID, 0)
	if err != nil {
		n.Status = "failed"
		return "", err
	}

	msg = strings.TrimSpace(strings.ToLower(msg))

	//---------------------------------------
	// 4) Check if user reply matches a key
	//---------------------------------------
	for key, nextName := range optionsMap {
		if strings.ToLower(key) == msg {
			n.Status = "done"
			return fmt.Sprintf("%v", nextName), nil
		}
	}

	//---------------------------------------
	// 5) If invalid → send fallback
	//---------------------------------------
	wapp.SendWhatsAppMessage(n.Data["from"].(string), invalidMessage)

	// Wait again → keep looping until user gives correct input
	return n.ID, nil
}
