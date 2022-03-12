package indexer

import (
	"github.com/fsnotify/fsnotify"
)

// Entry represent a new event when something happens in Download Folder
type Entry struct {
	event fsnotify.Event
	name  string
	path  string
	files []string
	isDir bool
}

// NewEntry permit to build an Entry object
func newEntry(event fsnotify.Event, name string, path string, files []string, isDir bool) Entry {
	return Entry{
		event: event,
		name:  name,
		path:  path,
		files: files,
		isDir: isDir,
	}
}

// IsDir says if the current Entry is a Directory or a File
func (e *Entry) IsDir() bool {
	return e.isDir
}

// GetPath return the root path of the Entry
func (e *Entry) GetPath() string {
	return e.path
}

// GetFiles return every MkvExt files from the root path
func (e *Entry) GetFiles() []string {
	return e.files
}

// GetEvent return the original Event of the Entry
func (e *Entry) GetEvent() fsnotify.Event {
	return e.event
}

func (e *Entry) GetName() string {
	return e.name
}
