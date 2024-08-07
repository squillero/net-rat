//        ____()()     NetRat v0.2
//       /      @@     ~~~~~~~~~~~
// `~~~~~\_;m__m._>o   A Go experiment
//
// Coded in July 2024, between Italy and Australia (34,138 km).
// Copyright © 2024 Giovanni Squillero <giovanni.squillero@polito.it>
// Released under 0BSD (see LICENSE).

package main

import (
	"flag"
	"fmt"
	"log"
	"log/slog"
	"os"
	"strings"
)

const NetRatVersion = "0.2"

func Squeeze(ips []IpInfo) string {
	var r []string
	for _, v := range ips {
		r = append(r, v.Describe())
	}
	return strings.Join(r[:], "/")
}

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
		slog.Debug("Removing cache", "file", CacheFile)
		os.Remove(CacheFile)
	}

	ni := getNetInfo()
	slog.Debug("NetInfo", "info", ni)
	loopback := ni.GetType(LoopbackIP, PublicIP)
	tunnel := ni.GetType(TunnelIP, PublicIP)
	local := ni.GetType(LocalIP, PublicIP)
	public := ni.GetType(PublicIP, 0)

	loopback_s := Squeeze(loopback)
	local_s := Squeeze(local)
	if len(tunnel) > 0 {
		local_s += " via " + Squeeze(tunnel)
	}
	var public_s string
	if len(public) == 1 {
		public_s = public[0].Describe()
	} else if len(public) > 1 {
		var d []string
		for _, v := range public {
			if v.Comment != "" {
				d = append(d, v.Comment)
			}
		}
		public_s = "Multi NAT"
		if len(d) > 0 {
			public_s += " [" + strings.Join(d[:], "; ") + "]"
		}
	}

	// 	return strings.Join(r[:], "/")

	if hostname, err := os.Hostname(); err == nil {
		fmt.Printf("%s @ ", hostname)
	}

	switch {
	case local_s == "" && public_s == "" && loopback_s == "":
		fmt.Print("no network")
	case local_s == "" && public_s == "" && loopback_s != "":
		fmt.Printf("%s [Loopback Only]", loopback_s)
	case local_s != "" && public_s == "":
		fmt.Printf("%s [Local Only]", loopback_s)
	case local_s != "" && public_s != "":
		fmt.Printf("%s // %s", local_s, public_s)
	case local_s == "" && public_s != "":
		fmt.Printf("%s", public_s)
	}
	fmt.Println()
}
