//        ____()()     NetRat v0.1 -- Yet another Go experiment
//       /      @@     Copyright (c) 2024 by Giovanni Squillero
// `~~~~~\_;m__m._>o   Distributed under 0BSD (see LICENSE)

package main

import (
	"flag"
	"log"
	"os"
	"path/filepath"
)

const NetRatVersion = "0.1"

var CacheFile = filepath.Join(os.TempDir(), "netrat.json")

// FLAGS
var Verbose bool = false
var NoCache bool = false

func main() {
	log.SetPrefix("[NetRat] ")
	log.Printf("This is NetRat v%s", NetRatVersion)

	// Parse flags
	flag.BoolVar(&Verbose, "v", false, "Verbose operations")
	flag.BoolVar(&NoCache, "n", false, "Don't use cache")

	getMullvad()
	getAirVPN()
}
