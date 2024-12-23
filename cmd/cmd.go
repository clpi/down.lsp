package cmd

import (
	"context"
	"log"

	"github.com/clpi/down.lsp/cmd/config"
	"github.com/clpi/down.lsp/cmd/delete"
	"github.com/clpi/down.lsp/cmd/export"
	"github.com/clpi/down.lsp/cmd/find"
	"github.com/clpi/down.lsp/cmd/initialize"
	"github.com/clpi/down.lsp/cmd/list"
	logc "github.com/clpi/down.lsp/cmd/log"
	lsc "github.com/clpi/down.lsp/cmd/lsp"
	"github.com/clpi/down.lsp/cmd/new"
	"github.com/clpi/down.lsp/cmd/note"
	"github.com/clpi/down.lsp/cmd/serve"
	"github.com/clpi/down.lsp/cmd/shell"
	"github.com/clpi/down.lsp/cmd/sync"
	cmdutil "github.com/clpi/down.lsp/cmd/util"
	"github.com/clpi/down.lsp/cmd/workspace"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

func flag() (cmd *cobra.Command, f pflag.Flag) {
	return
}

var downR = func(cmd *cobra.Command, args []string) {
	log.Println(`down`)
}

var Down = cmdutil.Cmd("down", []string{"d"}, "down", "down", downR)

func Configure() {
	cobra.EnableCommandSorting = true
	cobra.EnablePrefixMatching = true
	Down.AddCommand(&lsc.Lsp)
	Down.AddCommand(&initialize.Init)
	Down.AddCommand(&Runc)
	Down.AddCommand(&workspace.Workspace)
	Down.AddCommand(&find.Find)
	Down.AddCommand(&list.List)
	Down.AddCommand(&config.Config)
	Down.AddCommand(&logc.Log)
	Down.AddCommand(&Tag)
	Down.AddCommand(&new.New)
	Down.AddCommand(&note.Note)
	Down.AddCommand(&Link)
	Down.AddCommand(&shell.Shell)
	Down.AddCommand(&serve.Serve)
	Down.AddCommand(&delete.Delete)
	Down.AddCommand(&export.Export)
	Down.AddCommand(&sync.Sync)
	Down.AddCommand(&Snippet)
	Down.AddCommand(&Template)
}

func Run(c *context.Context) {
	Configure()
	Down.Execute()
}
