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

	ni := getNetInfo()
	loopback := ni.GetType(LoopbackIP)
	local := ni.GetType(LocalIP)
	public := ni.GetType(PublicIP)

	if hostname, err := os.Hostname(); err == nil {
		fmt.Printf("%s @ ", hostname)
	}

	if local == "" && public == "" {
		fmt.Println(loopback + " (loopback only)")
	} else if public == "" {
		fmt.Println(local + " (local only)")
	} else if local == "" {
		fmt.Println(public)
	} else {
		fmt.Println(local + " // " + public)
	}
}
