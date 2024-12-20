package files

import (
	"github.com/tliron/glsp"
	protocol "github.com/tliron/glsp/protocol_3_16"
)

func Watch(c *glsp.Context, p *protocol.DidChangeWatchedFilesParams) (*protocol.WorkspaceEdit, error) {
	return nil, nil
}

func DidWatch(c *glsp.Context, p *protocol.DidChangeWatchedFilesParams) error {
	return nil
}
