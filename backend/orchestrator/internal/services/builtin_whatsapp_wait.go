package services

import (
	wapp "github.com/Davanesh/auto-orchestrator/internal/executors"
)

func init() {
	RegisterExecutor("whatsapp_wait", &WhatsAppWaitExecutor{})
}

type WhatsAppWaitExecutor struct{}

// WAIT NODE:
// 1) Waits for incoming WhatsApp message
// 2) Stores that message inside n.Data["output"]
// 3) Marks node done
func (e *WhatsAppWaitExecutor) Execute(n *ExecNode, g *ExecGraph) (string, error) {
	n.Status = "waiting"

	// Get timeout (0 = infinite)
	timeout := 0
	if v, ok := n.Data["timeoutSeconds"]; ok {
		timeout = int(v.(float64))
	}

	// WAIT FOR WHATSAPP MESSAGE
	msg, err := wapp.WaitForWhatsAppMessage(g.RunID, n.ID, timeout)
	if err != nil {
		n.Status = "failed"
		return "", err
	}

	// PATCH: save incoming message inside output
	n.Data["output"] = msg
	n.Status = "done"

	return "", nil
}
