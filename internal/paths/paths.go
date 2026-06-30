// Package paths resolves the location of a Chronicle repository's files.
//
// Chronicle stores everything under a single .chronicle directory:
//
//	.chronicle/
//	    metadata.db
//	    events/
//	        <event-id>.patch
//
// This package is the one place that knows that layout. Nothing else should
// hard-code these names.
package paths

import (
	"os"
	"path/filepath"
)

// DirName is the name of the Chronicle directory created in a project root.
const DirName = ".chronicle"

// Paths holds the resolved locations of a Chronicle repository's files.
type Paths struct {
	// Root is the path to the .chronicle directory itself.
	Root string
}

// New returns the Paths for a .chronicle directory located directly inside dir.
func New(dir string) Paths {
	return Paths{Root: filepath.Join(dir, DirName)}
}

// ForCWD returns the Paths for a .chronicle directory in the current working
// directory.
func ForCWD() (Paths, error) {
	cwd, err := os.Getwd()
	if err != nil {
		return Paths{}, err
	}
	return New(cwd), nil
}

// DB is the path to the SQLite metadata database.
func (p Paths) DB() string {
	return filepath.Join(p.Root, "metadata.db")
}

// EventsDir is the directory holding per-event Git patches.
func (p Paths) EventsDir() string {
	return filepath.Join(p.Root, "events")
}

// Patch is the path to the patch file for a given event ID.
func (p Paths) Patch(eventID string) string {
	return filepath.Join(p.EventsDir(), eventID+".patch")
}

// Exists reports whether the .chronicle directory already exists.
func (p Paths) Exists() bool {
	info, err := os.Stat(p.Root)
	return err == nil && info.IsDir()
}
