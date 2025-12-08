package services

import (
	"log"
	"strconv"

	wapp "github.com/Davanesh/auto-orchestrator/internal/executors"
)

// Register wait node executor
func init() {
	RegisterExecutor("whatsapp_wait", &WhatsAppWaitExecutor{})
}

type WhatsAppWaitExecutor struct{}

func (e *WhatsAppWaitExecutor) Execute(n *ExecNode, g *ExecGraph) (string, error) {
	n.Status = "running"

	// get timeoutSeconds from node data (if present)
	timeout := 0
	if v, ok := n.Data["timeoutSeconds"]; ok {
		switch t := v.(type) {
		case float64:
			timeout = int(t)
		case string:
			if iv, err := strconv.Atoi(t); err == nil {
				timeout = iv
			}
		case int:
			timeout = t
		}
	}

	// Use graph.RunID and node.ID to register waiter
	runID := g.RunID
	nodeID := n.ID

	log.Printf("⏳ Waiting for WhatsApp message (run=%s node=%s timeout=%d)", runID, nodeID, timeout)

	msg, err := wapp.WaitForWhatsAppMessage(runID, nodeID, timeout)
	if err != nil {
		n.Status = "failed"
		return "", err
	}

	// store incoming message into node.Data["output"] or ["input"] for next node
	n.Data["output"] = msg
	n.Status = "done"

	// No direct override — engine will use n.Next[0]
	return "", nil
}
