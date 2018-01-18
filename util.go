package main

import "flag"
import "fmt"
import "os"
import "os/exec"
import "path/filepath"

func SetupWatch(paths []string, excludes []string) (int, *inotify.Watcher) {
	// How many directories are being watched
	var watchedCount int

	// Collect all subdirs of the watch and exclude roots
	paths = CollectPaths(paths)
	excludes = CollectPaths(excludes)

	watcher, err := inotify.NewWatcher()
	if err != nil {
		fmt.Println("Error establishing watcher: ", err)
	}

	// establish watches
	for _, path := range paths {
		if IndexOf(path, excludes) == -1 {
			err = watcher.Watch(path)
			if err != nil {
				fmt.Println("Error: ", err, "  establishing watch on: ", path)
			}
			watchedCount++
		}
	}
	return watchedCount, watcher
}

func IndexOf(element string, array []string) int {
	for i := 0; i < len(array); i++ {
		if array[i] == element {
			return i
		}
	}
	return -1
}

func CollectPaths(paths []string) []string {
	// paths to be returned
	collectedPaths := make([]string, 1, 1)

	for _, thisPath := range paths {
		err := filepath.Walk(filepath.Clean(thisPath),
			// Function arg for filepath.Walk
			func(path string, info os.FileInfo, err error) error {
				if info == nil {
					fmt.Println("File or directory does not exist.", path)
				} else if info.IsDir() {
					collectedPaths = append(collectedPaths, path)
				}
				return nil
			})

		if err != nil {
			fmt.Println(err)
		}
	}
	return collectedPaths
}
