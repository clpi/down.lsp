package handler

import (
	"github.com/tliron/glsp"
	// eterm "github.com/tliron/kutil/terminal"
	// event "go.lsp.dev/pkg/event"
	// golsp "go.lsp.dev/protocol"
	// stream "github.com/tliron/kutil/protobuf"
	// fswatch "github.com/tliron/kutil/fswatch"
	protocol "github.com/tliron/glsp/protocol_3_16"
	// util "github.com/tliron/kutil/util"
)

type (
	any = interface{}
	Env = map[string]string
)

var (
	t              bool                       = true
	f              bool                       = false
	v              protocol.Integer           = 0
	src                                       = protocol.CodeActionKindSource
	ActionProvider protocol.CodeActionOptions = protocol.CodeActionOptions{
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
			WorkDoneProgress: &t,
		},
		ResolveProvider: &t,
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
		IsPreferred: &t,
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
		IsPreferred: &t,
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

func (s *State) CodeAction(c *glsp.Context, p *protocol.CodeActionParams) (any, error) {
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

// ! Resolve
func (s *State) ActionResolve(c *glsp.Context, p *protocol.CodeAction) (*protocol.CodeAction, error) {
	return p, nil
}
