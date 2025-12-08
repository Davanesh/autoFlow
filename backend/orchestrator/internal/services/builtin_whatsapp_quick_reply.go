package services

import (
	"errors"
	"fmt"
	"log"
	"strings"

	wapp "github.com/Davanesh/auto-orchestrator/internal/executors"
)

func init() {
	RegisterExecutor("whatsapp_quick_reply", &WhatsAppQuickReplyExecutor{})
}

type WhatsAppQuickReplyExecutor struct{}

func (e *WhatsAppQuickReplyExecutor) Execute(n *ExecNode, g *ExecGraph) (string, error) {
	n.Status = "running"

	// options expected as map[string]interface{} : "1": "nodeA", "2": "nodeB"
	rawOptions, ok := n.Data["options"]
	if !ok {
		return "", errors.New("quick reply missing options")
	}
	opts, ok := rawOptions.(map[string]interface{})
	if !ok {
		return "", errors.New("options must be object")
	}

	fallback := fmt.Sprintf("%v", n.Data["invalid_message"])
	if fallback == "" {
		fallback = "Invalid option. Try again."
	}

	runID := g.RunID
	nodeID := n.ID

	// wait
	msg, err := wapp.WaitForWhatsAppMessage(runID, nodeID, 0)
	if err != nil {
		n.Status = "failed"
		return "", err
	}
	msgNorm := strings.TrimSpace(strings.ToLower(msg))

	// match
	for k, v := range opts {
		if strings.ToLower(strings.TrimSpace(k)) == msgNorm {
			// return target node id (string)
			target := fmt.Sprintf("%v", v)
			n.Status = "done"
			return target, nil
		}
	}

	// not matched -> send fallback and re-run the same node (return nodeID)
	to := fmt.Sprintf("%v", n.Data["from"])
	if to != "" {
		_ = wapp.SendWhatsAppMessage(to, fallback)
	} else {
		// try to get 'sender' from g or node data
	}

	log.Printf("⚠ quick-reply invalid input: %s", msg)
	return n.ID, nil
}
