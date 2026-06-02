package handler

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/tliron/glsp"
	protocol "github.com/tliron/glsp/protocol_3_16"
)

// WorkspaceDiagnosticKind classifies workspace-level issues.
type WorkspaceDiagnosticKind string

const (
	WsDiagOrphanDocument   WorkspaceDiagnosticKind = "orphan_document"
	WsDiagBrokenLink       WorkspaceDiagnosticKind = "broken_link"
	WsDiagCircularLink     WorkspaceDiagnosticKind = "circular_link"
	WsDiagStaleReference   WorkspaceDiagnosticKind = "stale_reference"
	WsDiagDuplicateHeading WorkspaceDiagnosticKind = "duplicate_heading"
	WsDiagEmptyDocument    WorkspaceDiagnosticKind = "empty_document"
)

// WorkspaceDiagnostic represents a workspace-level diagnostic.
type WorkspaceDiagnostic struct {
	Kind    WorkspaceDiagnosticKind `json:"kind"`
	URI     string                  `json:"uri"`
	Line    int                     `json:"line"`
	Message string                  `json:"message"`
}

// RunWorkspaceDiagnostics scans all open documents for workspace-level issues.
func (s *State) RunWorkspaceDiagnostics(ctx *glsp.Context) []WorkspaceDiagnostic {
	var diags []WorkspaceDiagnostic

	// Collect all known document URIs
	allDocs := make(map[string]bool)
	for uri := range s.Documents {
		allDocs[uri] = true
	}
	if s.Graph != nil {
		for _, ent := range s.Graph.Entities {
			for _, src := range ent.Sources {
				allDocs[src.URI] = true
			}
		}
	}

	// Check each document
	for uri, text := range s.Documents {
		lines := strings.Split(text, "\n")

		// Check for empty documents
		contentLines := 0
		for _, line := range lines {
			if strings.TrimSpace(line) != "" && !strings.HasPrefix(strings.TrimSpace(line), "---") {
				contentLines++
			}
		}
		if contentLines == 0 {
			diags = append(diags, WorkspaceDiagnostic{
				Kind:    WsDiagEmptyDocument,
				URI:     uri,
				Line:    0,
				Message: "Document is empty or contains only frontmatter",
			})
		}

		// Check for broken wiki links
		inCode := false
		for i, line := range lines {
			trimmed := strings.TrimSpace(line)
			if strings.HasPrefix(trimmed, "```") {
				inCode = !inCode
				continue
			}
			if inCode {
				continue
			}

			for _, m := range reLinkWiki.FindAllStringSubmatch(line, -1) {
				parts := strings.SplitN(m[1], "|", 2)
				target := strings.TrimSpace(parts[0])
				if !resolveWikiLinkExists(uri, target, allDocs) {
					diags = append(diags, WorkspaceDiagnostic{
						Kind:    WsDiagBrokenLink,
						URI:     uri,
						Line:    i,
						Message: "Broken wiki link: [[" + target + "]] — target not found",
					})
				}
			}
		}

		// Check for duplicate headings (same level, same text)
		headings := make(map[string]int)
		for i, line := range lines {
			trimmed := strings.TrimSpace(line)
			if strings.HasPrefix(trimmed, "#") {
				level := 0
				for _, ch := range trimmed {
					if ch == '#' {
						level++
					} else {
						break
					}
				}
				headingText := strings.TrimSpace(trimmed[level:])
				key := strings.ToLower(headingText)
				if prev, exists := headings[key]; exists {
					diags = append(diags, WorkspaceDiagnostic{
						Kind:    WsDiagDuplicateHeading,
						URI:     uri,
						Line:    i,
						Message: "Duplicate heading \"" + headingText + "\" (first at line " + intStr(prev+1) + ")",
					})
				} else {
					headings[key] = i
				}
			}
		}
	}

	// Check for orphan documents (not linked from anywhere)
	if s.Graph != nil {
		linkedDocs := make(map[string]bool)
		for _, rel := range s.Graph.Relations {
			if rel.Kind == "links_to" {
				if target, ok := s.Graph.Entities[rel.To]; ok {
					for _, src := range target.Sources {
						linkedDocs[src.URI] = true
					}
				}
			}
		}

		for uri := range s.Documents {
			if !linkedDocs[uri] {
				// Check if it's an index file (those are ok to be orphans)
				base := strings.ToLower(filepath.Base(strings.TrimPrefix(uri, "file://")))
				if base != "index.md" && base != "readme.md" && base != "home.md" {
					diags = append(diags, WorkspaceDiagnostic{
						Kind:    WsDiagOrphanDocument,
						URI:     uri,
						Line:    0,
						Message: "Orphan document: not linked from any other document",
					})
				}
			}
		}
	}

	// Publish diagnostics for broken links
	if ctx != nil {
		s.publishWorkspaceDiagnostics(ctx, diags)
	}

	return diags
}

// publishWorkspaceDiagnostics sends workspace-level diagnostics to the client.
func (s *State) publishWorkspaceDiagnostics(ctx *glsp.Context, wsDiags []WorkspaceDiagnostic) {
	// Group by URI
	byURI := make(map[string][]protocol.Diagnostic)
	source := "down-workspace"

	for _, d := range wsDiags {
		sev := protocol.DiagnosticSeverityInformation
		switch d.Kind {
		case WsDiagBrokenLink:
			sev = protocol.DiagnosticSeverityWarning
		case WsDiagCircularLink:
			sev = protocol.DiagnosticSeverityWarning
		case WsDiagOrphanDocument:
			sev = protocol.DiagnosticSeverityHint
		case WsDiagEmptyDocument:
			sev = protocol.DiagnosticSeverityHint
		case WsDiagDuplicateHeading:
			sev = protocol.DiagnosticSeverityInformation
		}

		byURI[d.URI] = append(byURI[d.URI], protocol.Diagnostic{
			Range: protocol.Range{
				Start: protocol.Position{Line: protocol.UInteger(d.Line), Character: 0},
				End:   protocol.Position{Line: protocol.UInteger(d.Line), Character: 100},
			},
			Severity: &sev,
			Source:   &source,
			Message:  d.Message,
		})
	}

	for uri, diags := range byURI {
		ctx.Notify(protocol.ServerTextDocumentPublishDiagnostics, protocol.PublishDiagnosticsParams{
			URI:         protocol.DocumentUri(uri),
			Diagnostics: diags,
		})
	}
}

// resolveWikiLinkExists checks if a wiki link target can be resolved.
func resolveWikiLinkExists(sourceURI, target string, knownDocs map[string]bool) bool {
	// Check if target matches any known document name
	targetLower := strings.ToLower(target)
	for uri := range knownDocs {
		base := filepath.Base(strings.TrimPrefix(uri, "file://"))
		nameNoExt := strings.TrimSuffix(base, filepath.Ext(base))
		if strings.EqualFold(nameNoExt, target) {
			return true
		}
		// Also check with dashes/underscores for spaces
		if strings.EqualFold(strings.ReplaceAll(nameNoExt, "-", " "), targetLower) {
			return true
		}
		if strings.EqualFold(strings.ReplaceAll(nameNoExt, "_", " "), targetLower) {
			return true
		}
	}

	// Check filesystem
	sourceDir := filepath.Dir(strings.TrimPrefix(sourceURI, "file://"))
	candidates := []string{
		filepath.Join(sourceDir, target+".md"),
		filepath.Join(sourceDir, target+".markdown"),
		filepath.Join(sourceDir, target),
		filepath.Join(sourceDir, strings.ReplaceAll(target, " ", "-")+".md"),
		filepath.Join(sourceDir, strings.ReplaceAll(target, " ", "_")+".md"),
	}
	for _, c := range candidates {
		if _, err := os.Stat(c); err == nil {
			return true
		}
	}

	return false
}
