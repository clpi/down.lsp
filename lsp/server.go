package lsp

import (
	"github.com/clpi/down.lsp/core/workspace"
	"github.com/clpi/down.lsp/lsp/handler"
	protocol "github.com/tliron/glsp/protocol_3_16"
	"github.com/tliron/glsp/server"
)

type (
	Info struct {
		Name    string
		Version string
	}
	Session any
	Server  struct {
		Server  *server.Server
		State   State
		Session Session
		Logger  *server.Logger
		Info    Info
	}
	State struct {
		Server     *server.Server
		Workspaces workspace.Workspaces
	}
)

func NewServer() (Server, error) {
	var handle protocol.Handler = handler.State{}.Handlers()
	return Server{
		Info: Info{
			Name:    Name,
			Version: Version,
		},
		Server: server.NewServer(&handle, Name, true),
		Logger: &server.Logger{},
	}, nil
}

func (s *Server) File(f string) {

}
func (s *Server) Handlers() {

}

