package handler

import (
	files "github.com/clpi/down.lsp/lsp/files"
	"github.com/tliron/glsp"
	protocol "github.com/tliron/glsp/protocol_3_16"
)

var (
	file       = "file"
	partialTok = "1"
)

type (
	DocProvider struct {
		Registration       protocol.TextDocumentRegistrationOptions
		Link               protocol.DocumentLinkOptions
		Highlight          protocol.DocumentHighlightOptions
		Implementation     protocol.ImplementationOptions
		References         protocol.ReferenceParams
		Declaration        protocol.DeclarationOptions
		Definition         protocol.DefinitionOptions
		TypeDefinition     protocol.TypeDefinitionOptions
		LinkedEditingRange protocol.LinkedEditingRangeOptions
		Moniker            protocol.MonikerOptions
		Symbol             protocol.DocumentSymbolOptions
		Color              protocol.DocumentColorOptions
		Format             protocol.DocumentFormattingOptions
		RangeFormat        protocol.DocumentRangeFormattingOptions
		OnType             protocol.DocumentOnTypeFormattingOptions
		Sync               protocol.TextDocumentSyncOptions
		Hover              protocol.HoverOptions
	}
)

var DocumentProvider = DocProvider{
	Registration: protocol.TextDocumentRegistrationOptions{
		DocumentSelector: &files.Filetypes,
	},
	Sync: protocol.TextDocumentSyncOptions{
		OpenClose:         &trueVal,
		WillSave:          &trueVal,
		WillSaveWaitUntil: &trueVal,
		Save: &protocol.SaveOptions{
			IncludeText: &trueVal,
		},
	},
	Highlight: protocol.DocumentHighlightOptions{
		WorkDoneProgressOptions: workDone,
	},
	Implementation: protocol.ImplementationOptions{
		WorkDoneProgressOptions: workDone,
	},
	References: protocol.ReferenceParams{
		WorkDoneProgressParams: protocol.WorkDoneProgressParams{
			WorkDoneToken: &protocol.ProgressToken{
				Value: partialTok,
			},
		},
		TextDocumentPositionParams: protocol.TextDocumentPositionParams{
			TextDocument: protocol.TextDocumentIdentifier{
				URI: "",
			},
			Position: protocol.Position{
				Line:      0,
				Character: 0,
			},
		},
		PartialResultParams: protocol.PartialResultParams{
			PartialResultToken: &protocol.ProgressToken{"1"},
		},
		Context: protocol.ReferenceContext{
			IncludeDeclaration: true,
		},
	},
	Declaration: protocol.DeclarationOptions{
		WorkDoneProgressOptions: workDone,
	},
	Definition: protocol.DefinitionOptions{
		WorkDoneProgressOptions: workDone,
	},
	TypeDefinition: protocol.TypeDefinitionOptions{
		WorkDoneProgressOptions: workDone,
	},
	LinkedEditingRange: protocol.LinkedEditingRangeOptions{
		WorkDoneProgressOptions: workDone,
	},
	Moniker: protocol.MonikerOptions{
		WorkDoneProgressOptions: workDone,
	},
	Symbol: protocol.DocumentSymbolOptions{
		WorkDoneProgressOptions: workDone,
	},
	Hover: protocol.HoverOptions{
		WorkDoneProgressOptions: workDone,
	},
	Link: protocol.DocumentLinkOptions{
		ResolveProvider:         &trueVal,
		WorkDoneProgressOptions: workDone,
	},
	Color: protocol.DocumentColorOptions{
		WorkDoneProgressOptions: workDone,
	},
	Format: protocol.DocumentFormattingOptions{
		WorkDoneProgressOptions: workDone,
	},
	RangeFormat: protocol.DocumentRangeFormattingOptions{
		WorkDoneProgressOptions: workDone,
	},
	OnType: protocol.DocumentOnTypeFormattingOptions{
		FirstTriggerCharacter: " ",
		MoreTriggerCharacter:  []string{" ", "\n"},
	},
}

func (s *State) DidSave(context *glsp.Context, params *protocol.DidSaveTextDocumentParams) error {
	return nil
}

func (s *State) DidClose(context *glsp.Context, params *protocol.DidCloseTextDocumentParams) error {
	return nil
}

func (s *State) WillSaveWaitUntil(context *glsp.Context, params *protocol.WillSaveTextDocumentParams) ([]protocol.TextEdit, error) {
	return nil, nil
}

func (s *State) WillSave(context *glsp.Context, params *protocol.WillSaveTextDocumentParams) error {
	return nil
}

func (s *State) DidOpen(context *glsp.Context, params *protocol.DidOpenTextDocumentParams) error {
	return nil
}

func (s *State) DidChange(context *glsp.Context, params *protocol.DidChangeTextDocumentParams) error {
	return nil
}

func (s *State) Rename(context *glsp.Context, p *protocol.RenameParams) (*protocol.WorkspaceEdit, error) {
	return nil, nil
}

func (s *State) PrepareRename(context *glsp.Context, p *protocol.PrepareRenameParams) (any, error) {
	return nil, nil
}

func (s *State) Moniker(c *glsp.Context, p *protocol.MonikerParams) ([]protocol.Moniker, error) {
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

func (s *State) References(c *glsp.Context, p *protocol.ReferenceParams) ([]protocol.Location, error) {
	var refs []protocol.Location
	return append(refs, protocol.Location{}), nil
}

func (s *State) Color(c *glsp.Context, p *protocol.DocumentColorParams) ([]protocol.ColorInformation, error) {
	var ci []protocol.ColorInformation
	return append(ci, protocol.ColorInformation{}), nil
}

func (s *State) ColorPresentation(c *glsp.Context, p *protocol.ColorPresentationParams) ([]protocol.ColorPresentation, error) {
	return []protocol.ColorPresentation{}, nil
}

func (s *State) Symbol(*glsp.Context, *protocol.DocumentSymbolParams) (any, error) {
	return nil, nil
}
