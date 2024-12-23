package workspace

import (
	"github.com/tliron/glsp"
	protocol "github.com/tliron/glsp/protocol_3_16"
)

var DownSettings map[string]interface{} = map[string]interface{}{
	"down": map[string]interface{}{
		"codeAction": map[string]interface{}{
			"enabled": true,
		},
		"codeLens": map[string]interface{}{
			"enabled": true,
		},
		"inlayHint": map[string]interface{}{
			"enabled": true,
		},
		"completion": map[string]interface{}{
			"enabled": true,
		},
		"enabled": true,
	},
	"markdown": map[string]interface{}{
		"enabled": true,
		"completion": map[string]interface{}{
			"enabled": true,
		},
	},
	"docdown": map[string]interface{}{
		"enabled": true,
		"codeAction": map[string]interface{}{
			"enabled": true,
		},
		"codeLens": map[string]interface{}{
			"enabled": true,
		},
		"inlayHint": map[string]interface{}{
			"enabled": true,
		},
		"completion": map[string]interface{}{
			"enabled": true,
		},
	},
}

func Configure(c *glsp.Context, p *protocol.DidChangeConfigurationParams) error {
	s := map[string]interface{}{
		"markdown": map[string]interface{}{},
		"down": map[string]interface{}{
			"enabled": true,
			"codeAction": map[string]interface{}{
				"enabled": true,
			},
			"codeLens": map[string]interface{}{
				"enabled": true,
			},
			"inlayHint": map[string]interface{}{
				"enabled": true,
			},
			"diagnostics": map[string]interface{}{
				"enabled": true,
			},
			"hover": map[string]interface{}{
				"enabled": true,
			},
			"completion": map[string]interface{}{
				"enabled": true,
			},
			"signatureHelp": map[string]interface{}{
				"enabled": true,
			},
		},
	}
	p.Settings = s
	p.Settings = s
	return nil
}
