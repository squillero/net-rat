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
		slog.Debug("Clearing cache", "file", CacheFile)
		os.Remove(CacheFile)
	}

	ni := getNetInfo()
	slog.Debug("NetInfo", "info", ni)
	loopback := ni.GetType(LoopbackIP, PublicIP)
	tunnel := ni.GetType(TunnelIP, PublicIP)
	local := ni.GetType(LocalIP, PublicIP)
	public := ni.GetType(PublicIP, 0)

	slog.Debug("LoopbackIP", "val", loopback)
	slog.Debug("TunnelIP", "val", tunnel)
	slog.Debug("LocalIP", "val", local)
	slog.Debug("PublicIP", "val", public)

	local_s := Squeeze(local)
	if len(tunnel) > 0 {
		local_s += " via " + Squeeze(tunnel)
	}
	var public_s string
	if len(public) == 1 {
		public_s = public[0].Describe()
	} else {
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
	slog.Debug("*", "local_s", local_s)
	slog.Debug("*", "public_s", public_s)

	// 	return strings.Join(r[:], "/")

	if hostname, err := os.Hostname(); err == nil {
		fmt.Printf("%s @ ", hostname)
	}
	if local_s != "" && public_s != "" {
		fmt.Printf("%s // %s", local_s, public_s)
	} else {
		fmt.Printf("%s%s", local_s, public_s)
	}
	fmt.Println()
}
