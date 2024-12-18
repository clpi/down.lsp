package handler

import (
	"github.com/tliron/commonlog"
	_ "github.com/tliron/commonlog/simple"
	"github.com/tliron/glsp"
	protocol "github.com/tliron/glsp/protocol_3_16"
	"log"
)

const (
	name = "down"
)

type (
	UserWorkspace struct {
		Name    string
		Path    protocol.URI
		Default bool
		Index   string
		Notes   string
		Config  interface{}
	}
	CodeTemplate struct {
		Body string
		URI  protocol.URI
	}
	CodeSnippet struct {
		Body        string
		Description string
	}
)

var (
	version  string = "0.1.0-alpha"
	handler  protocol.Handler
	trueVal  = true
	falseVal = false
	Emoji    = map[string]string{
		"happy":      "😀",
		"sad":        "😢",
		"angry":      "😠",
		"confused":   "😕",
		"excited":    "😆",
		"love":       "😍",
		"laughing":   "😂",
		"crying":     "😭",
		"sleepy":     "😴",
		"surprised":  "😮",
		"sick":       "🤒",
		"cool":       "😎",
		"nerd":       "🤓",
		"worried":    "😟",
		"scared":     "😨",
		"silly":      "🤪",
		"shocked":    "😱",
		"sunglasses": "😎",
		"tongue":     "😛",
		"thinking":   "🤔",
	}
	Snippets map[string]CodeSnippet = map[string]CodeSnippet{
		"#date": {
			Body:        "date +%Y-%m-%d",
			Description: "Insert the current date in the format YYYY-MM-DD",
		},
		"#time": {
			Body:        "date +%H:%M:%S",
			Description: "Insert the current time in the format HH:MM:SS",
		},
		"#datetime": {
			Body:        "date",
			Description: "Insert the current date and time",
		},
	}
)

func textDocumentCompletion(
	c *glsp.Context,
	p *protocol.CompletionParams,
) (interface{}, error) {
	var (
		preselect = true
		items     []protocol.CompletionItem
	)
	for s, sn := range Snippets {
		kind := protocol.CompletionItemKindSnippet
		items = append(items, protocol.CompletionItem{
			Label:     s,
			Kind:      &kind,
			Preselect: &preselect,
			Documentation: &protocol.MarkupContent{
				Value: "# Snippets\n\n## Snippet\n_ _ _\n### Snippet: " + sn.Description + "\n---\n" + sn.Body,
				Kind:  protocol.MarkupKindMarkdown,
			},
			Detail:     &sn.Description,
			InsertText: &sn.Body,
		})

	}
	for w, e := range Emoji {
		ec := e
		kind := protocol.CompletionItemKindConstant
		items = append(items, protocol.CompletionItem{
			Label:     w,
			Kind:      &kind,
			Preselect: &preselect,
			Documentation: &protocol.MarkupContent{
				Value: "# Emoji\n\n## Emoji\n_ _ _\n### Emoji: " + ec + "\n---\n" + ec,
				Kind:  protocol.MarkupKindMarkdown,
			},
			Detail:     &ec,
			InsertText: &ec,
		})
	}
	return items, nil
}

func documentHighlight(c *glsp.Context, p *protocol.DocumentHighlightParams) ([]protocol.DocumentHighlight, error) {
	var (
		hl []protocol.DocumentHighlight = []protocol.DocumentHighlight{}
		k1                              = protocol.DocumentHighlightKindText
		h1 protocol.DocumentHighlight   = protocol.DocumentHighlight{
			Kind: &k1,
			Range: protocol.Range{
				Start: protocol.Position{
					Line:      0,
					Character: 10,
				},
				End: protocol.Position{
					Line:      0,
					Character: 10,
				},
			},
		}
	)
	return append(hl, h1), nil
}
func codeLens(c *glsp.Context, p *protocol.CodeLensParams) ([]protocol.CodeLens, error) {
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

func initialized(context *glsp.Context, params *protocol.InitializedParams) error {
	return nil
}
func initialize(context *glsp.Context, params *protocol.InitializeParams) (any, error) {
	commonlog.NewInfoMessage(0, "down init...")
	capabilities := handler.CreateServerCapabilities()
	tru := true
	var (
		workspaceFoldersNotify protocol.BoolOrString = protocol.BoolOrString{
			Value: "workspace/didChangeFolders",
		}
	)
	capabilities.CompletionProvider = &protocol.CompletionOptions{
		ResolveProvider: &tru,
	}
	capabilities.CodeActionProvider = &protocol.CodeActionOptions{
		ResolveProvider: &tru,
	}
	capabilities.CodeLensProvider = &protocol.CodeLensOptions{
		ResolveProvider: &tru,
	}
	capabilities.TextDocumentSync = &protocol.TextDocumentSyncOptions{
		OpenClose: &tru,
		WillSave:  &tru,
		Save: &protocol.SaveOptions{
			IncludeText: &tru,
		},
	}
	capabilities.ExecuteCommandProvider = &protocol.ExecuteCommandOptions{
		WorkDoneProgressOptions: protocol.WorkDoneProgressOptions{},
		Commands: []string{
			"down.index",
			"down.workspace.open",
			"down.workspace.new",
			"down.workspace.delete",
			"down.link.create.cursor",
		},
	}
	capabilities.DocumentLinkProvider = &protocol.DocumentLinkOptions{
		ResolveProvider:         &tru,
		WorkDoneProgressOptions: protocol.WorkDoneProgressOptions{},
	}
	capabilities.Workspace = &protocol.ServerCapabilitiesWorkspace{
		WorkspaceFolders: &protocol.WorkspaceFoldersServerCapabilities{
			Supported:           &tru,
			ChangeNotifications: &workspaceFoldersNotify,
		},
		FileOperations: &protocol.ServerCapabilitiesWorkspaceFileOperations{},
	}
	capabilities.DocumentHighlightProvider = &protocol.DocumentHighlightOptions{
		WorkDoneProgressOptions: protocol.WorkDoneProgressOptions{},
	}
	capabilities.MonikerProvider = &protocol.MonikerOptions{
		WorkDoneProgressOptions: protocol.WorkDoneProgressOptions{},
	}
	capabilities.SemanticTokensProvider = &protocol.SemanticTokensOptions{
		WorkDoneProgressOptions: protocol.WorkDoneProgressOptions{},
	}
	capabilities.TextDocumentSync = &protocol.TextDocumentSyncOptions{}
	capabilities.HoverProvider = &protocol.HoverOptions{}
	capabilities.ColorProvider = &protocol.DocumentColorOptions{}
	capabilities.ReferencesProvider = &protocol.ReferenceOptions{}
	capabilities.Experimental = &map[string]interface{}{}
	// capabilities.HoverProvider = &protocol.HoverOptions{}
	// capabilities.DefinitionProvider = &protocol.DefinitionOptions{}
	// capabilities.DocumentSymbolProvider = &protocol.DocumentSymbolOptions{}
	// capabilities.DocumentFormattingProvider = &protocol.DocumentFormattingOptions{}
	// capabilities.DocumentRangeFormattingProvider = &protocol.DocumentRangeFormattingOptions{}
	return protocol.InitializeResult{
		Capabilities: capabilities,
		ServerInfo: &protocol.InitializeResultServerInfo{
			Name:    name,
			Version: &version,
		},
	}, nil
}

func moniker(c *glsp.Context, p *protocol.MonikerParams) ([]protocol.Moniker, error) {
	var (
		k1       protocol.MonikerKind = protocol.MonikerKindLocal
		monikers []protocol.Moniker
		m1       protocol.Moniker = protocol.Moniker{
			Unique:     "unique",
			Identifier: "identifier",
			Kind:       &k1,
			Scheme:     "file",
		}
	)
	return append(monikers, m1), nil
}

func codeAction(c *glsp.Context, p *protocol.CodeActionParams) (any, error) {
	var (
		actions   []protocol.CodeAction   = []protocol.CodeAction{}
		actionSrc protocol.CodeActionKind = protocol.CodeActionKindSource
		_         protocol.Range          = protocol.Range{
			Start: protocol.Position{
				Line:      10,
				Character: 10,
			},
			End: protocol.Position{
				Line:      10,
				Character: 20,
			},
		}
		cursorCreateLink protocol.CodeAction = protocol.CodeAction{
			Command: &protocol.Command{
				Arguments: []any{
					"dir",
				},
				Command: "down.link.create.cursor",
				Title:   "Create link on cursor word",
			},
			Diagnostics: nil,
			IsPreferred: &trueVal,
			Kind:        &actionSrc,
			Data:        nil,
			Edit:        &protocol.WorkspaceEdit{},
			Title:       "Create link on word/heading",
		}
	)
	return append(actions, cursorCreateLink), nil
}
func shutdown(context *glsp.Context) error {
	commonlog.NewInfoMessage(0, "down shutdown...")
	return nil
}

func setTrace(context *glsp.Context, params *protocol.SetTraceParams) error {
	protocol.SetTraceValue(params.Value)
	return nil
}
func semanticTokens(c *glsp.Context, p *protocol.SemanticTokensParams) (*protocol.SemanticTokens, error) {
	var (
		rid string                  = "result-id"
		st  protocol.SemanticTokens = protocol.SemanticTokens{
			Data: []protocol.UInteger{
				10, 20, 30,
			},
			ResultID: &rid,
		}
	)
	return &st, nil
}
func documentLink(c *glsp.Context, p *protocol.DocumentLinkParams) ([]protocol.DocumentLink, error) {
	var (
		dls []protocol.DocumentLink
		d1  protocol.DocumentLink = protocol.DocumentLink{
			Tooltip: &p.TextDocument.URI,
			Data:    nil,
			Target:  &p.TextDocument.URI,
			Range: protocol.Range{
				Start: protocol.Position{
					Line:      10,
					Character: 10,
				},
				End: protocol.Position{
					Line:      20,
					Character: 20,
				},
			},
		}
	)
	return append(dls, d1), nil
}
func linkedEditing(c *glsp.Context, p *protocol.LinkedEditingRangeParams) (*protocol.LinkedEditingRanges, error) {
	var (
		p1 string                       = ""
		r1 protocol.LinkedEditingRanges = protocol.LinkedEditingRanges{
			WordPattern: &p1,
			Ranges: []protocol.Range{
				{
					Start: protocol.Position{
						Line:      10,
						Character: 10,
					},
					End: protocol.Position{
						Line:      20,
						Character: 20,
					},
				},
				{
					Start: protocol.Position{
						Line:      10,
						Character: 10,
					},
					End: protocol.Position{
						Line:      10,
						Character: 20,
					},
				},
			},
		}
	)
	return &r1, nil
}

func workspaceCommand(c *glsp.Context, p *protocol.ExecuteCommandParams) (any, error) {
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

func documentColor(c *glsp.Context, p *protocol.DocumentColorParams) ([]protocol.ColorInformation, error) {
	var ci []protocol.ColorInformation
	return append(ci, protocol.ColorInformation{}), nil
}
func hover(c *glsp.Context, p *protocol.HoverParams) (*protocol.Hover, error) {
	return &protocol.Hover{
		Contents: protocol.MarkupContent{
			Kind:  protocol.MarkupKindMarkdown,
			Value: "# Info\n\n## Information\n\n### About\n\nThis is a hover tooltip",
		},
		Range: &protocol.Range{
			Start: protocol.Position{
				Line:      0,
				Character: 0,
			},
			End: protocol.Position{
				Line:      0,
				Character: 0,
			},
		},
	}, nil
}
func references(c *glsp.Context, p *protocol.ReferenceParams) ([]protocol.Location, error) {
	var refs []protocol.Location
	return append(refs, protocol.Location{}), nil
}
func documentLinkResolve(c *glsp.Context, p *protocol.DocumentLink) (*protocol.DocumentLink, error) {
	return p, nil
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
		Initialize:                     initialize,
		LogTrace:                       logTrace,
		Exit:                           shutdown,
		TextDocumentDocumentHighlight:  documentHighlight,
		TextDocumentLinkedEditingRange: linkedEditing,
		TextDocumentDocumentLink:       documentLink,
		TextDocumentCodeLens:           codeLens,
		SetTrace:                       setTrace,
		WorkspaceExecuteCommand:        workspaceCommand,
		Initialized:                    initialized,
		DocumentLinkResolve:            documentLinkResolve,
		TextDocumentCodeAction:         codeAction,
		TextDocumentColor:              documentColor,
		TextDocumentMoniker:            moniker,
		TextDocumentCompletion:         textDocumentCompletion,
		TextDocumentSemanticTokensFull: semanticTokens,
		TextDocumentHover:              hover,
		TextDocumentReferences:         references,
		Shutdown:                       shutdown,
	}
}

// Title: "Link current word to new file in workspace",
// IsPreferred: &trueVal,
// Disabled:    &falseVal,
// NewText: "[text](./text.md)",
