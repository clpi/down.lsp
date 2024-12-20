package entries

import protocol "github.com/tliron/glsp/protocol_3_16"

type (
	Template struct {
		Body        string
		Document    protocol.DocumentUri
		Workspace   string
		Description string
		URI         protocol.URI
	}
	Log struct {
		Body string
	}
	Note struct {
		Body string
	}
	Document struct {
	}
	Snippet struct {
		Body        string
		Description string
	}
	Workspace struct {
		Name    string
		Path    protocol.URI
		Default bool
		Index   string
		Notes   string
		Config  interface{}
	}
)
