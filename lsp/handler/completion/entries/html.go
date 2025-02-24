package entries

import (
	protocol "github.com/tliron/glsp/protocol_3_16"
)

var HtmlTagCompletions = func(i []protocol.CompletionItem) []protocol.CompletionItem {
	mode := protocol.InsertTextModeAsIs
	kind := protocol.CompletionItemKindInterface
	items := i
	htmltags := map[string]string{
		"a":     "<a href=\"\"></a>",
		"p":     "<p></p>",
		"h1":    "<h1></h1>",
		"h2":    "<h2></h2>",
		"h3":    "<h3></h3>",
		"h4":    "<h4></h4>",
		"h5":    "<h5></h5>",
		"h6":    "<h6></h6>",
		"ul":    "<ul></ul>",
		"ol":    "<ol></ol>",
		"li":    "<li></li>",
		"dl":    "<dl></dl>",
		"dt":    "<dt></dt>",
		"dd":    "<dd></dd>",
		"table": "<table></table>",
		"tr":    "<tr></tr>",
		"td":    "<td></td>",
		"th":    "<th></th>",
		"thead": "<thead></thead>",
		"tbody": "<tbody></tbody>",
		"tfoot": "<tfoot></tfoot>",
		"hr":    "<hr/>",
		"br":    "<br/>",
		"em":    "<em></em>",
		"div":   "<div></div>",
		"img":   "<img src=\"\" alt=\"\"/>",
	}
	t := true
	for tag, ins := range htmltags {
		desc := ins
		fmt := protocol.InsertTextFormatPlainText
		items = append(items, protocol.CompletionItem{
			Label:            tag,
			Documentation:    "# " + tag + "\n\n" + ins,
			InsertText:       &ins,
			Preselect:        &t,
			InsertTextFormat: &fmt,
			CommitCharacters: []string{"<"},
			InsertTextMode:   &mode,
			Detail:           &desc,
			Kind:             &kind,
		})
	}
	return items
}
