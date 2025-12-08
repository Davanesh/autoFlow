package services

import (
	"fmt"
	"regexp"
	"strings"
)

func init() {
	RegisterExecutor("whatsapp_static_reply", &WhatsAppStaticReplyExecutor{})
}

type WhatsAppStaticReplyExecutor struct{}

func (e *WhatsAppStaticReplyExecutor) Execute(n *ExecNode, g *ExecGraph) (string, error) {
	n.Status = "running"

	// input can come from previous node in graph (n.Data["input"]) or from wait node output
	input := fmt.Sprintf("%v", n.Data["input"])
	if input == "" {
		input = fmt.Sprintf("%v", n.Data["body"])
	}
	if input == "" {
		input = fmt.Sprintf("%v", n.Data["output"])
	}

	regex := fmt.Sprintf("%v", n.Data["match_regex"])
	template := fmt.Sprintf("%v", n.Data["reply_template"])

	out := ""
	if strings.TrimSpace(regex) == "" {
		out = strings.ReplaceAll(template, "${body}", input)
	} else {
		re, err := regexp.Compile(regex)
		if err != nil {
			n.Status = "failed"
			return "", err
		}
		matches := re.FindStringSubmatch(input)
		if matches == nil {
			// no match -> fallback: template with ${body}
			out = strings.ReplaceAll(template, "${body}", input)
		} else {
			out = template
			for i := 1; i < len(matches); i++ {
				out = strings.ReplaceAll(out, fmt.Sprintf("${%d}", i), matches[i])
			}
			out = strings.ReplaceAll(out, "${body}", input)
		}
	}

	n.Data["output"] = out
	n.Status = "done"
	return "", nil
}
