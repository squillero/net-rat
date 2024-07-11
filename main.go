//        ____()()     NetRat v0.1
//       /      @@     ~~~~~~~~~~~
// `~~~~~\_;m__m._>o   A Go hack
//
// Coded in July 2024, between Italy and Australia (34,138 km).
// Copyright © 2024 Giovanni Squillero <giovanni.squillero@polito.it>
// Released for educational and academic purposes under 0BSD.

package main

import (
	"flag"
	"fmt"
	"log"
	"log/slog"
	"os"
)

const NetRatVersion = "0.1"

// FLAGS
var Verbose bool = false
var NoCache bool = false

func main() {
	log.SetPrefix("🐀 ") // 🐁
	log.SetFlags(log.Ldate | log.Ltime | log.Lmicroseconds | log.Lmsgprefix | log.LUTC)

	// Parse flags
	verbose := flag.Bool("v", false, "Verbose operations")
	clearCache := flag.Bool("c", false, "Clear cache")
	flag.Parse()
	if *verbose {
		slog.SetLogLoggerLevel(slog.LevelDebug)
	}
	slog.Debug("This is NetRat v" + NetRatVersion)
	if *clearCache {
		slog.Debug("Clearing cache", "file", CacheFile)
		os.Remove(CacheFile)
	}

	ip := getNetInfo()
	if ip[0].Valid() && ip[1].Valid() && ip[0].RawIp != ip[1].RawIp {
		fmt.Printf("%s/%s\n", ip[0], ip[1])
	} else if ip[1].Valid() {
		fmt.Printf("%s\n", ip[1])
	} else if ip[0].Valid() {
		fmt.Printf("%s (local only)\n", ip[0])
	} else {
		fmt.Println("Not connected")
	}
}
