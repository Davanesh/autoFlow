package services

import (
	"errors"
	"fmt"
	"log"
	"strings"

	wapp "github.com/Davanesh/auto-orchestrator/internal/executors"
)

func init() {
	RegisterExecutor("whatsapp_keyword_router", &WhatsAppKeywordRouterExecutor{})
}

type WhatsAppKeywordRouterExecutor struct{}

func (e *WhatsAppKeywordRouterExecutor) Execute(n *ExecNode, g *ExecGraph) (string, error) {
	n.Status = "running"

	rawKeywords, ok := n.Data["keywords"]
	if !ok {
		return "", errors.New("keyword router missing 'keywords'")
	}
	kwMap, ok := rawKeywords.(map[string]interface{})
	if !ok {
		return "", errors.New("'keywords' must be an object map")
	}

	fallback := fmt.Sprintf("%v", n.Data["fallback_message"])
	if fallback == "" {
		fallback = "I didn't understand — try again."
	}

	runID := g.RunID
	nodeID := n.ID

	msg, err := wapp.WaitForWhatsAppMessage(runID, nodeID, 0)
	if err != nil {
		n.Status = "failed"
		return "", err
	}
	lmsg := strings.ToLower(strings.TrimSpace(msg))

	// iterate key groups
	for k, target := range kwMap {
		parts := strings.Split(k, ",")
		for _, p := range parts {
			p = strings.ToLower(strings.TrimSpace(p))
			if p == "" {
				continue
			}
			if lmsg == p || strings.Contains(lmsg, p) {
				n.Status = "done"
				return fmt.Sprintf("%v", target), nil
			}
		}
	}

	// no match
	to := fmt.Sprintf("%v", n.Data["from"])
	if to != "" {
		_ = wapp.SendWhatsAppMessage(to, fallback)
	} else {
		log.Printf("⚠ keyword_router no sender to send fallback")
	}
	return n.ID, nil
}
