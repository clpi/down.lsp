package cmd

import (
	"log"

	// dcmd "github.com/clpi/down.lsp/cli/cmd"
	"github.com/clpi/down.lsp/cli/cmd/config"
	"github.com/clpi/down.lsp/cli/cmd/delete"
	"github.com/clpi/down.lsp/cli/cmd/export"
	"github.com/clpi/down.lsp/cli/cmd/find"
	"github.com/clpi/down.lsp/cli/cmd/list"
	logc "github.com/clpi/down.lsp/cli/cmd/log"
	lsc "github.com/clpi/down.lsp/cli/cmd/lsp"
	"github.com/clpi/down.lsp/cli/cmd/new"
	"github.com/clpi/down.lsp/cli/cmd/note"
	"github.com/clpi/down.lsp/cli/cmd/shell"
	"github.com/clpi/down.lsp/cli/cmd/sync"
	"github.com/clpi/down.lsp/cli/cmd/workspace"
	"github.com/clpi/down.lsp/lsp"
	"github.com/spf13/cobra"
	// "github.com/tliron/commonlog"
	// "github.com/tliron/kutil/exec"
	// "github.com/tliron/kutil/terminal"
)

func Cmd(s string, l string, r func(*cobra.Command, []string)) cobra.Command {
	return cobra.Command{
		Aliases: []string{l},
		Use:     s,
		Long:    s,
		Short:   s,
		Run:     r,
	}
}

var (
	down = cobra.Command{
		Short:   "down",
		Use:     "down <command>",
		Example: "down lsp",
		Long:    "down",
		Version: lsp.Version,
		Hidden:  false,
		Args:    cobra.MinimumNArgs(0),
		PreRun: func(cmd *cobra.Command, args []string) {
			log.Print("down")
		},
		Run: func(cmd *cobra.Command, args []string) {
			log.Print("down")
		},
		PostRun: func(cmd *cobra.Command, args []string) {
			log.Print("down")
		},
	}
)

func Configure() {
	down.AddCommand(&lsc.Lsp)
	down.AddCommand(&Init)
	down.AddCommand(&Runc)
	down.AddCommand(&workspace.Workspace)
	down.AddCommand(&find.Find)
	down.AddCommand(&list.List)
	down.AddCommand(&config.Config)
	down.AddCommand(&logc.Log)
	down.AddCommand(&Tag)
	down.AddCommand(&new.New)
	down.AddCommand(&note.Note)
	down.AddCommand(&Link)
	down.AddCommand(&shell.Shell)
	down.AddCommand(&delete.Delete)
	down.AddCommand(&export.Export)
	down.AddCommand(&sync.Sync)
	down.AddCommand(&Snippet)
	down.AddCommand(&Template)
}
func Run() {
	Configure()
	down.Execute()
}
