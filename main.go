package main

import "flag"
import "log"
import "os"

import "github.com/fsnotify/fsnotify"

func main() {
	var target = flag.String("target", ".", "Specifiy a path to a file or directory to watch.")
	var recursive = flag.Bool("recursive", false, "Watch all subdirectories, if target is a directory.")
	var command = flag.String("command", "", "Command to execute as an event callback.")
	flag.Parse()
	log.Println(os.Args, *target, *recursive, *command)

	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Fatal(err)
	}
	defer watcher.Close()

	done := make(chan bool)
	go func() {
		for {
			select {
			case event := <-watcher.Events:
				switch {
				case event.Op&fsnotify.Write == fsnotify.Write:
					log.Println("WRITE:", event.Name)
				case event.Op&fsnotify.Create == fsnotify.Create:
					log.Println("CREAT:", event.Name)
				case event.Op&fsnotify.Remove == fsnotify.Remove:
					log.Println("REMOV:", event.Name)
				case event.Op&fsnotify.Rename == fsnotify.Rename:
					log.Println("RENAM:", event.Name)
				case event.Op&fsnotify.Chmod == fsnotify.Chmod:
					log.Println("CHMOD:", event.Name)
				}
			case err := <-watcher.Errors:
				log.Println("ERROR:", err)
			}
		}
	}()

	err = watcher.Add(*target)
	if err != nil {
		log.Fatal(err)
	}
	<-done
}
