package main

import "flag"
import "log"
import "os/exec"

import "github.com/fsnotify/fsnotify"

type Options struct {
	Recursive bool
	Command   string
	Paths     []string
}

func main() {
	var target = flag.String("target", ".", "Specifiy a path to a file or directory to watch.")
	var recursive = flag.Bool("recursive", false, "Watch all subdirectories, if target is a directory.")
	var command = flag.String("command", "", "Command to execute as an event callback.")
	flag.Parse()

	cmd, err := exec.LookPath(*command)
	if err != nil {
		log.Fatal(err)
	}

	opts := Options{
		Recursive: *recursive,
		Command:   cmd,
		Paths:     nil,
	}

	log.Println(opts)

	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Fatal(err)
	}
	defer watcher.Close()

	// TODO: Make recursive!
	if *recursive {
	} else {
		err = watcher.Add(*target)
		if err != nil {
			log.Fatal(err)
		}
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
				case event.Op&fsnotify.Remove == fsnotify.Remove:
					log.Println("REMOVE: ", event.Name)
				case event.Op&fsnotify.Rename == fsnotify.Rename:
					log.Println("RENAME: ", event.Name)
				case event.Op&fsnotify.Chmod == fsnotify.Chmod:
					log.Println("CHMOD:  ", event.Name)
				}
			case err := <-watcher.Errors:
				log.Println("ERROR:", err)
			}
		}
	}()
	// TODO: Make signals get handled and that is how to kill the process
	<-done
}
