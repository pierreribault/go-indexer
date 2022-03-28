package indexer

import (
	"errors"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/fsnotify/fsnotify"
)

const (
	MkvExtension = ".mkv"
	MP4Extension = ".mp4"
)

type Service struct {
	entry        Entry
	mediasPath   string
	downloadPath string
}

// New permit to build a new EntryService object
func New(event fsnotify.Event, mediasPath string, downloadPath string) Service {
	return Service{
		entry:        newEntry(event, extractName(event.Name), event.Name, []string{}, false),
		mediasPath:   mediasPath,
		downloadPath: downloadPath,
	}
}

func (s *Service) Analyse() error {
	file, _ := os.Open(s.entry.GetPath())
	fileInfo, _ := file.Stat()

	if s.entry.isDir = fileInfo.IsDir(); s.entry.IsDir() {
		s.entry.files, _ = glob(s.entry.GetPath())
		if len(s.entry.GetFiles()) == 0 {
			return errors.New(fmt.Sprintf("any mkv find in %s", s.entry.GetPath()))
		}
	} else {
		s.entry.files = append(s.entry.files, s.entry.GetPath())
	}

	return nil
}

// glob permit to find any files with specific ext recursively in dir folder.
func glob(dir string) ([]string, error) {
	var files []string
	err := filepath.Walk(dir, func(path string, f os.FileInfo, err error) error {
		switch filepath.Ext(path) {
		case MkvExtension, MP4Extension:
			files = append(files, path)
		}
		return nil
	})

	return files, err
}

// extractName get the text after the last /
func extractName(path string) string {
	position := strings.LastIndex(path, "/")

	if position == -1 {
		return path
	}

	return path[position+len("/"):]
}

// CreateSymbolicLink will create a ln -s on the right directory
func (s *Service) CreateSymbolicLink() error {
	s.RemoveSymbolicLink()

	target, output := s.generateTargetAndOutputPath()

	log.Println("Try to create a symlink to:", output)
	if err := os.Symlink(target, output); err != nil {
		return err
	}

	return nil
}

// RemoveSymbolicLink will remove every ln -s with a specific name in other directories
func (s *Service) RemoveSymbolicLink() {
	mediasPath := fmt.Sprintf("%s/%s", s.mediasPath, s.entry.GetName())

	if _, err := os.Lstat(mediasPath); err == nil {
		_ = os.Remove(mediasPath)
	}
}

// generateTargetAndOutputPath return the target and output name for the future symbolic link
func (s *Service) generateTargetAndOutputPath() (string, string) {
	target := fmt.Sprintf("../%s/%s", s.downloadPath, s.entry.GetName())
	output := fmt.Sprintf("%s/%s", s.mediasPath, s.entry.GetName())

	return target, output
}
