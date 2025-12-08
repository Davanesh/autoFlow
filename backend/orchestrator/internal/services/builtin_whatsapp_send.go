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

// WHATSAPP SEND NODE:
// mode = "static" → use template + regex
// mode = "ai" → generate using internal AI node
// Finally sends via Twilio and stores final output.
func (e *WhatsAppSendExecutor) Execute(n *ExecNode, g *ExecGraph) (string, error) {
	n.Status = "running"

	//---------------------------------------------
	// 1) Read destination number
	//---------------------------------------------
	to := fmt.Sprintf("%v", n.Data["to"])
	if to == "" {
		n.Status = "failed"
		return "", errors.New("missing 'to' in whatsapp_send node")
	}

	//---------------------------------------------
	// 2) Determine input message
	//---------------------------------------------
	input := ""
	if v, ok := n.Data["input"]; ok && v != nil {
		input = fmt.Sprintf("%v", v)
	}
	if input == "" {
		if v, ok := n.Data["output"]; ok && v != nil {
			input = fmt.Sprintf("%v", v)
		}
	}

	if input == "" {
		input = "(empty input)"
	}

	//---------------------------------------------
	// 3) Mode selection: "static" or "ai"
	//---------------------------------------------
	mode := "static"
	if v, ok := n.Data["mode"]; ok && v != nil {
		mode = fmt.Sprintf("%v", v)
	}

	regex := fmt.Sprintf("%v", n.Data["match_regex"])
	template := fmt.Sprintf("%v", n.Data["reply_template"])

	var finalOut string
	var err error

	if mode == "static" {
		// Use static reply builder
		finalOut, err = wapp.BuildStaticReply(regex, template, input)
		if err != nil {
			n.Status = "failed"
			return "", err
		}
	} else if mode == "ai" {
		// AI GENERATED MODE
		finalOut, err = wapp.ExecuteWhatsAppSendNode(to, "ai", "", template, input)
		if err != nil {
			n.Status = "failed"
			return "", err
		}
		// Already sent in ExecuteWhatsAppSendNode, so skip sending again
		n.Data["output"] = finalOut
		n.Status = "done"
		return "", nil
	} else {
		n.Status = "failed"
		return "", errors.New("unknown mode in whatsapp_send node")
	}

	//---------------------------------------------
	// 4) Send via Twilio
	//---------------------------------------------
	err = wapp.SendWhatsAppMessage(to, finalOut)
	if err != nil {
		n.Status = "failed"
		return "", err
	}

	//---------------------------------------------
	// 5) Save final output
	//---------------------------------------------
	n.Data["output"] = finalOut
	n.Status = "done"
	return "", nil
}
