package command

import (
	protocol "github.com/tliron/glsp/protocol_3_16"
)

var (
	trueVal  = true
	falseVal = false
	Commands = []string{
		"down.index",
		"down.log.new",
		"down.calendar.open",
		"down.save",
		"down.template.new",
		"down.template.open",
		"down.template.delete",
		"down.template.index",
		"down.snippet.new",
		"down.snippet.open",
		"down.snippet.delete",
		"down.snippet.index",
		"down.snippet.cursor",
		"down.load",
		"down.capture",
		"down.note.index",
		"down.note.today",
		"down.note.yesterday",
		"down.note.tomorrow",
		"down.note.month",
		"down.note.year",
		"down.task.index",
		"down.task.new",
		"down.task.today",
		"down.task.list",
		"down.task.delete",
		"down.log.index",
		"down.log.delete",
		"down.workspace.open",
		"down.workspace.new",
		"down.workspace.delete",
		"down.link.backlinks",
		"down.link.create",
		"down.link.create.cursor",
	}
	Provider protocol.ExecuteCommandOptions = protocol.ExecuteCommandOptions{
		WorkDoneProgressOptions: protocol.WorkDoneProgressOptions{
			WorkDoneProgress: &trueVal,
		},
		Commands: Commands,
	}
)
