package indexer

import (
	"errors"
	"fmt"
	"github.com/fsnotify/fsnotify"
	"log"
	"os"
	"path/filepath"
	"strings"
)

const (
	MkvExtension = ".mkv"
)

type Service struct {
	entry        Entry
	moviesPath   string
	tvShowsPath  string
	downloadPath string
}

// New permit to build a new EntryService object
func New(event fsnotify.Event, moviesPath string, tvShowsPath string, downloadPath string) Service {
	return Service{
		entry:        newEntry(event, extractName(event.Name), event.Name, []string{}, false),
		moviesPath:   moviesPath,
		tvShowsPath:  tvShowsPath,
		downloadPath: downloadPath,
	}
}

func (s *Service) Analyse() error {
	file, _ := os.Open(s.entry.GetPath())
	fileInfo, _ := file.Stat()

	if s.entry.isDir = fileInfo.IsDir(); s.entry.IsDir() {
		s.entry.files, _ = glob(s.entry.GetPath(), MkvExtension)
		if len(s.entry.GetFiles()) == 0 {
			return errors.New(fmt.Sprintf("any mkv find in %s", s.entry.GetPath()))
		}
	} else {
		s.entry.files = append(s.entry.files, s.entry.GetPath())
	}

	return nil
}

// glob permit to find any files with specify ext recursively in dir folder.
func glob(dir string, ext string) ([]string, error) {
	var files []string
	err := filepath.Walk(dir, func(path string, f os.FileInfo, err error) error {
		if filepath.Ext(path) == ext {
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

// IsMovie return true if the entry get only one file
func (s *Service) IsMovie() bool {
	return len(s.entry.GetFiles()) == 1
}

// IsTVShow return true if the entry get more than one file
func (s *Service) IsTVShow() bool {
	return len(s.entry.GetFiles()) > 1
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
	moviesPath := fmt.Sprintf("%s/%s", s.moviesPath, s.entry.GetName())
	tvShowPath := fmt.Sprintf("%s/%s", s.tvShowsPath, s.entry.GetName())

	if _, err := os.Lstat(moviesPath); err == nil {
		_ = os.Remove(moviesPath)
	}

	if _, err := os.Lstat(tvShowPath); err == nil {
		_ = os.Remove(tvShowPath)
	}
}

// generateTargetAndOutputPath return the target and output name for the future symbolic link
func (s *Service) generateTargetAndOutputPath() (string, string) {
	var path string

	if s.IsMovie() {
		path = s.moviesPath
	} else {
		path = s.tvShowsPath
	}

	target := fmt.Sprintf("../%s/%s", s.downloadPath, s.entry.GetName())
	output := fmt.Sprintf("%s/%s", path, s.entry.GetName())

	return target, output
}
