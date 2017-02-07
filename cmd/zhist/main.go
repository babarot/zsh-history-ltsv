package main

import (
	"log"
	"os"

	"gopkg.in/alecthomas/kingpin.v2"
)

func main() {
	var (
		search        = kingpin.Command("search", "Search timezones")
		searchFilters = search.Arg("filter", "Filter strings").Strings()
		update        = kingpin.Command("update", "Updates repositories from geonames")
	)

	cmd := kingpin.Parse()
	log.SetOutput(os.Stderr)

	switch cmd {
	case update.FullCommand():
		updateCommand()
	case search.FullCommand():
		// updateCommand()
		out := &terminalOutputer{}
		searchCommand(out, *searchFilters)
	}
}
