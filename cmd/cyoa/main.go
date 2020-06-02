package main

import "flag"

var storiesDir = flag.String("storiesDir", "stories/", "Directory in which stories (*.json) are stored.")

func main() {
	flag.Parse()
}
