package services

import (
	"errors"
	"fmt"
	"strings"

	wapp "github.com/Davanesh/auto-orchestrator/internal/executors"
)

func init() {
	RegisterExecutor("whatsapp_keyword_router", &WhatsAppKeywordRouterExecutor{})
}

type WhatsAppKeywordRouterExecutor struct{}

func (e *WhatsAppKeywordRouterExecutor) Execute(n *ExecNode, g *ExecGraph) (string, error) {
	n.Status = "running"

	runID := g.RunID
	nodeID := n.ID

	//--------------------------------------------
	// 1) Load keyword map
	//--------------------------------------------
	rawKeywords, ok := n.Data["keywords"]
	if !ok {
		return "", errors.New("keyword_router node missing 'keywords'")
	}

	keywordMap, ok := rawKeywords.(map[string]interface{})
	if !ok {
		return "", errors.New("'keywords' must be a map")
	}

	//--------------------------------------------
	// 2) Load fallback message
	//--------------------------------------------
	fallback := fmt.Sprintf("%v", n.Data["fallback_message"])
	if fallback == "" {
		fallback = "Invalid choice. Try again."
	}

	//--------------------------------------------
	// 3) Wait for incoming message
	//--------------------------------------------
	msg, err := wapp.WaitForWhatsAppMessage(runID, nodeID, 0)
	if err != nil {
		n.Status = "failed"
		return "", err
	}

	lowerMsg := strings.ToLower(strings.TrimSpace(msg))

	//--------------------------------------------
	// 4) Try to match keywords
	//--------------------------------------------
	for key, nextNode := range keywordMap {
		// key may contain comma-separated keywords
		group := strings.Split(key, ",")
		for _, word := range group {
			keyword := strings.ToLower(strings.TrimSpace(word))

			// match: exact OR contains
			if lowerMsg == keyword || strings.Contains(lowerMsg, keyword) {
				n.Status = "done"
				return fmt.Sprintf("%v", nextNode), nil
			}
		}
	}

	//--------------------------------------------
	// 5) No match → send fallback message
	//--------------------------------------------
	from := fmt.Sprintf("%v", n.Data["from"])
	if from != "" {
		wapp.SendWhatsAppMessage(from, fallback)
	}

	//--------------------------------------------
	// Loop again (re-run same node)
	//--------------------------------------------
	return n.ID, nil
}
