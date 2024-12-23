package handler

import (
	"github.com/clpi/down.lsp/lsp/files"
	"github.com/tliron/glsp"
	protocol "github.com/tliron/glsp/protocol_3_16"
)

var (
	WorkspaceFilesProvider = protocol.ServerCapabilitiesWorkspaceFileOperations{
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

func (s *State) WsCreate(c *glsp.Context, p *protocol.CreateFilesParams) (*protocol.WorkspaceEdit, error) {
	return nil, nil
}

func (s *State) WsDidCreate(c *glsp.Context, p *protocol.CreateFilesParams) error {
	return nil
}

func (s *State) WsDelete(c *glsp.Context, p *protocol.DeleteFilesParams) (*protocol.WorkspaceEdit, error) {
	return nil, nil
}

func (s *State) WsDidDelete(c *glsp.Context, p *protocol.DeleteFilesParams) error {
	return nil
}
func (s *State) WsWatch(c *glsp.Context, p *protocol.DidChangeWatchedFilesParams) (*protocol.WorkspaceEdit, error) {
	return nil, nil
}

func (s *State) WsDidWatch(c *glsp.Context, p *protocol.DidChangeWatchedFilesParams) error {
	return nil
}
func (s *State) WsRename(c *glsp.Context, p *protocol.RenameFilesParams) (*protocol.WorkspaceEdit, error) {
	return nil, nil
}

func (s *State) WsDidRename(c *glsp.Context, p *protocol.RenameFilesParams) error {
	return nil
}
func (s *State) WsWillCreate(c *glsp.Context, p *protocol.CreateFilesParams) (*protocol.WorkspaceEdit, error) {
	return nil, nil

}
func (s *State) Configure(c *glsp.Context, p *protocol.DidChangeConfigurationParams) error {
	// s := map[string]interface{}{
	// 	"markdown": map[string]interface{}{},
	// 	"down": map[string]interface{}{
	// 		"enabled": true,
	// 		"codeAction": map[string]interface{}{
	// 			"enabled": true,
	// 		},
	// 		"codeLens": map[string]interface{}{
	// 			"enabled": true,
	// 		},
	// 		"inlayHint": map[string]interface{}{
	// 			"enabled": true,
	// 		},
	// 		"diagnostics": map[string]interface{}{
	// 			"enabled": true,
	// 		},
	// 		"hover": map[string]interface{}{
	// 			"enabled": true,
	// 		},
	// 		"completion": map[string]interface{}{
	// 			"enabled": true,
	// 		},
	// 		"signatureHelp": map[string]interface{}{
	// 			"enabled": true,
	// 		},
	// 	},
	// }
	p.Settings = s
	p.Settings = s
	return nil
}
func (s *State) WorkspaceSymbol(*glsp.Context, *protocol.WorkspaceSymbolParams) ([]protocol.SymbolInformation, error) {
	return []protocol.SymbolInformation{}, nil
}
func (s *State) ChangeWorkspaceFolders(c *glsp.Context, p *protocol.DidChangeWorkspaceFoldersParams) error {
	return nil
}
