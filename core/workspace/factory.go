package workspace

import protocol "github.com/tliron/glsp/protocol_3_16"

type (
	Factory interface {
		Invoke()
	}
)
type (
	AddNoteFunc func(name string, dir protocol.URI, template string) error
)
