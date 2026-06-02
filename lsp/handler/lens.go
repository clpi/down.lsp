package handler

import (
	"fmt"
	"strings"

	"github.com/clpi/down.lsp/lsp/files"
	"github.com/tliron/glsp"
	protocol "github.com/tliron/glsp/protocol_3_16"
)

var (
	LensProvider = protocol.CodeLensOptions{
		ResolveProvider: &trueVal,
		WorkDoneProgressOptions: protocol.WorkDoneProgressOptions{
			WorkDoneProgress: &trueVal,
		},
	}
	LensRegistration = protocol.CodeLensRegistrationOptions{
		TextDocumentRegistrationOptions: protocol.TextDocumentRegistrationOptions{
			DocumentSelector: &files.Filetypes,
		},
		CodeLensOptions: LensProvider,
	}
)

func (s *State) CodeLens(c *glsp.Context, p *protocol.CodeLensParams) ([]protocol.CodeLens, error) {
	var lens []protocol.CodeLens
	uri := string(p.TextDocument.URI)
	text, ok := s.Documents[uri]
	if !ok {
		return lens, nil
	}

	// Document metadata lens (line 0)
	meta := s.computeDocMetadata(uri, text)
	if meta != "" {
		lens = append(lens, protocol.CodeLens{
			Range: protocol.Range{
				Start: protocol.Position{Line: 0, Character: 0},
				End:   protocol.Position{Line: 0, Character: 0},
			},
			Command: &protocol.Command{
				Command: "down.document.info",
				Title:   meta,
			},
		})
	}

	// Document type lens
	info := s.DetectDocumentType(uri)
	if info != nil && info.Type != DocTypeGeneric {
		typeLens := fmt.Sprintf("%s %s", info.Icon, info.Type)
		lens = append(lens, protocol.CodeLens{
			Range: protocol.Range{
				Start: protocol.Position{Line: 0, Character: 0},
				End:   protocol.Position{Line: 0, Character: 0},
			},
			Command: &protocol.Command{
				Command: "down.document.type",
				Title:   typeLens,
			},
		})
	}

	// Knowledge graph lens
	if s.Graph != nil {
		entities := s.Graph.EntitiesByDocument(uri)
		if len(entities) > 0 {
			title := fmt.Sprintf("🧠 %d entities", len(entities))
			lens = append(lens, protocol.CodeLens{
				Range: protocol.Range{
					Start: protocol.Position{Line: 0, Character: 0},
					End:   protocol.Position{Line: 0, Character: 0},
				},
				Command: &protocol.Command{
					Command: "down.knowledge.summary",
					Title:   title,
				},
			})
		}
	}

	// Backlinks lens
	bl := s.ComputeBacklinks(uri)
	if bl.Count > 0 {
		blTitle := fmt.Sprintf("← %d backlinks", bl.Count)
		lens = append(lens, protocol.CodeLens{
			Range: protocol.Range{
				Start: protocol.Position{Line: 0, Character: 0},
				End:   protocol.Position{Line: 0, Character: 0},
			},
			Command: &protocol.Command{
				Command:   "down.backlinks",
				Title:     blTitle,
				Arguments: []interface{}{uri},
			},
		})
	}

	return lens, nil
}

// computeDocMetadata returns a metadata summary string for the code lens.
func (s *State) computeDocMetadata(uri string, text string) string {
	words := len(strings.Fields(text))
	if words == 0 {
		return ""
	}

	// Reading time (average 200 words/min)
	readingMinutes := words / 200
	if readingMinutes == 0 {
		readingMinutes = 1
	}

	// Count lines
	lines := strings.Count(text, "\n") + 1

	return fmt.Sprintf("📊 %d words · %d min read · %d lines", words, readingMinutes, lines)
}

func (s *State) LensResolve(c *glsp.Context, p *protocol.CodeLens) (*protocol.CodeLens, error) {
	return p, nil
}
