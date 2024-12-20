package lsp

import (
	"github.com/clpi/down.lsp/lsp/handler"
	proto "github.com/tliron/glsp/protocol_3_16"
	"github.com/tliron/glsp/server"
)

const (
	Name     string = "down"
	Path     string = "down.lsp"
	Version  string = "0.1.0-alpha"
	Protocol string = "3.16"
)

type (
	Client struct {
	}
	Lsp struct {
		Name    string
		Version string
		Server  *server.Server
		Logger  *server.Logger
	}
)

func NewLsp() (Lsp, error) {
	var handle proto.Handler = handler.LspHandler()
	return Lsp{
		Name:    Name,
		Version: Version,
		Server:  server.NewServer(&handle, Name, true),
		Logger:  &server.Logger{},
	}, nil
}

func (lsp *Lsp) NodeJs() error {
	return lsp.Server.RunNodeJs()
}

func (lsp *Lsp) Stdio() error {
	return lsp.Server.RunStdio()
}
