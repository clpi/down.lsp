package semantic

import (
	"github.com/clpi/down.lsp/lsp/files"
	protocol "github.com/tliron/glsp/protocol_3_16"
)

var (
	id         string = "semanticTokens"
	trueVal           = true
	falseVal          = false
	TokenTypes        = []string{
		"namespace",
		"type",
		"class",
		"enum",
		"interface",
		"struct",
		"typeParameter",
		"parameter",
		"variable",
		"property",
		"enumMember",
		"event",
		"function",
		"method",
		"macro",
		"keyword",
		"modifier",
		"comment",
		"string",
		"number",
		"regexp",
		"operator",
	}
	workDone = protocol.WorkDoneProgressOptions{
		WorkDoneProgress: &trueVal,
	}
	Provider = protocol.SemanticTokensOptions{
		Full:                    true,
		Range:                   true,
		WorkDoneProgressOptions: workDone,
		Legend: protocol.SemanticTokensLegend{
			TokenTypes: TokenTypes,
			TokenModifiers: []string{
				"declaration",
				"definition",
				"readonly",
				"static",
				"deprecated",
				"abstract",
				"async",
				"modification",
				"documentation",
				"defaultLibrary",
			},
		},
	}
	Registration = protocol.SemanticTokensRegistrationOptions{
		TextDocumentRegistrationOptions: protocol.TextDocumentRegistrationOptions{
			DocumentSelector: &files.Filetypes,
		},
		SemanticTokensOptions: Provider,
		StaticRegistrationOptions: protocol.StaticRegistrationOptions{
			ID: &id,
		},
	}
)
