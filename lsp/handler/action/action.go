package action

import (
	"github.com/tliron/glsp"
	protocol "github.com/tliron/glsp/protocol_3_16"
)

var (
	trueVal  bool                       = true
	falseVal bool                       = false
	v        protocol.Integer           = 0
	src                                 = protocol.CodeActionKindSource
	Provider protocol.CodeActionOptions = protocol.CodeActionOptions{
		CodeActionKinds: []protocol.CodeActionKind{
			protocol.CodeActionKindSource,
			protocol.CodeActionKindQuickFix,
			protocol.CodeActionKindRefactor,
			protocol.CodeActionKindRefactorExtract,
			protocol.CodeActionKindRefactorInline,
			protocol.CodeActionKindRefactorRewrite,
			protocol.CodeActionKindSourceOrganizeImports,
		},
		WorkDoneProgressOptions: protocol.WorkDoneProgressOptions{
			WorkDoneProgress: &trueVal,
		},
		ResolveProvider: &trueVal,
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
		Kind:        &src,
		Data:        nil,
		Edit:        &protocol.WorkspaceEdit{},
		Title:       "Create link on word/heading",
	}
	generateToc protocol.CodeAction = protocol.CodeAction{
		Command: &protocol.Command{
			Arguments: []any{
				"dir",
			},
			Title:   "Generate Table of Contents",
			Command: "down.toc.generate",
		},
		Kind:        &src,
		Title:       "Generate Table of Contents",
		Diagnostics: nil,
		IsPreferred: &trueVal,
		Data:        nil,
		Edit: &protocol.WorkspaceEdit{
			DocumentChanges: []any{
				protocol.TextDocumentEdit{
					TextDocument: protocol.OptionalVersionedTextDocumentIdentifier{
						Version: &v,
						TextDocumentIdentifier: protocol.TextDocumentIdentifier{
							URI: "file:///path/to/file.md",
						},
					},
					Edits: []any{
						protocol.TextEdit{
							Range: protocol.Range{
								Start: protocol.Position{
									Line:      0,
									Character: 0,
								},
								End: protocol.Position{
									Line:      0,
									Character: 0,
								},
							},
							NewText: "## Table of Contents\n\n",
						},
					},
				},
			},
			Changes:           nil,
			ChangeAnnotations: nil,
		},
		Disabled: &struct {
			Reason string `json:"reason"`
		}{},
	}
)

func CodeAction(c *glsp.Context, p *protocol.CodeActionParams) (any, error) {
	var (
		actions []protocol.CodeAction = []protocol.CodeAction{}
		_       protocol.Range        = protocol.Range{
			Start: protocol.Position{
				Line:      10,
				Character: 10,
			},
			End: protocol.Position{
				Line:      10,
				Character: 20,
			},
		}
	)
	return append(actions, cursorCreateLink, generateToc), nil
}
