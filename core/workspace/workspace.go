package workspace

import (
	"path/filepath"

	store "github.com/clpi/down.lsp/core/store"
	_ "go.lsp.dev/protocol"
)

type (
	WorkspaceConfig = map[string]interface{}
)

var (
	DefaultWorkspaceConfig = map[string]interface{}{
		"workspace": map[string]interface{}{
			"indexName": "index.md",
			"directories": map[string]interface{}{
				"store": map[string]interface{}{
					"default": ".down",
				},
				"notes": map[string]interface{}{
					"default": "notes",
				},
				"templates": map[string]interface{}{
					"default": "templates",
				},
				"snippets": map[string]interface{}{
					"default": "snippets",
				},
			},
		},
	}
)

type (
	Name       string
	Identifier string
)

var ()

type Workspace store.Store[Name, store.Store[Identifier, interface{}]]

type Workspaces store.Store[string, Workspace]

func NewWorkspace(id string) (Workspace, error) {
	return Workspace{}, nil
}

func (w Workspace) Path(d ...string) string {
	return filepath.Join(string(w.Uri), filepath.Join(d...))
}
func (w Workspace) Index() string {
	return w.Path("index.md")

}
