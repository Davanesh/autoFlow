package services

import (
	"errors"
	"log"
)

/*
----------------------------------------------------
    EXECUTION ENGINE  (PATCHED WITH DATA PIPELINING)
----------------------------------------------------
*/

func RunWorkflow(g *ExecGraph) error {
	if g == nil {
		return errors.New("nil graph")
	}
	if g.Start == "" {
		return errors.New("no start node defined")
	}

	log.Println("🚀 Starting workflow execution (new engine)")
	current := g.Start

	for {
		n, ok := g.Nodes[current]
		if !ok {
			return errors.New("node not found: " + current)
		}

		executor, err := GetExecutor(n.Type)
		if err != nil {
			return errors.New("no executor for node type: " + n.Type)
		}

		// Execute node
		nextOverride, err := executor.Execute(n, g)
		if err != nil {
			n.Status = "failed"
			log.Printf("❌ Node failed: %s (%v)", n.ID, err)
			return err
		}

		// If node has exactly one next node, pass output -> next.input
		// DATA PIPELINING
if len(n.Next) == 1 {
    nextNodeID := n.Next[0]
    if nextNode, ok := g.Nodes[nextNodeID]; ok {

        log.Println("📤 Passing data FROM:", n.ID, "TO:", nextNodeID)

        if out, exists := n.Data["output"]; exists {

            log.Println("   └─ output:", out)

            if nextNode.Data == nil {
                nextNode.Data = map[string]interface{}{}
            }
            nextNode.Data["input"] = out

            log.Println("📥 next.input =", nextNode.Data["input"])
        } else {
            log.Println("⚠️ NO output found in node:", n.ID)
        }
    }
}


		// If EXECUTOR decided next node → override
		if nextOverride != "" {
			current = nextOverride
			continue
		}

		// No next nodes → end of workflow
		if len(n.Next) == 0 {
			log.Println("🏁 Workflow complete!")
			return nil
		}

		// Normal nodes must have only 1 next
		if len(n.Next) > 1 {
			return errors.New("node has multiple next branches but is not a decision: " + n.ID)
		}

		// Move forward
		current = n.Next[0]
	}
}
