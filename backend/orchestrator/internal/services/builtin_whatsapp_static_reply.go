package services

import (
	"fmt"
	wapp "github.com/Davanesh/auto-orchestrator/internal/executors"
)

func init() {
	RegisterExecutor("whatsapp_static_reply", &WhatsAppStaticReplyExecutor{})
}

type WhatsAppStaticReplyExecutor struct{}

// STATIC REPLY NODE:
// 1) Takes incoming message from input OR output
// 2) Applies regex + template
// 3) Saves formatted response into output
func (e *WhatsAppStaticReplyExecutor) Execute(n *ExecNode, g *ExecGraph) (string, error) {
	n.Status = "running"

	// PATCH: Determine input (try input → fallback to output)
	input := ""
	if v, ok := n.Data["input"]; ok && v != nil {
		input = fmt.Sprintf("%v", v)
	}
	if input == "" {
		if v, ok := n.Data["output"]; ok && v != nil {
			input = fmt.Sprintf("%v", v)
		}
	}

	regex := fmt.Sprintf("%v", n.Data["match_regex"])
	template := fmt.Sprintf("%v", n.Data["reply_template"])

	out, err := wapp.BuildStaticReply(regex, template, input)
	if err != nil {
		n.Status = "failed"
		return "", err
	}

	// PATCH: Save final message to output
	n.Data["output"] = out
	n.Status = "done"

	return "", nil
}
