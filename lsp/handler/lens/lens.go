package lens

import (
	"github.com/clpi/down.lsp/lsp/files"
	"github.com/tliron/glsp"
	protocol "github.com/tliron/glsp/protocol_3_16"
)

var (
	trueVal  bool = true
	falseVal bool = false
	Provider      = protocol.CodeLensOptions{
		ResolveProvider: &trueVal,
		WorkDoneProgressOptions: protocol.WorkDoneProgressOptions{
			WorkDoneProgress: &trueVal,
		},
	}
	Registration = protocol.CodeLensRegistrationOptions{
		TextDocumentRegistrationOptions: protocol.TextDocumentRegistrationOptions{
			DocumentSelector: &files.Filetypes,
		},
		CodeLensOptions: Provider,
	}
)

func CodeLens(c *glsp.Context, p *protocol.CodeLensParams) ([]protocol.CodeLens, error) {
	var lens []protocol.CodeLens
	var (
		wsOpen protocol.CodeLens = protocol.CodeLens{
			Data: nil,
			Command: &protocol.Command{
				Command:   "down.workspace.open",
				Arguments: nil,
				Title:     "Open workspace",
			},
		}
		wsNew protocol.CodeLens = protocol.CodeLens{
			Data: nil,
			Command: &protocol.Command{
				Arguments: nil,
				Title:     "new workspace",
				Command:   "down.workspace.new",
			},
		}
	)
	return append(lens, wsOpen, wsNew), nil

}
