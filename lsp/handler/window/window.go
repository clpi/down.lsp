package window

import (
	"github.com/tliron/glsp"
	protocol "github.com/tliron/glsp/protocol_3_16"
)

func Cancel(c *glsp.Context, p *protocol.WorkDoneProgressCancelParams) error {
	return nil
}

func Progress(c *glsp.Context, p *protocol.ProgressParams) error {
	return nil
}
