package entries

import (
	protocol "github.com/tliron/glsp/protocol_3_16"
)

var (
	FileCompletions = func(i []protocol.CompletionItem) []protocol.CompletionItem {
		items := i
		return items
	}
)
