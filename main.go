package main

import "flag"
import "fmt"
import "os"

//import "github.com/fsnotify/fsnotify"

func main() {
	var watchTarget = flag.String("target", ".", "Specifiy a path to a file or directory to watch.")
	var recursive = flag.Bool("recursive", false, "Watch all subdirectories, if target is a directory.")
	var command = flag.String("command", "", "Command to execute as an event callback.")
	flag.Parse()

	fmt.Println(os.Args, *watchTarget, *recursive, *command)
}
