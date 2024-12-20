package files

import (
	"github.com/clpi/down.lsp/lsp/files"
	protocol "github.com/tliron/glsp/protocol_3_16"
)

var (
	Provider = protocol.ServerCapabilitiesWorkspaceFileOperations{
		DidCreate: &protocol.FileOperationRegistrationOptions{
			Filters: files.FileOps,
		},
		WillCreate: &protocol.FileOperationRegistrationOptions{
			Filters: files.FileOps,
		},
		DidRename: &protocol.FileOperationRegistrationOptions{
			Filters: files.FileOps,
		},
		WillRename: &protocol.FileOperationRegistrationOptions{
			Filters: files.FileOps,
		},
		WillDelete: &protocol.FileOperationRegistrationOptions{
			Filters: files.FileOps,
		},
		DidDelete: &protocol.FileOperationRegistrationOptions{
			Filters: files.FileOps,
		},
	}
)
