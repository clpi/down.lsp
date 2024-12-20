package workspace

import (
	wfiles "github.com/clpi/down.lsp/lsp/handler/workspace/files"
	protocol "github.com/tliron/glsp/protocol_3_16"
)

var (
	trueVal  = true
	falseVal = false
)

type (
	WorkspaceProvider struct {
		Files   protocol.ServerCapabilitiesWorkspaceFileOperations
		Folders protocol.WorkspaceFoldersServerCapabilities
		Symbol  protocol.WorkspaceSymbolOptions
	}
)

var (
	Provider = WorkspaceProvider{
		Files: wfiles.Provider,
		Folders: protocol.WorkspaceFoldersServerCapabilities{
			ChangeNotifications: &protocol.BoolOrString{true},
			Supported:           &trueVal,
		},
		Symbol: protocol.WorkspaceSymbolOptions{
			WorkDoneProgressOptions: protocol.WorkDoneProgressOptions{
				WorkDoneProgress: &trueVal,
			},
		},
	}
	Capabilities = protocol.ServerCapabilitiesWorkspace{
		FileOperations:   &Provider.Files,
		WorkspaceFolders: &Provider.Folders,
	}
)
