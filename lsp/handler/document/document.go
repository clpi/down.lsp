package document

import (
	files "github.com/clpi/down.lsp/lsp/files"
	protocol "github.com/tliron/glsp/protocol_3_16"
)

var (
	trueVal  = true
	workDone = protocol.WorkDoneProgressOptions{
		WorkDoneProgress: &trueVal,
	}
	falseVal          = false
	file              = "file"
	partialTok        = "1"
	TriggerCharacters = []string{
		" ", "@", "#", "$", "%", "&",
		"*", "+", "-", "/", "<", "=",
		">", "?", "^", "|", "~",
		"[", "(", "<", "{", "`",
		"]", ")", ">", "}",
		":", "=", ",",
	}
)

type (
	DocumentProvider struct {
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
		SignatureHelp      protocol.SignatureHelpOptions
	}
)

var (
	Provider = DocumentProvider{
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
		SignatureHelp: protocol.SignatureHelpOptions{
			WorkDoneProgressOptions: workDone,
			RetriggerCharacters:     TriggerCharacters,
			TriggerCharacters:       TriggerCharacters,
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
)
