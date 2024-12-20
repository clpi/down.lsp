package handler

import (
	"log"

	"github.com/clpi/down.lsp/lsp/handler/window"
	workspacecommand "github.com/clpi/down.lsp/lsp/handler/workspace/command"
	workspacefiles "github.com/clpi/down.lsp/lsp/handler/workspace/files"
	"github.com/tliron/commonlog"
	_ "github.com/tliron/commonlog/simple"
	"github.com/tliron/glsp"
	protocol "github.com/tliron/glsp/protocol_3_16"

	"github.com/clpi/down.lsp/lsp/handler/action"
	"github.com/clpi/down.lsp/lsp/handler/completion"
	"github.com/clpi/down.lsp/lsp/handler/document"
	"github.com/clpi/down.lsp/lsp/handler/lens"
	"github.com/clpi/down.lsp/lsp/handler/semantic"
	"github.com/clpi/down.lsp/lsp/handler/workspace"
	"github.com/clpi/down.lsp/lsp/handler/workspace/command"
)

var (
	name    string = "down"
	version string = "0.1.0-alpha"
	handler protocol.Handler
)

func initialized(context *glsp.Context, params *protocol.InitializedParams) error {
	return nil
}
func initialize(context *glsp.Context, params *protocol.InitializeParams) (any, error) {
	commonlog.NewInfoMessage(0, "down init...")
	capabilities := handler.CreateServerCapabilities()
	capabilities.CompletionProvider = &completion.Provider
	capabilities.CodeActionProvider = &action.Provider
	capabilities.CodeLensProvider = &lens.Provider
	capabilities.TextDocumentSync = &document.Provider.Sync
	capabilities.ExecuteCommandProvider = &command.Provider
	capabilities.DocumentLinkProvider = &document.Provider.Link
	capabilities.DocumentHighlightProvider = &document.Provider.Highlight
	capabilities.MonikerProvider = &document.Provider.Moniker
	capabilities.SemanticTokensProvider = &semantic.Provider
	capabilities.TextDocumentSync = &document.Provider.Sync
	capabilities.HoverProvider = &document.Provider.Hover
	capabilities.ColorProvider = &document.Provider.Color
	capabilities.DefinitionProvider = &document.Provider.Definition
	capabilities.DocumentSymbolProvider = &document.Provider.Symbol
	capabilities.WorkspaceSymbolProvider = &workspace.Provider
	capabilities.LinkedEditingRangeProvider = &document.Provider.LinkedEditingRange
	capabilities.SignatureHelpProvider = &document.Provider.SignatureHelp
	capabilities.Workspace = &workspace.Capabilities
	capabilities.ReferencesProvider = &document.Provider.References
	capabilities.Experimental = &map[string]interface{}{}
	capabilities.DocumentFormattingProvider = &document.Provider.Format
	capabilities.DocumentOnTypeFormattingProvider = &document.Provider.OnType
	capabilities.DocumentRangeFormattingProvider = &document.Provider.RangeFormat
	return protocol.InitializeResult{
		Capabilities: capabilities,
		ServerInfo: &protocol.InitializeResultServerInfo{
			Name:    name,
			Version: &version,
		},
	}, nil
}

func shutdown(context *glsp.Context) error {
	commonlog.NewInfoMessage(0, "down shutdown...")
	return nil
}

func setTrace(context *glsp.Context, params *protocol.SetTraceParams) error {
	protocol.SetTraceValue(params.Value)
	return nil
}
func logTrace(c *glsp.Context, p *protocol.LogTraceParams) error {
	if p.Verbose != nil {
		log.Print(p.Message, c.Method)
	} else {
		log.Print(p.Message, c.Method)
	}
	return nil
}

func LspHandler() protocol.Handler {
	return protocol.Handler{
		Initialize:                      initialize,
		Initialized:                     initialized,
		LogTrace:                        logTrace,
		SetTrace:                        setTrace,
		Shutdown:                        shutdown,
		Exit:                            shutdown,
		TextDocumentLinkedEditingRange:  document.LinkedEditing,
		TextDocumentDocumentHighlight:   document.DocumentHighlight,
		TextDocumentDocumentLink:        document.Link,
		TextDocumentCodeLens:            lens.CodeLens,
		WorkspaceDidChangeConfiguration: workspace.Configure,
		WorkspaceDidChangeWatchedFiles:  workspacefiles.DidWatch,
		WorkspaceDidCreateFiles:         workspacefiles.DidCreate,
		WorkspaceDidDeleteFiles:         workspacefiles.DidDelete,
		WorkspaceDidRenameFiles:         workspacefiles.DidRename,
		WorkspaceWillCreateFiles:        workspacefiles.Create,
		WorkspaceWillDeleteFiles:        workspacefiles.Delete,
		WorkspaceWillRenameFiles:        workspacefiles.Rename,
		WorkspaceExecuteCommand:         workspacecommand.Execute,
		DocumentLinkResolve:             document.LinkResolve,
		TextDocumentCodeAction:          action.CodeAction,
		TextDocumentColor:               document.Color,
		TextDocumentDocumentSymbol:      document.Symbol,
		WorkspaceSymbol:                 workspace.Symbol,
		Progress:                        window.Progress,

		WorkspaceDidChangeWorkspaceFolders: workspace.ChangeWorkspaceFolders,
		TextDocumentMoniker:                document.Moniker,
		CompletionItemResolve:              completion.ItemResolve,
		CodeActionResolve:                  action.Resolve,
		TextDocumentColorPresentation:      document.ColorPresentation,
		CodeLensResolve:                    lens.Resolve,
		TextDocumentCompletion:             completion.Completion,
		TextDocumentHover:                  document.Hover,
		TextDocumentReferences:             document.References,
		WindowWorkDoneProgressCancel:       window.Cancel,
		// TextDocumentSemanticTokensFullDelta: semantic.Delta,
		// TextDocumentSemanticTokensFull:      semantic.Full,
		// TextDocumentSemanticTokensRange:     semantic.Range,
		// WorkspaceSemanticTokensRefresh:      semantic.Workspace,
		// TextDocumentDeclaration:            document.Declaration,
		// TextDocumentDefinition:             document.Definition,
		// TextDocumentFormatting: document.Format,
	}
}

// Title: "Link current word to new file in workspace",
// IsPreferred: &trueVal,
// Disabled:    &falseVal,
// NewText: "[text](./text.md)",
