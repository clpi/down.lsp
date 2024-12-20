package document

import (
	"log"
	"math"
	"os"
	"strings"

	"github.com/clpi/down.lsp/core/entities"
	"github.com/clpi/down.lsp/lsp/files"
	"github.com/clpi/down.lsp/lsp/util"
	"github.com/tliron/glsp"
	protocol "github.com/tliron/glsp/protocol_3_16"
)

const (
	Text  = 1
	Read  = 2
	Write = 3
)

func Hi(p *protocol.DocumentHighlightParams) protocol.DocumentHighlight {
	u := p.TextDocumentPositionParams.TextDocument.URI
	md, e := os.ReadFile(u)
	if e != nil {
		log.Println(e)
	}
	r := util.Range(p.Position.Line, 0, p.Position.Line, p.Position.EndOfLineIn(string(md)).Character)
	i := protocol.DocumentHighlightKindText
	return protocol.DocumentHighlight{
		Kind:  &i,
		Range: r,
	}

}
func Tx(p *protocol.DocumentHighlightParams) protocol.DocumentHighlight {
	return Hi(p)
}
func Rd(p *protocol.DocumentHighlightParams) protocol.DocumentHighlight {
	return Hi(p)
}
func Wr(p *protocol.DocumentHighlightParams) protocol.DocumentHighlight {
	return Hi(p)
}

const (
	WikiLinkOpenStart  = 1
	WikiLinkOpenEnd    = 2
	WikiLinkCloseStart = -1
	WikiLinkCloseEnd   = -2
	WikiLinkText       = 0
	WikiLinkWaiting    = math.MaxInt
	WikiLinkError      = -math.MaxInt
)

func Wikilink(p *protocol.DocumentHighlightParams) protocol.DocumentHighlight {
	u := p.TextDocument.URI
	md, e := os.ReadFile(u)
	ms := string(md)
	if e != nil {
		log.Println(e)
	}
	// s1 := ms[p.Position.IndexIn(ms)]
	// Init to iniinity, means is awaiting start of wikilink
	_ = entities.WikiLink{
		Range: util.Range(0, 0, 0, 0),
		Value: "",
	}
	inwl := 0

	for _, c := range md {
		switch c {
		case '[':
			switch inwl {
			// In link text
			case math.MaxInt:
				continue
			// case -math.MaxInt:
			// 	println("Error: unexpected '['")
			// 	inwl = 0
			// 	break
			// Not started parsing wikilink
			case 0 | 1:
				inwl += 1
			// Have parsed '[[', now in text
			case 2:
				inwl = -1
			// Have parse '[[...]', so '[' is error
			case 3:
				inwl = -math.MaxInt
			case 4:
				inwl = 0

			default:
				continue
			}
		case ']':
			switch inwl {
			case 0:
			case math.MaxInt:
				inwl = 1
			case 0 | 1:
				inwl = -1
			case 3 | 4:
				inwl = -1
			default:
				continue
			}
		default:
			break

		}
		if c == '\n' {

		}
	}
	_ = ms[p.Position.IndexIn(ms):p.Position.EndOfLineIn(ms).Character]
	strings.Split(ms, " ")
	_ = ms[p.Position.EndOfLineIn(ms).Character]

	// r := util.Rng(util.Start(p.Position.Line), util.Dc(p.Position.Line, end))

	if files.IsMarkdown(p.TextDocument.URI) {
		return Hi(p)
	}
	k := protocol.DocumentHighlightKindText
	return protocol.DocumentHighlight{
		Kind: &k,
		Range: protocol.Range{
			Start: protocol.Position{
				Line:      0,
				Character: 5,
			},
			End: protocol.Position{
				Line:      0,
				Character: 5,
			},
		},
	}
}

func DocumentHighlight(c *glsp.Context, p *protocol.DocumentHighlightParams) ([]protocol.DocumentHighlight, error) {
	var (
		kk                              = protocol.DocumentHighlightKindText
		hl []protocol.DocumentHighlight = []protocol.DocumentHighlight{}
		k1                              = protocol.DocumentHighlightKindText
		_                               = protocol.DocumentHighlight{
			Kind: &kk,
			Range: protocol.Range{
				Start: protocol.Position{
					Line:      0,
					Character: 5,
				},
				End: protocol.Position{
					Line:      0,
					Character: 5,
				},
			},
		}
		h1 protocol.DocumentHighlight = protocol.DocumentHighlight{
			Kind: &k1,
			Range: protocol.Range{
				Start: protocol.Position{
					Line:      0,
					Character: 10,
				},
				End: protocol.Position{
					Line:      0,
					Character: 10,
				},
			},
		}
	)
	return append(hl, h1), nil
}
