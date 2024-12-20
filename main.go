package main

import (
	c "github.com/clpi/down.lsp/cli/cmd"
	"github.com/tliron/commonlog"
	_ "github.com/tliron/commonlog/simple"
	"log"
)

func main() {
	commonlog.Configure(2, nil)
	log.Print("Starting down.lsp...")
	c.Run()

}
