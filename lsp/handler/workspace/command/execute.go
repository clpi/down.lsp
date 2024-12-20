package command

import (
	"log"

	"github.com/tliron/glsp"
	protocol "github.com/tliron/glsp/protocol_3_16"
)

func Execute(c *glsp.Context, p *protocol.ExecuteCommandParams) (any, error) {
	var args = p.Arguments
	log.Print(p.Command, p.Arguments)
	switch p.Command {
	case "down.index":
		if len(args) == 0 {
			const _ = "default"
		} else {
			const _ = "default"
		}
	case "down.workspace.open":
	case "down.workspace.new":
	default:
	}
	return nil, nil
}
