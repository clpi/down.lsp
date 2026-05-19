package handler

import (
	"strings"

	files "github.com/clpi/down.lsp/lsp/files"
	"github.com/clpi/down.lsp/lsp/knowledge"
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

func (s *State) DidSave(ctx *glsp.Context, params *protocol.DidSaveTextDocumentParams) error {
	uri := string(params.TextDocument.URI)
	if params.Text != nil {
		s.Documents[uri] = *params.Text
		knowledge.ExtractFromDocument(s.Graph, uri, *params.Text)
		s.Graph.Save()
		s.publishDiagnostics(ctx, uri, *params.Text)
	}
	return nil
}

func (s *State) DidClose(ctx *glsp.Context, params *protocol.DidCloseTextDocumentParams) error {
	uri := string(params.TextDocument.URI)
	delete(s.Documents, uri)
	ctx.Notify(protocol.ServerTextDocumentPublishDiagnostics, protocol.PublishDiagnosticsParams{
		URI:         protocol.DocumentUri(uri),
		Diagnostics: []protocol.Diagnostic{},
	})
	return nil
}

func (s *State) WillSaveWaitUntil(context *glsp.Context, params *protocol.WillSaveTextDocumentParams) ([]protocol.TextEdit, error) {
	return nil, nil
}

func (s *State) WillSave(context *glsp.Context, params *protocol.WillSaveTextDocumentParams) error {
	return nil
}

func (s *State) DidOpen(ctx *glsp.Context, params *protocol.DidOpenTextDocumentParams) error {
	uri := string(params.TextDocument.URI)
	text := params.TextDocument.Text
	s.Documents[uri] = text
	knowledge.ExtractFromDocument(s.Graph, uri, text)
	s.publishDiagnostics(ctx, uri, text)
	return nil
}

func (s *State) DidChange(context *glsp.Context, params *protocol.DidChangeTextDocumentParams) error {
	uri := string(params.TextDocument.URI)
	for _, change := range params.ContentChanges {
		if c, ok := change.(protocol.TextDocumentContentChangeEventWhole); ok {
			s.Documents[uri] = c.Text
			knowledge.ExtractFromDocument(s.Graph, uri, c.Text)
		}
	}
	return nil
}

func (s *State) Rename(_ *glsp.Context, p *protocol.RenameParams) (*protocol.WorkspaceEdit, error) {
	uri := string(p.TextDocument.URI)
	text, ok := s.Documents[uri]
	if !ok {
		return nil, nil
	}

	lines := strings.Split(text, "\n")
	line := int(p.Position.Line)
	if line >= len(lines) {
		return nil, nil
	}

	word := wordAtPosition(lines[line], int(p.Position.Character))
	if word == "" {
		return nil, nil
	}

	newName := p.NewName

	// Find all documents containing this word and build edits
	changes := map[protocol.DocumentUri][]protocol.TextEdit{}
	for docURI, docText := range s.Documents {
		docLines := strings.Split(docText, "\n")
		for lineIdx, l := range docLines {
			col := 0
			lower := strings.ToLower(l)
			target := strings.ToLower(word)
			for {
				idx := strings.Index(lower[col:], target)
				if idx < 0 {
					break
				}
				pos := col + idx
				// Check word boundaries
				if pos > 0 && isWordChar(l[pos-1]) {
					col = pos + 1
					continue
				}
				end := pos + len(word)
				if end < len(l) && isWordChar(l[end]) {
					col = pos + 1
					continue
				}
				changes[protocol.DocumentUri(docURI)] = append(
					changes[protocol.DocumentUri(docURI)],
					protocol.TextEdit{
						Range: protocol.Range{
							Start: protocol.Position{Line: protocol.UInteger(lineIdx), Character: protocol.UInteger(pos)},
							End:   protocol.Position{Line: protocol.UInteger(lineIdx), Character: protocol.UInteger(end)},
						},
						NewText: newName,
					},
				)
				col = end
			}
		}
	}

	if len(changes) == 0 {
		return nil, nil
	}

	return &protocol.WorkspaceEdit{Changes: changes}, nil
}

func (s *State) PrepareRename(_ *glsp.Context, p *protocol.PrepareRenameParams) (any, error) {
	uri := string(p.TextDocument.URI)
	text, ok := s.Documents[uri]
	if !ok {
		return nil, nil
	}

	lines := strings.Split(text, "\n")
	line := int(p.Position.Line)
	if line >= len(lines) {
		return nil, nil
	}

	word := wordAtPosition(lines[line], int(p.Position.Character))
	if word == "" {
		return nil, nil
	}

	// Find the word start position
	col := int(p.Position.Character)
	start := col
	for start > 0 && isWordChar(lines[line][start-1]) {
		start--
	}

	return protocol.Range{
		Start: protocol.Position{Line: p.Position.Line, Character: protocol.UInteger(start)},
		End:   protocol.Position{Line: p.Position.Line, Character: protocol.UInteger(start + len(word))},
	}, nil
}

func wordAtPosition(line string, col int) string {
	if col >= len(line) {
		return ""
	}
	start := col
	for start > 0 && isWordChar(line[start-1]) {
		start--
	}
	end := col
	for end < len(line) && isWordChar(line[end]) {
		end++
	}
	if start == end {
		return ""
	}
	return line[start:end]
}

func (s *State) Moniker(_ *glsp.Context, _ *protocol.MonikerParams) ([]protocol.Moniker, error) {
	return nil, nil
}

func (s *State) References(_ *glsp.Context, p *protocol.ReferenceParams) ([]protocol.Location, error) {
	if s.Graph == nil {
		return nil, nil
	}

	uri := string(p.TextDocument.URI)
	doc, ok := s.Documents[uri]
	if !ok {
		return nil, nil
	}

	lines := strings.Split(doc, "\n")
	lineIdx := int(p.Position.Line)
	if lineIdx >= len(lines) {
		return nil, nil
	}
	line := lines[lineIdx]
	col := int(p.Position.Character)
	if col >= len(line) {
		return nil, nil
	}

	wordStart, wordEnd := col, col
	for wordStart > 0 && isWordChar(line[wordStart-1]) {
		wordStart--
	}
	for wordEnd < len(line) && isWordChar(line[wordEnd]) {
		wordEnd++
	}
	if wordStart >= wordEnd {
		return nil, nil
	}
	word := line[wordStart:wordEnd]

	results := s.Graph.Search(word)
	var locations []protocol.Location
	for _, ent := range results {
		if strings.ToLower(ent.Name) != strings.ToLower(word) {
			continue
		}
		for _, src := range ent.Sources {
			if !p.Context.IncludeDeclaration && src.URI == uri && src.Line == lineIdx {
				continue
			}
			locations = append(locations, protocol.Location{
				URI: protocol.DocumentUri(src.URI),
				Range: protocol.Range{
					Start: protocol.Position{Line: protocol.UInteger(src.Line), Character: 0},
					End:   protocol.Position{Line: protocol.UInteger(src.Line), Character: protocol.UInteger(len(ent.Name))},
				},
			})
		}
	}
	return locations, nil
}

func (s *State) Color(c *glsp.Context, p *protocol.DocumentColorParams) ([]protocol.ColorInformation, error) {
	var ci []protocol.ColorInformation
	return append(ci, protocol.ColorInformation{}), nil
}

func (s *State) ColorPresentation(c *glsp.Context, p *protocol.ColorPresentationParams) ([]protocol.ColorPresentation, error) {
	return []protocol.ColorPresentation{}, nil
}

func (s *State) Symbol(_ *glsp.Context, p *protocol.DocumentSymbolParams) (any, error) {
	if s.Graph == nil {
		return nil, nil
	}

	uri := string(p.TextDocument.URI)
	entities := s.Graph.EntitiesByDocument(uri)
	if len(entities) == 0 {
		return nil, nil
	}

	symbolKindMap := map[knowledge.EntityKind]protocol.SymbolKind{
		knowledge.KindPerson:   protocol.SymbolKindVariable,
		knowledge.KindConcept:  protocol.SymbolKindClass,
		knowledge.KindProject:  protocol.SymbolKindPackage,
		knowledge.KindAction:   protocol.SymbolKindFunction,
		knowledge.KindTag:      protocol.SymbolKindKey,
		knowledge.KindDocument: protocol.SymbolKindFile,
		knowledge.KindDate:     protocol.SymbolKindEvent,
		knowledge.KindPlace:    protocol.SymbolKindNamespace,
		knowledge.KindCode:     protocol.SymbolKindObject,
	}

	var symbols []protocol.DocumentSymbol
	for _, ent := range entities {
		kind, ok := symbolKindMap[ent.Kind]
		if !ok {
			kind = protocol.SymbolKindString
		}

		line := protocol.UInteger(0)
		for _, src := range ent.Sources {
			if src.URI == uri {
				line = protocol.UInteger(src.Line)
				break
			}
		}

		detail := string(ent.Kind)
		symbols = append(symbols, protocol.DocumentSymbol{
			Name:   ent.Name,
			Detail: &detail,
			Kind:   kind,
			Range: protocol.Range{
				Start: protocol.Position{Line: line, Character: 0},
				End:   protocol.Position{Line: line, Character: protocol.UInteger(len(ent.Name))},
			},
			SelectionRange: protocol.Range{
				Start: protocol.Position{Line: line, Character: 0},
				End:   protocol.Position{Line: line, Character: protocol.UInteger(len(ent.Name))},
			},
		})
	}

	return symbols, nil
}

func (s *State) Definition(_ *glsp.Context, p *protocol.DefinitionParams) (any, error) {
	if s.Graph == nil {
		return nil, nil
	}

	uri := string(p.TextDocument.URI)
	doc, ok := s.Documents[uri]
	if !ok {
		return nil, nil
	}

	lines := strings.Split(doc, "\n")
	lineIdx := int(p.Position.Line)
	if lineIdx >= len(lines) {
		return nil, nil
	}
	line := lines[lineIdx]
	col := int(p.Position.Character)
	if col >= len(line) {
		return nil, nil
	}

	wordStart, wordEnd := col, col
	for wordStart > 0 && isWordChar(line[wordStart-1]) {
		wordStart--
	}
	for wordEnd < len(line) && isWordChar(line[wordEnd]) {
		wordEnd++
	}
	if wordStart >= wordEnd {
		return nil, nil
	}
	word := line[wordStart:wordEnd]

	results := s.Graph.Search(word)
	if len(results) == 0 {
		return nil, nil
	}

	var locations []protocol.Location
	for _, ent := range results {
		if strings.ToLower(ent.Name) != strings.ToLower(word) {
			continue
		}
		for _, src := range ent.Sources {
			if src.URI == uri && src.Line == lineIdx {
				continue
			}
			locations = append(locations, protocol.Location{
				URI: protocol.DocumentUri(src.URI),
				Range: protocol.Range{
					Start: protocol.Position{Line: protocol.UInteger(src.Line), Character: 0},
					End:   protocol.Position{Line: protocol.UInteger(src.Line), Character: protocol.UInteger(len(ent.Name))},
				},
			})
		}
	}

	if len(locations) == 0 {
		return nil, nil
	}
	if len(locations) == 1 {
		return locations[0], nil
	}
	return locations, nil
}

func isWordChar(b byte) bool {
	return (b >= 'a' && b <= 'z') || (b >= 'A' && b <= 'Z') ||
		(b >= '0' && b <= '9') || b == '_' || b == '-' || b == '.'
}

func (s *State) publishDiagnostics(ctx *glsp.Context, uri string, text string) {
	if s.Graph == nil {
		return
	}

	kDiags := knowledge.DiagnoseDocument(s.Graph, uri, text)
	lspDiags := make([]protocol.Diagnostic, len(kDiags))

	sevMap := map[knowledge.DiagSeverity]protocol.DiagnosticSeverity{
		knowledge.SeverityError:   protocol.DiagnosticSeverityError,
		knowledge.SeverityWarning: protocol.DiagnosticSeverityWarning,
		knowledge.SeverityInfo:    protocol.DiagnosticSeverityInformation,
		knowledge.SeverityHint:    protocol.DiagnosticSeverityHint,
	}

	source := "down"
	for i, d := range kDiags {
		sev := sevMap[d.Severity]
		lspDiags[i] = protocol.Diagnostic{
			Range: protocol.Range{
				Start: protocol.Position{Line: protocol.UInteger(d.Line), Character: protocol.UInteger(d.ColStart)},
				End:   protocol.Position{Line: protocol.UInteger(d.Line), Character: protocol.UInteger(d.ColEnd)},
			},
			Severity: &sev,
			Source:   &source,
			Message:  d.Message,
		}
	}

	ctx.Notify(protocol.ServerTextDocumentPublishDiagnostics, protocol.PublishDiagnosticsParams{
		URI:         protocol.DocumentUri(uri),
		Diagnostics: lspDiags,
	})
}
