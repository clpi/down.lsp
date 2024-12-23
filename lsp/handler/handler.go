package handler

import (
	"log"
	"sync"

	"github.com/tliron/commonlog"
	_ "github.com/tliron/commonlog/simple"
	"github.com/tliron/glsp"
	protocol "github.com/tliron/glsp/protocol_3_16"

	"github.com/clpi/down.lsp/lsp/handler/completion"
	"github.com/clpi/down.lsp/lsp/handler/semantic"
)

var (
	name    string = "down"
	version string = "0.1.0-alpha"
	ph      protocol.Handler
)

type Handler struct {
	protocol.Handler
	initialized bool
	lock        sync.Mutex
}

func Capabilities() protocol.ServerCapabilities {
	cb := ph.CreateServerCapabilities()
	cb.CompletionProvider = &completion.Provider
	cb.CodeActionProvider = &ActionProvider
	cb.CodeLensProvider = &LensProvider
	cb.TextDocumentSync = &DocumentProvider.Sync
	cb.ExecuteCommandProvider = &CommandProvider
	cb.DocumentLinkProvider = &DocumentProvider.Link
	cb.DeclarationProvider = &DocumentProvider.Declaration
	cb.TypeDefinitionProvider = &DocumentProvider.TypeDefinition
	cb.ImplementationProvider = &DocumentProvider.Implementation
	cb.DocumentHighlightProvider = &DocumentProvider.Highlight
	cb.MonikerProvider = &DocumentProvider.Moniker
	cb.SemanticTokensProvider = &semantic.Provider
	cb.TextDocumentSync = &DocumentProvider.Sync
	cb.HoverProvider = &DocumentProvider.Hover
	cb.ColorProvider = &DocumentProvider.Color
	cb.DefinitionProvider = &DocumentProvider.Definition
	cb.DocumentSymbolProvider = &DocumentProvider.Symbol
	cb.WorkspaceSymbolProvider = protocol.WorkspaceSymbolOptions{}
	cb.Workspace = &protocol.ServerCapabilitiesWorkspace{
		WorkspaceFolders: &protocol.WorkspaceFoldersServerCapabilities{},
		FileOperations:   &WorkspaceFilesProvider,
	}
	cb.LinkedEditingRangeProvider = &DocumentProvider.LinkedEditingRange
	cb.SignatureHelpProvider = &SignatureOptions
	cb.ReferencesProvider = &DocumentProvider.References
	cb.Experimental = &map[string]interface{}{}
	cb.DocumentFormattingProvider = &DocumentProvider.Format
	cb.DocumentOnTypeFormattingProvider = &DocumentProvider.OnType
	cb.DocumentRangeFormattingProvider = &DocumentProvider.RangeFormat
	return cb
}

func ServerInfo() *protocol.InitializeResultServerInfo {
	return &protocol.InitializeResultServerInfo{
		Name:    name,
		Version: &version,
	}
}

func (s *State) Shutdown(context *glsp.Context) error {
	protocol.SetTraceValue(protocol.TraceValueOff)
	commonlog.NewInfoMessage(0, "down shutdown...")
	return nil
}

func (s *State) SetTrace(context *glsp.Context, params *protocol.SetTraceParams) error {
	protocol.SetTraceValue(params.Value)
	return nil
}

func (s *State) Cancel(c *glsp.Context, p *protocol.CancelParams) error {
	return nil
}

func (s *State) LogTrace(c *glsp.Context, p *protocol.LogTraceParams) error {
	if p.Verbose != nil {
		log.Print(p.Message, c.Method)
	} else {
		log.Print(p.Message, c.Method)
	}
	return nil
}

func (s *State) Initialized(c *glsp.Context, p *protocol.InitializedParams) error {
	protocol.SetTraceValue(protocol.TraceValueVerbose)
	return nil
}

func (s State) Handlers() protocol.Handler {
	return protocol.Handler{
		Initialize:                          s.Initialize,
		Initialized:                         s.Initialized,
		LogTrace:                            s.LogTrace,
		SetTrace:                            s.SetTrace,
		Shutdown:                            s.Shutdown,
		Exit:                                s.Exit,
		TextDocumentLinkedEditingRange:      s.LinkedEditing,
		TextDocumentDocumentHighlight:       s.DocumentHighlight,
		TextDocumentDocumentLink:            s.Links,
		TextDocumentCodeLens:                s.CodeLens,
		TextDocumentSemanticTokensFullDelta: s.Delta,
		WorkspaceSemanticTokensRefresh:      s.Refresh,
		TextDocumentSemanticTokensFull:      s.Full,
		TextDocumentSemanticTokensRange:     s.Range,
		TextDocumentSignatureHelp:           s.SignatureHelp,

		WorkspaceDidChangeConfiguration: s.Configure,
		WorkspaceDidChangeWatchedFiles:  s.WsDidWatch,
		WorkspaceDidCreateFiles:         s.WsDidCreate,
		WorkspaceDidDeleteFiles:         s.WsDidDelete,
		WorkspaceDidRenameFiles:         s.WsDidRename,
		WorkspaceWillCreateFiles:        s.WsWillCreate,
		WorkspaceWillDeleteFiles:        s.WsDelete,
		WorkspaceWillRenameFiles:        s.WsRename,
		WorkspaceExecuteCommand:         s.Command,
		DocumentLinkResolve:             s.LinkResolve,
		TextDocumentCodeAction:          s.CodeAction,
		TextDocumentRename:              s.Rename,
		TextDocumentPrepareRename:       s.PrepareRename,
		TextDocumentColor:               s.Color,
		TextDocumentDocumentSymbol:      s.Symbol,
		WorkspaceSymbol:                 s.WorkspaceSymbol,
		Progress:                        s.Progress,

		WorkspaceDidChangeWorkspaceFolders: s.ChangeWorkspaceFolders,
		TextDocumentMoniker:                s.Moniker,
		CompletionItemResolve:              s.ItemResolve,
		CodeActionResolve:                  s.ActionResolve,
		TextDocumentColorPresentation:      s.ColorPresentation,
		CodeLensResolve:                    s.LensResolve,
		TextDocumentCompletion:             s.Completion,
		TextDocumentHover:                  s.Hover,
		TextDocumentReferences:             s.References,
		WindowWorkDoneProgressCancel:       s.CancelWorkDoneProgresss,
		TextDocumentDidOpen:                s.DidOpen,
		TextDocumentDidChange:              s.DidChange,

		TextDocumentDidClose:          s.DidClose,
		TextDocumentWillSave:          s.WillSave,
		TextDocumentWillSaveWaitUntil: s.WillSaveWaitUntil,
		TextDocumentDidSave:           s.DidSave,
		CancelRequest:                 s.Cancel,
		// TextDocumentSignatureHelp:          s.SignatureHelp,

		// TextDocumentTypeDefinition:          document.TypeDefinition,
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
// Disabled:    &falseVal,
// IsPreferred: &trueVal,
// NewText: "[text](./text.md)",

func NewState() State {
	return State{}
}

func (s *State) Exit(context *glsp.Context) error {
	commonlog.NewInfoMessage(0, "down exitj...")
	return nil
}

func (s *State) Initialize(context *glsp.Context, params *protocol.InitializeParams) (any, error) {
	commonlog.NewInfoMessage(0, "down init...")
	protocol.SetTraceValue(protocol.TraceValueVerbose)
	return protocol.InitializeResult{
		Capabilities: Capabilities(),
		ServerInfo:   ServerInfo(),
	}, nil
}

func (s *State) Progress(context *glsp.Context, params *protocol.ProgressParams) error {
	return nil
}

func (s *State) CancelWorkDoneProgresss(context *glsp.Context, params *protocol.WorkDoneProgressCancelParams) error {
	return nil
}
