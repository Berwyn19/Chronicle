package cmd

import (
	"fmt"

	"github.com/berwyn/chronicle/internal/paths"
	"github.com/berwyn/chronicle/internal/store"
)

// openStore resolves the .chronicle directory for the current working
// directory and opens its database. It returns a helpful error if Chronicle
// has not been initialized here. Read commands (list, show, diff, file) use
// this rather than creating anything.
func openStore() (*store.Store, paths.Paths, error) {
	p, err := paths.ForCWD()
	if err != nil {
		return nil, paths.Paths{}, err
	}
	if !p.Exists() {
		return nil, paths.Paths{}, fmt.Errorf("not a chronicle project (no %s directory) — run 'chronicle init' first", paths.DirName)
	}
	s, err := store.Open(p.DB())
	if err != nil {
		return nil, paths.Paths{}, err
	}
	return s, p, nil
}
