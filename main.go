package main

import (
	"flag"
	"log"
	"os/exec"
	"path/filepath"

	"github.com/fsnotify/fsnotify"
)

type Options struct {
	Recursive bool
	Command   string
	Paths     []string
}

func parseOptionsFromFlags() Options {
	var target = flag.String("target", ".", "Specifiy a path to a file or directory to watch.")
	var recursive = flag.Bool("recursive", false, "Watch all subdirectories, if target is a directory.")
	var command = flag.String("command", "echo", "Command to execute as an event callback.")
	flag.Parse()

	cmd, err := exec.LookPath(*command)
	if err != nil {
		log.Fatal("Option Parsing error: ", err)
	}

	path, err := filepath.Abs(*target)
	if err != nil {
		log.Fatal("Error expanding target: ", *target, "\n", err)
	}

	opts := Options{
		Recursive: *recursive,
		Command:   cmd,
		Paths:     []string{path},
	}

	// TODO: opts.Paths = []string of path(s).
	// if a directory was given: contains all files that directory
	//   and if recursive, traverse all subdirectories
	// if a file was given, only watch that file

	return opts
}

func newWatcher(o *Options) *fsnotify.Watcher {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Fatal("newWatcher: ", err)
	}

	if o.Recursive {
		// iterate over o.Paths and recruse down each subirectory unless it is a symlink
	} else {
		// iterate over o.Paths and add each to watcher
	}

	return watcher
}

func main() {
	var options Options = parseOptionsFromFlags()
	var watcher *fsnotify.Watcher = newWatcher(&options)
	defer watcher.Close()

	// TODO: Make recursive!
	err := watcher.Add(options.Paths[0])
	if err != nil {
		log.Fatal("Add to watcher error: ", options.Paths[0], " -- ", err)
	}

	done := make(chan bool)
	go func() {
		for {
			select {
			case event := <-watcher.Events:
				switch {
				case event.Op&fsnotify.Write == fsnotify.Write:
					log.Println("WRITE:  ", event.Name)
				case event.Op&fsnotify.Create == fsnotify.Create:
					log.Println("CREATE: ", event.Name)
					// TODO: add to watch and options.Paths
					// IF recursive and isDir(event.Name), add directory to watch
				case event.Op&fsnotify.Remove == fsnotify.Remove:
					log.Println("REMOVE: ", event.Name)
					// TODO: Remove from watch and options.Paths
				case event.Op&fsnotify.Rename == fsnotify.Rename:
					log.Println("RENAME: ", event.Name)
					// TODO: update options.Paths array with new name
					// Remove this file, a CREATE will also happen and, if need be, add to paths
				case event.Op&fsnotify.Chmod == fsnotify.Chmod:
					log.Println("CHMOD:  ", event.Name)
				}
				// TODO: run command
			case err := <-watcher.Errors:
				log.Println("ERROR:", err)
			}
		}
	}()
	// TODO: Make signals get handled and that is how to kill the process
	<-done
}
