package main

import (
	"log"
	// cmd "github.com/clpi/down.lsp/cmd"
	lsp "github.com/clpi/down.lsp/lsp"
	"github.com/tliron/commonlog"
	_ "github.com/tliron/commonlog/simple"
)

func args() any {
	return nil
}
func main() {
	commonlog.Configure(2, nil)
	log.Print("Starting down.lsp...")
	down, error := lsp.NewLsp()
	if error != nil {
		log.Fatal(error)
	}
	err := down.Stdio()
	if err != nil {
		log.Fatal(err)
	}

}
