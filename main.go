//        ____()()     NetRat v0.1
//       /      @@     ~~~~~~~~~~~
// `~~~~~\_;m__m._>o   Yet another little Go experiment
//
// Copyright © 2024 Giovanni Squillero <giovanni.squillero@polito.it>
// This code is being released for EDUCATIONAL and ACADEMIC purposes.
// +----------------------------------------------------------------+
// | <=*=> COMMERCIAL USE OF THE CODE IS EXPRESSLY PROHIBITED <=*=> |
// +----------------------------------------------------------------+
// The code is made available "AS-IS" without any express or implied
// guarantees such as fitness, merchantability, or non-infringement.

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
	log.SetFlags(0)
	log.Printf("This is NetRat v%s", NetRatVersion)

	// Parse flags
	flag.Parse()
	flag.BoolVar(&Verbose, "v", false, "Verbose operations")
	flag.BoolVar(&NoCache, "n", false, "Don't use cache")

	eager()
}
