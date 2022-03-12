package main

import (
	"github.com/fsnotify/fsnotify"
	"github.com/pierreribault/go-plex-transfer/pkg/indexer"
	"log"
	"os"
	"time"
)

var (
	HostFolderName = os.Getenv("DOWNLOAD_FOLDER_NAME")
)

const (
	DownloadHostFolder = "Téléchargement"
	Downloads          = "/app/storage/downloads"
	Movies             = "/app/storage/movies"
	TVShows            = "/app/storage/tvshows"
)

func main() {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Fatalln("ERROR", err)
	}

	defer watcher.Close()

	log.Println("Indexer started, watching: ", Downloads)

	done := make(chan bool)
	go func() {
		for {
			select {
			case event, ok := <-watcher.Events:
				if !ok {
					return
				}
				time.Sleep(2 * time.Second)
				NewEventProvided(event)
			case err, ok := <-watcher.Errors:
				if !ok {
					return
				}
				log.Println("error:", err)
			}
		}
	}()

	err = watcher.Add(Downloads)
	if err != nil {
		log.Fatal(err)
	}
	<-done
}

func NewEventProvided(event fsnotify.Event) {
	if event.Name[len(event.Name)-1:] == "~" {
		return
	}

	if event.Op&fsnotify.Create == fsnotify.Create {
		CreateEvent(event)
	}

	if event.Op&fsnotify.Remove == fsnotify.Remove {
		RemoveEvent(event)
	}
}

func CreateEvent(event fsnotify.Event) {
	log.Println("New created event provided: ", event.Name)

	idx := CreateIndexer(event)
	if err := idx.Analyse(); err != nil {
		log.Println(err)
	}

	if err := idx.CreateSymbolicLink(); err != nil {
		log.Println(err)
	} else {
		log.Println("Symlink created")
	}
}

func RemoveEvent(event fsnotify.Event) {
	log.Println("New removed event provided: ", event.Name)

	idx := CreateIndexer(event)

	idx.RemoveSymbolicLink()
	log.Println("Symlink removed")
}

func CreateIndexer(event fsnotify.Event) indexer.Service {
	return indexer.New(event, Movies, TVShows, HostFolderName)
}
